package server

import (
	"context"
	"errors"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/grpc/interceptor"
	"github.com/devldavydov/gophkeeper/internal/server/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	MetaUserID = "USER_ID"

	_userTokenExpiration = 24 * time.Hour

	_msgPingFailed                 = "ping failed"
	_msgUserCredentialsBadRequest  = "invalid credentials" //nolint:gosec // Ok
	_msgUserFailedToCreate         = "failed to create user"
	_msgUserAlreadyExists          = "user already exists"
	_msgUserFailedToCreateToken    = "failed to create user token"
	_msgUserInvalidLoginOrPassword = "user invalid login or password"
	_msgUserNotFound               = "user not found"
	_msgUserFailedToLogin          = "failed to login"
)

// GrpcServer represents gRPC server.
type GrpcServer struct {
	pb.UnimplementedGophKeeperServiceServer
	stg          storage.Storage
	serverSecret []byte
	logger       *logrus.Logger
}

// NewGrpcServer creates new GRPCServer object.
func NewGrpcServer(
	stg storage.Storage,
	tlsCredentials credentials.TransportCredentials,
	serverSecret []byte,
	logger *logrus.Logger,
) (*grpc.Server, *GrpcServer) {
	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			interceptor.NewAuthTokenInterceptor(
				[]string{
					pb.GophKeeperService_SecretGetList_FullMethodName,
					pb.GophKeeperService_SecretGet_FullMethodName,
					pb.GophKeeperService_SecretCreate_FullMethodName,
					pb.GophKeeperService_SecretUpdate_FullMethodName,
					pb.GophKeeperService_SecretDelete_FullMethodName,
					pb.GophKeeperService_Ping_FullMethodName,
				},
				serverSecret).Handle),
	}

	opts = append([]grpc.ServerOption{grpc.Creds(tlsCredentials)}, opts...)

	grpcSrv := grpc.NewServer(opts...)
	srv := &GrpcServer{stg: stg, serverSecret: serverSecret, logger: logger}
	pb.RegisterGophKeeperServiceServer(grpcSrv, srv)
	return grpcSrv, srv
}

func (g *GrpcServer) UserCreate(ctx context.Context, user *pb.User) (*pb.UserAuthToken, error) {
	if user.Login == "" || user.Password == "" {
		g.logger.Errorf("invalid user credentials: user='%s' password='%s'", user.Login, user.Password)
		return nil, status.Error(codes.InvalidArgument, _msgUserCredentialsBadRequest)
	}

	pwdHash, err := hashPassword(user.Password)
	if err != nil {
		g.logger.Errorf("create user [%s] password hash error: %v", user.Login, err)
		return nil, status.Error(codes.Internal, _msgUserFailedToCreate)
	}

	userID, err := g.stg.CreateUser(ctx, user.Login, pwdHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserAlreadyExists) {
			g.logger.Errorf("failed to create user [%s]: already exists", user.Login)
			return nil, status.Error(codes.AlreadyExists, _msgUserAlreadyExists)
		}
		g.logger.Errorf("failed to create user [%s]: %v", user.Login, err)
		return nil, status.Error(codes.Internal, _msgUserFailedToCreate)
	}

	token, err := token.NewJWTForUser(userID, _userTokenExpiration, g.serverSecret)
	if err != nil {
		g.logger.Errorf("failed to create token for user [%s, %d]: %v", user.Login, userID, err)
		return nil, status.Error(codes.Unavailable, _msgUserFailedToCreateToken)
	}

	return &pb.UserAuthToken{Token: token}, nil
}

func (g *GrpcServer) UserLogin(ctx context.Context, user *pb.User) (*pb.UserAuthToken, error) {
	userID, pwdHash, err := g.stg.FindUser(ctx, user.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			g.logger.Errorf("login error for [%s]: not found", user.Login)
			return nil, status.Error(codes.NotFound, _msgUserNotFound)
		}

		g.logger.Errorf("login error for [%s]: %v", user.Login, err)
		return nil, status.Error(codes.Internal, _msgUserFailedToLogin)
	}

	if err = checkPassword(user.Password, pwdHash); err != nil {
		g.logger.Errorf("login error for [%s]: wrong password", user.Login)
		return nil, status.Error(codes.PermissionDenied, _msgUserInvalidLoginOrPassword)
	}

	token, err := token.NewJWTForUser(userID, _userTokenExpiration, g.serverSecret)
	if err != nil {
		g.logger.Errorf("failed to create token for user [%s, %d]: %v", user.Login, userID, err)
		return nil, status.Error(codes.Unavailable, _msgUserFailedToCreateToken)
	}

	return &pb.UserAuthToken{Token: token}, nil
}

func (g *GrpcServer) Ping(ctx context.Context, _ *pb.Empty) (*pb.Empty, error) {
	res := g.stg.Ping(ctx)
	if res {
		return &pb.Empty{}, nil
	}
	return nil, status.Error(codes.Internal, _msgPingFailed)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 13)
	return string(bytes), err
}

func checkPassword(password, pwdHash string) error {
	return bcrypt.CompareHashAndPassword([]byte(pwdHash), []byte(password))
}
