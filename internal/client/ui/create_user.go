package ui

import (
	"errors"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/rivo/tview"
)

func (r *UIApp) createUserPage() {
	r.frmCreateUser = tview.NewForm().
		AddInputField("Login", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddButton("Create", r.doCreateUser).
		AddButton("Back to Login", r.showLogin).
		AddButton("Exit", r.app.Stop)
	r.frmCreateUser.
		SetBorder(true).
		SetTitle("Create new user")

	flex := uiCenteredWidget(r.frmCreateUser, 10, 1)
	r.uiPages.AddPage(_pageCreateUser, flex, true, false)
}

func (r *UIApp) showCreateUser() {
	r.frmCreateUser.GetFormItemByLabel("Login").(*tview.InputField).SetText("")
	r.frmCreateUser.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	r.frmCreateUser.SetFocus(r.frmCreateUser.GetFormItemIndex("Login"))
	r.uiPages.SwitchToPage(_pageCreateUser)
}

func (r *UIApp) doCreateUser() {
	userLogin := r.frmCreateUser.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.frmCreateUser.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	token, err := r.tr.UserCreate(userLogin, userPassword)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.showError(_msgInternalServerError, r.showCreateUser)
		case errors.Is(err, transport.ErrUserAlreadyExists):
			r.showError(_msgUserAlreadyExists, r.showCreateUser)
		case errors.Is(err, transport.ErrUserInvalidCredentials):
			r.showError(_msgUserInvalidCreds, r.showCreateUser)
		}
		return
	}

	r.cltToken = token
	r.wdgUser.SetText(userLogin)
	r.doReloadUserSecrets()
}
