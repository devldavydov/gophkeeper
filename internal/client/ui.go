package client

import (
	"github.com/rivo/tview"
)

const (
	_pageLogin       = "login"
	_pageCreateUser  = "create user"
	_pageUserSecrets = "user secrets"
	_pageError       = "error"

	_msgInternalServerError = "Internal server error"
	_msgUserAlreadyExists   = "User already exists"
	_msgUserInvalidCreds    = "User invalid credentials" //nolint:gosec // OK
	_msgUserNotFound        = "User not found"
	_msgUserLoginFailed     = "User with login/password failed"
)

func (r *Client) createUIApplication() {
	r.uiApp = tview.NewApplication()
	r.uiPages = tview.NewPages()

	r.uiCreateLoginPage()
	r.uiCreateUserPage()
	r.uiCreateUserSecretsPage()
	r.uiCreateErrorPage()

	r.uiApp.
		SetRoot(r.uiPages, true).
		SetFocus(r.uiPages).
		EnableMouse(true)
}

func (r *Client) uiCreateLoginPage() {
	r.wdgLogin = tview.NewForm().
		AddInputField("Login", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddButton("Login", r.doLogin).
		AddButton("Create user", r.uiSwitchToCreateUser).
		AddButton("Exit", r.uiApp.Stop)
	r.wdgLogin.
		SetBorder(true).
		SetTitle("GophKeeper - Login")
	r.uiPages.AddPage(_pageLogin, r.wdgLogin, true, true)
}

func (r *Client) uiCreateUserPage() {
	r.wdgCreateUser = tview.NewForm().
		AddInputField("Login", "", 30, nil, nil).
		AddPasswordField("Password", "", 30, '*', nil).
		AddButton("Create", r.doCreateUser).
		AddButton("Back to Login", r.uiSwitchToLogin).
		AddButton("Exit", r.uiApp.Stop)
	r.wdgCreateUser.
		SetBorder(true).
		SetTitle("GophKeeper - Create new user")
	r.uiPages.AddPage(_pageCreateUser, r.wdgCreateUser, true, false)
}

func (r *Client) uiCreateUserSecretsPage() {
	r.wdgUserSecrets = tview.NewForm().
		AddTextView("Token", "", 100, 20, true, true).
		AddButton("Create secret", func() {

		}).
		AddButton("Logout", r.uiSwitchToLogin)
	r.wdgUserSecrets.
		SetBorder(true).
		SetTitle("GophKeeper - Secrets list")
	r.uiPages.AddPage(_pageUserSecrets, r.wdgUserSecrets, true, false)
}

func (r *Client) uiCreateErrorPage() {
	r.wdgError = tview.NewForm().
		AddTextView("Error", "", 50, 5, true, true).
		AddButton("Back", nil)
	r.wdgUserSecrets.
		SetBorder(true).
		SetTitle("GophKeeper - Error")
	r.uiPages.AddPage(_pageError, r.wdgError, true, false)
}

func (r *Client) uiSwitchToLogin() {
	r.wdgLogin.GetFormItemByLabel("Login").(*tview.InputField).SetText("")
	r.wdgLogin.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	r.wdgLogin.SetFocus(r.wdgLogin.GetFormItemIndex("Login"))
	r.uiPages.SwitchToPage(_pageLogin)
}

func (r *Client) uiSwitchToCreateUser() {
	r.wdgCreateUser.GetFormItemByLabel("Login").(*tview.InputField).SetText("")
	r.wdgCreateUser.GetFormItemByLabel("Password").(*tview.InputField).SetText("")
	r.wdgCreateUser.SetFocus(r.wdgCreateUser.GetFormItemIndex("Login"))
	r.uiPages.SwitchToPage(_pageCreateUser)
}

func (r *Client) uiSwitchToUserSecrets() {
	r.wdgUserSecrets.GetFormItemByLabel("Token").(*tview.TextView).SetText(r.cltToken)
	r.uiPages.SwitchToPage(_pageUserSecrets)
}

func (r *Client) uiSwitchToError(errMsg string, fnCallback func()) {
	r.wdgError.GetFormItemByLabel("Error").(*tview.TextView).SetText(errMsg)
	r.wdgError.SetFocus(r.wdgError.GetFormItemIndex("Error"))

	btnBackInd := r.wdgError.GetButtonIndex("Back")
	r.wdgError.GetButton(btnBackInd).SetSelectedFunc(fnCallback)
	r.uiPages.SwitchToPage(_pageError)
}
