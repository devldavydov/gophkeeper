package transport

import (
	"context"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/sirupsen/logrus"

	"github.com/devldavydov/gophkeeper/internal/common/nettools"
	gkTLS "github.com/devldavydov/gophkeeper/internal/common/tls"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const _serverRequestTimeout = 15 * time.Second

// GrpcTransport is a gRPC implementation of Transport interface.
type GrpcTransport struct {
	gClt   pb.GophKeeperServiceClient
	logger *logrus.Logger
}

var _ Transport = (*GrpcTransport)(nil)

// NewGrpcTransport creates new GrpcTransport object.
func NewGrpcTransport(
	serverAddress *nettools.Address,
	tlsCACertPath string,
	logger *logrus.Logger,
) (*GrpcTransport, error) {
	tlsCredentials, err := gkTLS.LoadCACert(tlsCACertPath, "")
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(
		serverAddress.String(),
		grpc.WithTransportCredentials(tlsCredentials),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(model.MaxPayloadSizeBytes+1024),
			grpc.MaxCallSendMsgSize(model.MaxPayloadSizeBytes+1024),
		),
	)
	if err != nil {
		return nil, err
	}

	return newGrpcTransport(pb.NewGophKeeperServiceClient(conn), logger), nil
}

func newGrpcTransport(clt pb.GophKeeperServiceClient, logger *logrus.Logger) *GrpcTransport {
	return &GrpcTransport{gClt: clt, logger: logger}
}

// UserCreate is a gRPC implemention of user creation method. Accepts user login and password.
//
// Returns user token or error:
//
// - ErrUserAlreadyExists - user with given login/password already exists on server.
//
// - ErrUserInvalidCredentials - provided user login or password invalid.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) UserCreate(userLogin, userPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pbToken, err := gt.gClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
		gt.logger.Errorf("gRPC user create request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return "", ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.AlreadyExists:
			return "", ErrUserAlreadyExists
		case codes.InvalidArgument:
			return "", ErrUserInvalidCredentials
		default:
			return "", ErrInternalServerError
		}
	}

	return pbToken.Token, nil
}

// UserLogin is a gRPC implemention of user login method. Accepts user login and password.
//
// Returns user token or error:
//
// - ErrUserNotFound - user with given login/password not found on server.
//
// - ErrUserPermissionDenied - provided user login or password not match and permission denied.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) UserLogin(userLogin, userPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	pbToken, err := gt.gClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
		gt.logger.Errorf("gRPC user login request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return "", ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return "", ErrUserNotFound
		case codes.PermissionDenied:
			return "", ErrUserPermissionDenied
		default:
			return "", ErrInternalServerError
		}
	}

	return pbToken.Token, nil
}

// SecretGetList is a gRPC implemention of user's secret list method. Accepts authenticated user token.
//
// Returns list of secrets or error:
//
// - ErrUserPermissionDenied - provided token not valid and permission denied.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) SecretGetList(token string) ([]model.SecretInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	lstSrvSecrets, err := gt.gClt.SecretGetList(contextWithToken(ctx, token), &pb.Empty{})
	if err != nil {
		gt.logger.Errorf("gRPC secret get list request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return nil, ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return make([]model.SecretInfo, 0), nil
		case codes.PermissionDenied:
			return nil, ErrUserPermissionDenied
		default:
			return nil, ErrInternalServerError
		}
	}

	lstSecretInfo := make([]model.SecretInfo, 0, len(lstSrvSecrets.Items))
	for _, srvSecret := range lstSrvSecrets.Items {
		secretInfo := model.SecretInfo{
			Name:    srvSecret.Name,
			Version: srvSecret.Version,
			Type:    model.SecretType(srvSecret.Type),
		}
		lstSecretInfo = append(lstSecretInfo, secretInfo)
	}

	return lstSecretInfo, nil
}

