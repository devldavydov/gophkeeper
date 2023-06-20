package client

import (
	"context"
	"fmt"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/rivo/tview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
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
	r.wdgUser.SetText(userLogin)
	r.doReload()
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
	r.wdgUser.SetText(userLogin)
	r.doReload()
}

func (r *Client) doReload() {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	// Load secrets from server
	lstSrvSecrets, err := r.gClt.SecretGetList(r.contextWithToken(ctx), &pb.Empty{})
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.NotFound:
			r.uiSwitchToUserSecrets()
		default:
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToLogin)
		}
		return
	}

	// Store internal
	r.wdgLstSecrets.Clear()
	r.lstSecrets = make([]model.SecretInfo, 0)
	for _, srvSecret := range lstSrvSecrets.Items {
		secretInfo := &model.SecretInfo{
			Name:    srvSecret.Name,
			Version: srvSecret.Version,
			Type:    model.SecretType(srvSecret.Type),
		}
		r.lstSecrets = append(r.lstSecrets, *secretInfo)

		r.wdgLstSecrets.AddItem(fmt.Sprintf("%s (%s)", secretInfo.Name, secretInfo.Type), "", 0, nil)
	}

	r.uiSwitchToUserSecrets()
}

func (r *Client) contextWithToken(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{token.HeaderName: r.cltToken})
	return metadata.NewOutgoingContext(ctx, md)
}
