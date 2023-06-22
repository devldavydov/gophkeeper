package transport

import (
	"errors"
	"fmt"
	"testing"

	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/devldavydov/gophkeeper/internal/grpc/mocks"
	"github.com/golang/mock/gomock"
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
	gt.gmckCtrl = gomock.NewController(gt.T())
	gt.gCltMock = mocks.NewMockGophKeeperServiceClient(gt.gmckCtrl)
	gt.tr = newGrpcTransport(gt.gCltMock)
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

func TestGrpcTransportSuite(t *testing.T) {
	suite.Run(t, new(GrpcTransportSuite))
}
