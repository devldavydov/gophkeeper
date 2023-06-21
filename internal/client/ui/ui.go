// Package ui contains user intreface functionality.
//
//nolint:gosec // OK
package ui

import (
	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/rivo/tview"
)

const (
	_pageLogin            = "login"
	_pageCreateUser       = "create user"
	_pageUserSecrets      = "user secrets"
	_pageCreateUserSecret = "create user secret"
	_pageEditUserSecret   = "edit user secret"
	_pageError            = "error"

	_msgInternalServerError = "Internal server error"
	_msgClientError         = "Internal client error"
	_msgUserAlreadyExists   = "User already exists"
	_msgUserInvalidCreds    = "User invalid credentials"
	_msgUserNotFound        = "User not found"
	_msgUserLoginFailed     = "User wrong login/password"
	_msgSecretAlreadyExists = "Secret already exists"
	_msgSecretNotFound      = "Secret not found. Maybe it was removed in another session."
	_msgSecretOutdated      = "Secret outdated. It was changed in another session."
)

// App represents user interface application.
type App struct {
	cltToken   string
	tr         transport.Transport
	lstSecrets []model.SecretInfo
	//
	app                 *tview.Application
	uiPages             *tview.Pages
	frmLogin            *tview.Form
	frmCreateUser       *tview.Form
	frmCreateUserSecret *tview.Form
	frmEditUserSecret   *tview.Form
	frmError            *tview.Form
	wdgLstSecrets       *tview.List
	wdgUser             *tview.TextView
}

func NewApp(tr transport.Transport) *App {
	return &App{tr: tr, app: tview.NewApplication()}
}

func (r *App) Run() error {
	r.uiPages = tview.NewPages()

	r.createLoginPage()
	r.createUserPage()
	r.createUserSecretsPage()
	r.createCreateUserSecretPage()
	r.createEditUserSecretPage()
	r.createErrorPage()

	r.app.
		SetRoot(r.uiPages, true).
		SetFocus(r.uiPages).
		EnableMouse(true)

	return r.app.Run()
}

func uiCenteredWidget(wdg tview.Primitive, wdgFixed, wdgPropotion int) *tview.Flex {
	flexCommon := tview.NewFlex()
	flexCommon.SetBorder(true).SetTitle("GophKeeper")

	flexWdg := tview.NewFlex().SetDirection(tview.FlexRow)
	flexWdg.AddItem(nil, 0, 1, false)
	flexWdg.AddItem(wdg, wdgFixed, wdgPropotion, true)
	flexWdg.AddItem(nil, 0, 1, false)

	flexCommon.AddItem(nil, 0, 1, false)
	flexCommon.AddItem(flexWdg, 0, 1, true)
	flexCommon.AddItem(nil, 0, 1, false)

	return flexCommon
}
