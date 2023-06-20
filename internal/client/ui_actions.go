//nolint:gosec // OK
package client

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/devldavydov/gophkeeper/internal/common/model"
	gkMsgp "github.com/devldavydov/gophkeeper/internal/common/msgp"
	"github.com/devldavydov/gophkeeper/internal/common/token"
	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/rivo/tview"
	"github.com/tinylib/msgp/msgp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	_msgInternalServerError     = "Internal server error"
	_msgClientError             = "Internal client error"
	_msgUserAlreadyExists       = "User already exists"
	_msgUserInvalidCreds        = "User invalid credentials"
	_msgUserNotFound            = "User not found"
	_msgUserLoginFailed         = "User wrong login/password"
	_msgUserSecretAlreadyExists = "User secret already exists"
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

	r.wdgLstSecrets.Clear()

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
	r.lstSecrets = make([]model.SecretInfo, 0)
	for _, srvSecret := range lstSrvSecrets.Items {
		secretInfo := &model.SecretInfo{
			Name:    srvSecret.Name,
			Version: srvSecret.Version,
			Type:    model.SecretType(srvSecret.Type),
		}
		r.lstSecrets = append(r.lstSecrets, *secretInfo)

		r.wdgLstSecrets.AddItem(
			fmt.Sprintf("%s (%s)", secretInfo.Name, secretInfo.Type),
			"",
			0,
			func() {
				r.doEditSecret(secretInfo.Name)
			})
	}

	r.uiSwitchToUserSecrets()
}

func (r *Client) doCreateSecret() {
	ctx, cancel := context.WithTimeout(context.Background(), _serverRequestTimeout)
	defer cancel()

	fName := r.wdgCreateUserSecret.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	fMeta := r.wdgCreateUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).GetText()
	_, fType := r.wdgCreateUserSecret.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()

	var sType pb.SecretType
	var payload msgp.Encodable

	switch fType {
	case model.CredsSecret.String():
		sType = pb.SecretType_CREDS
		payload = model.NewCredsPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Login").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Password").(*tview.InputField).GetText(),
		)
	case model.TextSecret.String():
		sType = pb.SecretType_TEXT
		payload = model.NewTextPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Text").(*tview.TextArea).GetText(),
		)
	case model.BinarySecret.String():
		sType = pb.SecretType_BINARY
		binData, _ := hex.DecodeString(r.wdgCreateUserSecret.GetFormItemByLabel("Binary").(*tview.TextArea).GetText())
		payload = model.NewBinaryPayload(binData)
	case model.CardSecret.String():
		sType = pb.SecretType_CARD
		payload = model.NewCardPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Card number").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Card holder").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Valid thru").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("CVV").(*tview.InputField).GetText(),
		)
	}

	payloadRaw, err := gkMsgp.Serialize(payload)
	if err != nil {
		r.uiSwitchToError(_msgClientError, r.uiSwitchToCreateUser)
	}

	secretReq := &pb.SecretCreateRequest{
		Secret: &pb.Secret{
			Name:       fName,
			Meta:       fMeta,
			Type:       sType,
			PayloadRaw: payloadRaw,
			Version:    time.Now().UTC().Unix(),
		},
	}

	_, err = r.gClt.SecretCreate(r.contextWithToken(ctx), secretReq)
	if err != nil {
		status, ok := status.FromError(err)
		if !ok {
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		}

		switch status.Code() { //nolint:exhaustive // OK
		case codes.AlreadyExists:
			r.uiSwitchToError(_msgUserSecretAlreadyExists, r.uiSwitchToCreateUserSecret)
		default:
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUserSecret)
		}
		return
	}

	r.doReload()
}

func (r *Client) doEditSecret(secretName string) {
	panic(secretName)
}

func (r *Client) contextWithToken(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{token.HeaderName: r.cltToken})
	return metadata.NewOutgoingContext(ctx, md)
}
