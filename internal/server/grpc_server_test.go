package server

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/server/storage"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const (
	_envTestDatabaseDsn = "TEST_DATABASE_DSN"
)

type GrpcServerSuite struct {
	suite.Suite
	testSrv   *GrpcServer
	testClt   pb.GophKeeperServiceClient
	stg       storage.Storage
	logger    *logrus.Logger
	fTeardown func()
}

func (gs *GrpcServerSuite) SetupSuite() {
	var err error
	gs.logger = logrus.New()
	gs.stg, err = storage.NewPgStorage(os.Getenv(_envTestDatabaseDsn), gs.logger)
	require.NoError(gs.T(), err)
}

func (gs *GrpcServerSuite) SetupTest() {
	gs.createTestServer()
}

func (gs *GrpcServerSuite) TearDownTest() {
	gs.fTeardown()
}

func (gs *GrpcServerSuite) TestCreateUser() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	gs.Run("create new user", func() {
		token, err := gs.testClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
		gs.NoError(err)
		gs.NotNil(token)
		gs.NotEqual("", token.Token)
	})

	gs.Run("create user with invalid credentials", func() {
		_, err := gs.testClt.UserCreate(ctx, &pb.User{Login: "", Password: ""})
		gs.Error(err)
		status, ok := status.FromError(err)
		gs.True(ok)
		gs.Equal(codes.InvalidArgument, status.Code())
	})

	gs.Run("create same user twice", func() {
		_, err := gs.testClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
		gs.Error(err)
		status, ok := status.FromError(err)
		gs.True(ok)
		gs.Equal(codes.AlreadyExists, status.Code())
	})
}

func (gs *GrpcServerSuite) TestLoginUser() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	userLogin, userPassword := uuid.NewString(), uuid.NewString()

	gs.Run("login, user not exists", func() {
		_, err := gs.testClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
		gs.Error(err)
		status, ok := status.FromError(err)
		gs.True(ok)
		gs.Equal(codes.NotFound, status.Code())
	})

	gs.Run("create new user", func() {
		token, err := gs.testClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
		gs.NoError(err)
		gs.NotNil(token)
		gs.NotEqual("", token.Token)
	})

	gs.Run("login, succesfull", func() {
		token, err := gs.testClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
		gs.NoError(err)
		gs.NotNil(token)
		gs.NotEqual("", token.Token)
	})

	gs.Run("login, wrong password", func() {
		_, err := gs.testClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: "foobar"})
		gs.Error(err)
		status, ok := status.FromError(err)
		gs.True(ok)
		gs.Equal(codes.PermissionDenied, status.Code())
	})
}

func (gs *GrpcServerSuite) TestSecretGetList() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, token := gs.createTestUser(ctx)

	gs.Run("get empty secret list", func() {
		_, err := gs.testClt.SecretGetList(contextWithToken(ctx, token), &pb.Empty{})
		gs.Error(err)
		status, ok := status.FromError(err)
		gs.True(ok)
		gs.Equal(codes.NotFound, status.Code())
	})

	// TODO create add secret
	// Check list after create secrets
}

func (gs *GrpcServerSuite) TestStoragePing() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, _, token := gs.createTestUser(ctx)
	_, err := gs.testClt.Ping(contextWithToken(ctx, token), &pb.Empty{})
	gs.NoError(err)
}

func TestGrpcServerSuite(t *testing.T) {
	// TODO: REMOVE AFTER ALL TESTS
	t.Setenv(_envTestDatabaseDsn, "postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable")
	//
	_, ok := os.LookupEnv(_envTestDatabaseDsn)
	if !ok {
		t.Skip("Test environment not set")
		return
	}
	suite.Run(t, new(GrpcServerSuite))
}

func (gs *GrpcServerSuite) createTestServer() {
	buffer := 101024 * 1024
	lis := bufconn.Listen(buffer)
	serverSecret := []byte("GophKeeperTest")

	srvCredentials, cltCredentials := getServerCredentials(), getClientCredentials()

	var grpcSrv *grpc.Server
	grpcSrv, gs.testSrv = NewGrpcServer(gs.stg, srvCredentials, serverSecret, gs.logger)

	go func() {
		_ = grpcSrv.Serve(lis)
	}()

	conn, err := grpc.DialContext(context.Background(), "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}),
		grpc.WithTransportCredentials(cltCredentials))
	require.NoError(gs.T(), err)

	gs.fTeardown = func() {
		lis.Close()
		grpcSrv.Stop()
	}

	gs.testClt = pb.NewGophKeeperServiceClient(conn)
}

func (gs *GrpcServerSuite) createTestUser(ctx context.Context) (string, string, string) {
	userLogin, userPassword := uuid.NewString(), uuid.NewString()
	token, err := gs.testClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
	require.NoError(gs.T(), err)

	return userLogin, userPassword, token.Token
}

func contextWithToken(ctx context.Context, tokenString string) context.Context {
	md := metadata.New(map[string]string{token.HeaderName: tokenString})
	return metadata.NewOutgoingContext(ctx, md)
}

func getServerCredentials() credentials.TransportCredentials {
	tlsServerSettings, _ := gkTLS.NewServerSettings(getTLSFile("server-cert.pem"), getTLSFile("server-key.pem"))
	tlsCredentials, _ := tlsServerSettings.Load()
	return tlsCredentials
}

func getClientCredentials() credentials.TransportCredentials {
	tlsCredentials, _ := gkTLS.LoadCACert(getTLSFile("ca-cert.pem"), "127.0.0.1")
	return tlsCredentials
}

func getTLSFile(fileName string) string {
	_, this, _, _ := runtime.Caller(0)
	tlsRoot := filepath.Join(this, "../../../tls")
	return filepath.Join(tlsRoot, fileName)
}
