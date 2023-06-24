//nolint:dupl // OK
package ui

import (
	"errors"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/rivo/tview"
)

func (r *App) createLoginPage() {
	r.frmLogin = tview.NewForm().
		AddInputField("Login", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddButton("Login", r.doLogin).
		AddButton("Create user", r.showCreateUser).
		AddButton("Exit", r.app.Stop)
	r.frmLogin.
		SetBorder(true).
		SetTitle("Login")

	flex := uiCenteredWidget(r.frmLogin, 10, 1)
	r.uiPages.AddPage(_pageLogin, flex, true, true)
}

func (r *App) showLogin() {
	r.frmLogin.GetFormItemByLabel("Login").(*tview.InputField).SetText("")
	r.frmLogin.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	r.frmLogin.SetFocus(r.frmLogin.GetFormItemIndex("Login"))
	r.uiPages.SwitchToPage(_pageLogin)
}

func (r *App) doLogin() {
	userLogin := r.frmLogin.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.frmLogin.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	token, err := r.tr.UserLogin(userLogin, userPassword)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.showError(_msgInternalServerError, r.showCreateUser)
		case errors.Is(err, transport.ErrUserNotFound):
			r.showError(_msgUserNotFound, r.showLogin)
		case errors.Is(err, transport.ErrUserPermissionDenied):
			r.showError(_msgUserLoginFailed, r.showLogin)
		}
		return
	}

	r.cltToken = token
	r.wdgUser.SetText(userLogin)
	r.doReloadUserSecrets()
}
