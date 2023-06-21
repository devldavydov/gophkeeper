//nolint:gosec // OK
package client

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/rivo/tview"
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

	token, err := r.tr.UserCreate(userLogin, userPassword)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		case errors.Is(err, transport.ErrUserAlreadyExists):
			r.uiSwitchToError(_msgUserAlreadyExists, r.uiSwitchToCreateUser)
		case errors.Is(err, transport.ErrUserInvalidCredentials):
			r.uiSwitchToError(_msgUserInvalidCreds, r.uiSwitchToCreateUser)
		}
		return
	}

	r.cltToken = token
	r.wdgUser.SetText(userLogin)
	r.doReload()
}

func (r *Client) doLogin() {
	userLogin := r.wdgLogin.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.wdgLogin.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	token, err := r.tr.UserLogin(userLogin, userPassword)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		case errors.Is(err, transport.ErrUserNotFound):
			r.uiSwitchToError(_msgUserNotFound, r.uiSwitchToLogin)
		case errors.Is(err, transport.ErrUserLoginFailed):
			r.uiSwitchToError(_msgUserLoginFailed, r.uiSwitchToLogin)
		}
		return
	}

	r.cltToken = token
	r.wdgUser.SetText(userLogin)
	r.doReload()
}

func (r *Client) doReload() {
	r.wdgLstSecrets.Clear()

	// Load secrets from server
	lstSecrets, err := r.tr.SecretGetList(r.cltToken)
	if err != nil {
		r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUser)
		return
	}

	// Store internal
	r.lstSecrets = lstSecrets
	for _, scrt := range r.lstSecrets {
		r.wdgLstSecrets.AddItem(
			fmt.Sprintf("%s (%s)", scrt.Name, scrt.Type),
			"",
			0,
			func() {
				r.doEditSecret(scrt.Name)
			})
	}

	r.uiSwitchToUserSecrets()
}

func (r *Client) doCreateSecret() {
	secret := &model.Secret{
		Name: r.wdgCreateUserSecret.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		Meta: r.wdgCreateUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).GetText(),
	}
	var payload model.Payload

	_, fType := r.wdgCreateUserSecret.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()

	switch fType {
	case model.CredsSecret.String():
		secret.Type = model.CredsSecret
		payload = model.NewCredsPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Login").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Password").(*tview.InputField).GetText(),
		)
	case model.TextSecret.String():
		secret.Type = model.TextSecret
		payload = model.NewTextPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Text").(*tview.TextArea).GetText(),
		)
	case model.BinarySecret.String():
		secret.Type = model.BinarySecret
		binData, _ := hex.DecodeString(r.wdgCreateUserSecret.GetFormItemByLabel("Binary").(*tview.TextArea).GetText())
		payload = model.NewBinaryPayload(binData)
	case model.CardSecret.String():
		secret.Type = model.CardSecret
		payload = model.NewCardPayload(
			r.wdgCreateUserSecret.GetFormItemByLabel("Card number").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Card holder").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("Valid thru").(*tview.InputField).GetText(),
			r.wdgCreateUserSecret.GetFormItemByLabel("CVV").(*tview.InputField).GetText(),
		)
	}

	err := r.tr.SecretCreate(r.cltToken, secret, payload)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrSecretAlreadyExists):
			r.uiSwitchToError(_msgUserSecretAlreadyExists, r.uiSwitchToCreateUserSecret)
		case errors.Is(err, transport.ErrInternalError):
			r.uiSwitchToError(_msgClientError, r.uiSwitchToCreateUserSecret)
		case errors.Is(err, transport.ErrInternalServerError):
			r.uiSwitchToError(_msgInternalServerError, r.uiSwitchToCreateUserSecret)
		}
		return
	}

	r.doReload()
}

func (r *Client) doEditSecret(secretName string) {
	panic(secretName)
}