// SecretGet is a gRPC implemention of user's get secret method. Accepts authenticated user token and secret name.
//
// Returns secret or error:
//
// - ErrSecretNotFound - secret with given name not found.
//
// - ErrUserPermissionDenied - provided user token not valid and permission denied.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) SecretGet(token, name string) (*model.Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	pbSecret, err := gt.gClt.SecretGet(contextWithToken(ctx, token), &pb.SecretGetRequest{Name: name})
	if err != nil {
		gt.logger.Errorf("gRPC secret get request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return nil, ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return nil, ErrSecretNotFound
		case codes.PermissionDenied:
			return nil, ErrUserPermissionDenied
		default:
			return nil, ErrInternalServerError
		}
	}

	return &model.Secret{
		Type:       model.SecretType(pbSecret.Type),
		Name:       pbSecret.Name,
		Meta:       pbSecret.Meta,
		Version:    pbSecret.Version,
		PayloadRaw: pbSecret.PayloadRaw,
	}, nil
}

// SecretCreate is a gRPC implemention of user's create secret method. Accepts authenticated user token and new secret.
//
// Returns nil or error:
//
// - ErrSecretAlreadyExists - secret with given name already exists.
//
// - ErrUserPermissionDenied - provided user token not valid and permission denied.
//
// - ErrSecretPayloadSizeExceeded - provided file in binary secret too big.
//
// - ErrSecretInvalid - provided secret not valid.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) SecretCreate(token string, secret *model.Secret) error {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	secretReq := &pb.SecretCreateRequest{
		Secret: &pb.Secret{
			Name:       secret.Name,
			Meta:       secret.Meta,
			Type:       pb.SecretType(secret.Type),
			PayloadRaw: secret.PayloadRaw,
			Version:    1,
		},
	}

	_, err := gt.gClt.SecretCreate(contextWithToken(ctx, token), secretReq)
	if err != nil {
		gt.logger.Errorf("gRPC secret create request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.AlreadyExists:
			return ErrSecretAlreadyExists
		case codes.ResourceExhausted:
			return ErrSecretPayloadSizeExceeded
		case codes.PermissionDenied:
			return ErrUserPermissionDenied
		case codes.InvalidArgument:
			return ErrSecretInvalid
		default:
			return ErrInternalServerError
		}
	}

	return nil
}

// SecretUpdate is a gRPC implemention of user's update secret method.
//
// Accepts authenticated user token and secret update.
//
// Returns nil or error:
//
// - ErrSecretNotFound - secret with given name not found.
//
// - ErrUserPermissionDenied - provided user token not valid and permission denied.
//
// - ErrSecretPayloadSizeExceeded - provided file in binary secret too big.
//
// - ErrSecretInvalid - provided secret not valid.
//
// - ErrSecretOutdated - secret outdated (was changed in another session).
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) SecretUpdate(token, name string, updSecret *model.SecretUpdate) error {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	updReq := &pb.SecretUpdateRequest{
		Name:          name,
		Meta:          updSecret.Meta,
		Version:       updSecret.Version,
		PayloadRaw:    updSecret.PayloadRaw,
		UpdatePayload: updSecret.UpdatePayload,
	}

	_, err := gt.gClt.SecretUpdate(contextWithToken(ctx, token), updReq)
	if err != nil {
		gt.logger.Errorf("gRPC secret update request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return ErrSecretNotFound
		case codes.FailedPrecondition:
			return ErrSecretOutdated
		case codes.InvalidArgument:
			return ErrSecretInvalid
		case codes.PermissionDenied:
			return ErrUserPermissionDenied
		case codes.ResourceExhausted:
			return ErrSecretPayloadSizeExceeded
		default:
			return ErrInternalServerError
		}
	}

	return nil
}

// SecretDelete is a gRPC implemention of user's delete secret method.
//
// Accepts authenticated user token and secret name.
//
// Returns nil or error:
//
// - ErrUserPermissionDenied - provided user token not valid and permission denied.
//
// - ErrInternalServerError - unexpected server error.
func (gt *GrpcTransport) SecretDelete(token, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	_, err := gt.gClt.SecretDelete(contextWithToken(ctx, token), &pb.SecretDeleteRequest{Name: name})
	if err != nil {
		gt.logger.Errorf("gRPC secret delete request error: %v", err)
		status, ok := status.FromError(err)
		if !ok {
			return ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.PermissionDenied:
			return ErrUserPermissionDenied
		default:
			return ErrInternalServerError
		}
	}

	return nil
}

func contextWithToken(ctx context.Context, cltToken string) context.Context {
	md := metadata.New(map[string]string{token.HeaderName: cltToken})
	return metadata.NewOutgoingContext(ctx, md)
}
