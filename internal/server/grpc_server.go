package server

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/cipher"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/grpc/interceptor"
	"github.com/devldavydov/gophkeeper/internal/server/storage"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	MetaUserID = "USER_ID"

	_userTokenExpiration = 24 * time.Hour

	_msgPingFailed = "ping failed"
	//
	_msgUserCredentialsBadRequest  = "invalid credentials" //nolint:gosec // Ok
	_msgUserFailedToCreate         = "failed to create user"
	_msgUserAlreadyExists          = "user already exists"
	_msgUserFailedToCreateToken    = "failed to create user token"
	_msgUserTokenError             = "user token error"
	_msgUserInvalidLoginOrPassword = "user invalid login or password"
	_msgUserNotFound               = "user not found"
	_msgUserFailedToLogin          = "failed to login"
	//
	_msgSecretsNotFound      = "secrets not found" //nolint:gosec // Ok
	_msgSecretsFailedToGet   = "failed to get secrets"
	_msgSecretBadRequest     = "invalid secret"
	_msgSecretFailedToCreate = "failed to create secret"
	_msgSecretAlreadyExists  = "secret already exists"
	_msgSecretNotFound       = "secret not found"
	_msgSecretFailedToDelete = "failed to delete secret"
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

func (g *GrpcServer) SecretGetList(ctx context.Context, _ *pb.Empty) (*pb.SecretGetListResponse, error) {
	userID, err := g.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dbSecretList, err := g.stg.GetAllSecrets(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrNoSecrets) {
			return nil, status.Error(codes.NotFound, _msgSecretsNotFound)
		}
		g.logger.Errorf("[user=%d] secret get list error %v", userID, err)
		return nil, status.Error(codes.Internal, _msgSecretsFailedToGet)
	}

	respList := &pb.SecretGetListResponse{Items: make([]*pb.SecretListItem, 0, len(dbSecretList))}
	for _, dbItem := range dbSecretList {
		respList.Items = append(respList.Items, &pb.SecretListItem{
			Name:    dbItem.Name,
			Type:    pb.SecretType(dbItem.Type),
			Version: dbItem.Version,
		})
	}

	return respList, nil
}

func (g *GrpcServer) SecretGet(ctx context.Context, in *pb.SecretGetRequest) (*pb.Secret, error) {
	userID, err := g.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Get from storage
	dbSecret, err := g.stg.GetSecret(ctx, userID, in.Name)
	if err != nil {
		if errors.Is(err, storage.ErrSecretNotFound) {
			g.logger.Errorf("[user=%d] secret [%s] get error: not found", userID, in.Name)
			return nil, status.Error(codes.NotFound, _msgSecretNotFound)
		}
		g.logger.Errorf("[user=%d] secret [%s] get error: %v", userID, in.Name, err)
		return nil, status.Error(codes.Internal, _msgSecretsFailedToGet)
	}

	// Decode payload
	payloadRaw, err := cipher.AESDecrypt(dbSecret.PayloadRaw, g.serverSecret)
	if err != nil {
		g.logger.Errorf("[user=%d] secret [%s] get error: decrypt payload error: %v", userID, in.Name, err)
		return nil, status.Error(codes.Internal, _msgSecretsFailedToGet)
	}

	secret := &pb.Secret{
		Name:       dbSecret.Name,
		Type:       pb.SecretType(dbSecret.Type),
		Version:    dbSecret.Version,
		Meta:       dbSecret.Meta,
		PayloadRaw: payloadRaw,
	}

	return secret, nil
}

func (g *GrpcServer) SecretCreate(ctx context.Context, in *pb.SecretCreateRequest) (*pb.Empty, error) {
	userID, err := g.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Check secret
	if err = model.ValidSecretType(model.SecretType(in.Secret.Type)); err != nil {
		g.logger.Errorf("[user=%d] secret create error: invalid type [%d]", userID, in.Secret.Type)
		return nil, status.Error(codes.InvalidArgument, _msgSecretBadRequest)
	}

	if in.Secret.Name == "" {
		g.logger.Errorf("[user=%d] secret create error: empty name", userID)
		return nil, status.Error(codes.InvalidArgument, _msgSecretBadRequest)
	}

	// Encrypt payload
	encPayload, err := cipher.AESEncrpyt(in.Secret.PayloadRaw, g.serverSecret)
	if err != nil {
		g.logger.Errorf("[user=%d] secret [%s] create error: encrypt payload error: %v", userID, in.Secret.Name, err)
		return nil, status.Error(codes.Internal, _msgSecretFailedToCreate)
	}

	// Create secret
	secret := &model.Secret{
		Type:       model.SecretType(in.Secret.Type),
		Name:       in.Secret.Name,
		Meta:       in.Secret.Meta,
		Version:    in.Secret.Version,
		PayloadRaw: encPayload,
	}

	err = g.stg.CreateSecret(ctx, userID, secret)
	if err != nil {
		if errors.Is(err, storage.ErrSecretAlreadyExists) {
			g.logger.Errorf("[user=%d] secret [%s] create error: already exists", userID, in.Secret.Name)
			return nil, status.Error(codes.AlreadyExists, _msgSecretFailedToCreate)
		}

		g.logger.Errorf("[user=%d] secret [%s] create error: %v", userID, in.Secret.Name, err)
		return nil, status.Error(codes.Internal, _msgSecretFailedToCreate)
	}

	return &pb.Empty{}, nil
}

func (g *GrpcServer) SecretUpdate(ctx context.Context, in *pb.SecretUpdateRequest) (*pb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SecretUpdate not implemented")
}
func (g *GrpcServer) SecretDelete(ctx context.Context, in *pb.SecretDeleteRequest) (*pb.Empty, error) {
	userID, err := g.getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	err = g.stg.DeleteSecret(ctx, userID, in.Name)
	if err != nil {
		g.logger.Errorf("[user=%d] secret [%s] delete error: %v", userID, in.Name, err)
		return nil, status.Error(codes.Internal, _msgSecretFailedToDelete)
	}

	return &pb.Empty{}, nil
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

func (g *GrpcServer) getUserIDFromContext(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		g.logger.Error("user token error: no metadata")
		return 0, status.Error(codes.Internal, _msgUserTokenError)
	}

	sUserID := md.Get(interceptor.MetaUserID)
	if len(sUserID) != 1 {
		g.logger.Error("user token error: wrong metadata length")
		return 0, status.Error(codes.Internal, _msgUserTokenError)
	}

	userID, err := strconv.ParseInt(sUserID[0], 10, 64)
	if err != nil {
		g.logger.Errorf("user token error: %v", err)
		return 0, status.Error(codes.Internal, _msgUserTokenError)
	}

	return userID, nil
}
