package transport

import (
	"context"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"

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

// GrpcTransport is a gRPC implementation for Transport interface.
type GrpcTransport struct {
	gClt pb.GophKeeperServiceClient
}

var _ Transport = (*GrpcTransport)(nil)

func NewGrpcTransport(serverAddress *nettools.Address, tlsCACertPath string) (*GrpcTransport, error) {
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

	return newGrpcTransport(pb.NewGophKeeperServiceClient(conn)), nil
}

func newGrpcTransport(clt pb.GophKeeperServiceClient) *GrpcTransport {
	return &GrpcTransport{gClt: clt}
}

func (gt *GrpcTransport) UserCreate(userLogin, userPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pbToken, err := gt.gClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
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

func (gt *GrpcTransport) UserLogin(userLogin, userPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	pbToken, err := gt.gClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			return "", ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return "", ErrUserNotFound
		case codes.PermissionDenied:
			return "", ErrUserLoginFailed
		default:
			return "", ErrInternalServerError
		}
	}

	return pbToken.Token, nil
}

func (gt *GrpcTransport) SecretGetList(token string) ([]model.SecretInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	lstSrvSecrets, err := gt.gClt.SecretGetList(contextWithToken(ctx, token), &pb.Empty{})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			return nil, ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return make([]model.SecretInfo, 0), nil
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

func (gt *GrpcTransport) SecretGet(token, name string) (*model.Secret, error) {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	pbSecret, err := gt.gClt.SecretGet(contextWithToken(ctx, token), &pb.SecretGetRequest{Name: name})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			return nil, ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return nil, ErrSecretNotFound
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
		status, ok := status.FromError(err)
		if !ok {
			return ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.AlreadyExists:
			return ErrSecretAlreadyExists
		case codes.ResourceExhausted:
			return ErrSecretPayloadSizeExceeded
		case codes.InvalidArgument:
			return ErrSecretInvalid
		default:
			return ErrInternalServerError
		}
	}

	return nil
}

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
		status, ok := status.FromError(err)
		if !ok {
			return ErrInternalServerError
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			return ErrSecretNotFound
		case codes.FailedPrecondition:
			return ErrSecretOutdated
		default:
			return ErrInternalServerError
		}
	}

	return nil
}

func (gt *GrpcTransport) SecretDelete(token, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	_, err := gt.gClt.SecretDelete(contextWithToken(ctx, token), &pb.SecretDeleteRequest{Name: name})
	if err != nil {
		return ErrInternalServerError
	}

	return nil
}

func contextWithToken(ctx context.Context, cltToken string) context.Context {
	md := metadata.New(map[string]string{token.HeaderName: cltToken})
	return metadata.NewOutgoingContext(ctx, md)
}
