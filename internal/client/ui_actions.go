package client

import (
	"context"
	"time"

	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/rivo/tview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (r *Client) doCreateUser() {
	userLogin := r.wdgCreateUser.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.wdgCreateUser.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pbToken, err := r.gClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.AlreadyExists:
			r.uiSwitchToError(_msgUserAlreadyExists, r.uiSwitchToCreateUser)
		case codes.InvalidArgument:
			r.uiSwitchToError(_msgUserInvalidCreds, r.uiSwitchToCreateUser)
		default:
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		}
		return
	}

	r.cltToken = pbToken.Token
	r.uiSwitchToUserSecrets()
}

func (r *Client) doLogin() {
	userLogin := r.wdgLogin.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.wdgLogin.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	pbToken, err := r.gClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			r.uiSwitchToError(_msgUserNotFound, r.uiSwitchToLogin)
		case codes.PermissionDenied:
			r.uiSwitchToError(_msgUserLoginFailed, r.uiSwitchToLogin)
		default:
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToLogin)
		}
		return
	}

	r.cltToken = pbToken.Token
	r.uiSwitchToUserSecrets()
}
