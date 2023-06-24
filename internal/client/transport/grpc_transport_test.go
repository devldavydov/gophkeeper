package transport

import (
	"errors"
	"fmt"
	"testing"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/grpc/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcTransportSuite struct {
	suite.Suite
	gmckCtrl *gomock.Controller
	tr       Transport
	gCltMock *mocks.MockGophKeeperServiceClient
}

func (gt *GrpcTransportSuite) SetupTest() {
	logger := logrus.New()
	gt.gmckCtrl = gomock.NewController(gt.T())
	gt.gCltMock = mocks.NewMockGophKeeperServiceClient(gt.gmckCtrl)
	gt.tr = newGrpcTransport(gt.gCltMock, logger)
}

func (gt *GrpcTransportSuite) TearDownTest() {
	gt.gmckCtrl.Finish()
}

func (gt *GrpcTransportSuite) TestUserCreate() {
	fCreate := func() (string, error) {
		return gt.tr.UserCreate("login", "password")
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			UserCreate(gomock.Any(), &pb.User{Login: "login", Password: "password"}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expToken  string
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expToken:  "",
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.AlreadyExists, "")},
			expToken:  "",
			expErr:    ErrUserAlreadyExists,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.InvalidArgument, "")},
			expToken:  "",
			expErr:    ErrUserInvalidCredentials,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.Internal, "")},
			expToken:  "",
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{&pb.UserAuthToken{Token: "foobar"}, nil},
			expToken:  "foobar",
			expErr:    nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			token, err := fCreate()
			gt.Equal(tt.expToken, token)
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestUserLogin() {
	fCreate := func() (string, error) {
		return gt.tr.UserLogin("login", "password")
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			UserLogin(gomock.Any(), &pb.User{Login: "login", Password: "password"}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expToken  string
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expToken:  "",
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.NotFound, "")},
			expToken:  "",
			expErr:    ErrUserNotFound,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expToken:  "",
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.Internal, "")},
			expToken:  "",
			expErr:    ErrInternalServerError,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			token, err := fCreate()
			gt.Equal(tt.expToken, token)
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestSecretGetList() {
	fGetList := func() ([]model.SecretInfo, error) {
		return gt.tr.SecretGetList("token")
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			SecretGetList(gomock.Any(), &pb.Empty{}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expList   []model.SecretInfo
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expList:   nil,
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.Internal, "")},
			expList:   nil,
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expList:   nil,
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.NotFound, "")},
			expList:   make([]model.SecretInfo, 0),
			expErr:    nil,
		},
		{
			fMockArgs: []any{
				&pb.SecretGetListResponse{Items: []*pb.SecretListItem{
					{Name: "foo", Version: 1, Type: pb.SecretType_BINARY},
					{Name: "bar", Version: 2, Type: pb.SecretType_TEXT},
					{Name: "fuzz", Version: 3, Type: pb.SecretType_CARD},
					{Name: "buzz", Version: 4, Type: pb.SecretType_CREDS},
				}},
				nil},
			expList: []model.SecretInfo{
				{Type: model.BinarySecret, Name: "foo", Version: 1},
				{Type: model.TextSecret, Name: "bar", Version: 2},
				{Type: model.CardSecret, Name: "fuzz", Version: 3},
				{Type: model.CredsSecret, Name: "buzz", Version: 4},
			},
			expErr: nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			lst, err := fGetList()
			gt.Equal(tt.expList, lst)
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestSecretGet() {
	fGet := func() (*model.Secret, error) {
		return gt.tr.SecretGet("token", "name")
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			SecretGet(gomock.Any(), &pb.SecretGetRequest{Name: "name"}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expItem   *model.Secret
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expItem:   nil,
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.Internal, "")},
			expItem:   nil,
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expItem:   nil,
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.NotFound, "")},
			expItem:   nil,
			expErr:    nil,
		},
		{
			fMockArgs: []any{
				&pb.Secret{
					Name:       "foo",
					Type:       pb.SecretType_CREDS,
					Version:    1,
					Meta:       "meta",
					PayloadRaw: []byte("test"),
				},
				nil},
			expItem: &model.Secret{
				Type:       model.CredsSecret,
				Name:       "foo",
				Meta:       "meta",
				Version:    1,
				PayloadRaw: []byte("test"),
			},
			expErr: nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			lst, err := fGet()
			gt.Equal(tt.expItem, lst)
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestSecretCreate() {
	fCreate := func() error {
		return gt.tr.SecretCreate("token", &model.Secret{
			Type:       model.CredsSecret,
			Name:       "foo",
			Meta:       "meta",
			Version:    1,
			PayloadRaw: []byte("foo"),
		})
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			SecretCreate(gomock.Any(), &pb.SecretCreateRequest{
				Secret: &pb.Secret{
					Name:       "foo",
					Type:       pb.SecretType_CREDS,
					Version:    1,
					Meta:       "meta",
					PayloadRaw: []byte("foo"),
				},
			}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.AlreadyExists, "")},
			expErr:    ErrSecretAlreadyExists,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.ResourceExhausted, "")},
			expErr:    ErrSecretPayloadSizeExceeded,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.InvalidArgument, "")},
			expErr:    ErrSecretInvalid,
		},
		{
			fMockArgs: []any{nil, nil},
			expErr:    nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			err := fCreate()
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestSecretUpdate() {
	fUpdate := func() error {
		return gt.tr.SecretUpdate("token", "name", &model.SecretUpdate{
			Meta:          "meta",
			Version:       2,
			PayloadRaw:    []byte("foo"),
			UpdatePayload: true,
		})
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			SecretUpdate(gomock.Any(), &pb.SecretUpdateRequest{
				Name:          "name",
				Version:       2,
				Meta:          "meta",
				PayloadRaw:    []byte("foo"),
				UpdatePayload: true,
			}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.NotFound, "")},
			expErr:    ErrSecretNotFound,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.FailedPrecondition, "")},
			expErr:    ErrSecretOutdated,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.ResourceExhausted, "")},
			expErr:    ErrSecretPayloadSizeExceeded,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.InvalidArgument, "")},
			expErr:    ErrSecretInvalid,
		},
		{
			fMockArgs: []any{nil, nil},
			expErr:    nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			err := fUpdate()
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func (gt *GrpcTransportSuite) TestSecretDelete() {
	fDelete := func() error {
		return gt.tr.SecretDelete("token", "name")
	}

	fMock := func(args ...any) {
		gt.gCltMock.EXPECT().
			SecretDelete(gomock.Any(), &pb.SecretDeleteRequest{Name: "name"}).
			Return(args...)
	}

	for i, tt := range []struct {
		fMockArgs []any
		expErr    error
	}{
		{
			fMockArgs: []any{nil, errors.New("Not gRPC error")},
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.Internal, "")},
			expErr:    ErrInternalServerError,
		},
		{
			fMockArgs: []any{nil, status.Error(codes.PermissionDenied, "")},
			expErr:    ErrUserPermissionDenied,
		},
		{
			fMockArgs: []any{nil, nil},
			expErr:    nil,
		},
	} {
		tt := tt
		gt.Run(fmt.Sprintf("Run %d", i), func() {
			fMock(tt.fMockArgs...)
			err := fDelete()
			if tt.expErr != nil {
				gt.ErrorIs(err, tt.expErr)
			}
		})
	}
}

func TestGrpcTransportSuite(t *testing.T) {
	suite.Run(t, new(GrpcTransportSuite))
}
