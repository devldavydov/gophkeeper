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
	_msgUserLoginFailed     = "User wrong login/password"
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
		AddInputField("Login", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddButton("Login", r.doLogin).
		AddButton("Create user", r.uiSwitchToCreateUser).
		AddButton("Exit", r.uiApp.Stop)
	r.wdgLogin.
		SetBorder(true).
		SetTitle("Login")

	flex := uiFlexWithCenteredWidget("GophKeeper", r.wdgLogin)
	r.uiPages.AddPage(_pageLogin, flex, true, true)
}

func (r *Client) uiCreateUserPage() {
	r.wdgCreateUser = tview.NewForm().
		AddInputField("Login", "", 0, nil, nil).
		AddPasswordField("Password", "", 0, '*', nil).
		AddButton("Create", r.doCreateUser).
		AddButton("Back to Login", r.uiSwitchToLogin).
		AddButton("Exit", r.uiApp.Stop)
	r.wdgCreateUser.
		SetBorder(true).
		SetTitle("Create new user")

	flex := uiFlexWithCenteredWidget("GophKeeper", r.wdgCreateUser)
	r.uiPages.AddPage(_pageCreateUser, flex, true, false)
}

func (r *Client) uiCreateUserSecretsPage() {
	flexCommon := tview.NewFlex().SetDirection(tview.FlexRow)
	flexCommon.SetTitle("GophKeeper")
	flexCommon.SetBorder(true)

	r.wdgLstSecrets = tview.NewList().ShowSecondaryText(false)
	r.wdgLstSecrets.SetBorder(true).SetTitle("Secrets")

	flexSecrets := tview.NewFlex()
	flexSecrets.AddItem(r.wdgLstSecrets, 0, 1, false)
	flexSecrets.AddItem(tview.NewBox().SetBorder(true).SetTitle("Secret details"), 0, 4, false)

	formActions := tview.NewForm().
		AddTextView("User", "", 30, 1, false, false).
		AddButton("Create secret", nil).
		AddButton("Reload", r.doReload).
		AddButton("Logout", r.uiSwitchToLogin)
	r.wdgUser, _ = formActions.GetFormItemByLabel("User").(*tview.TextView)

	flexCommon.AddItem(flexSecrets, 0, 1, true)
	flexCommon.AddItem(formActions, 5, 1, true)

	r.uiPages.AddPage(_pageUserSecrets, flexCommon, true, false)
}

func (r *Client) uiCreateErrorPage() {
	r.wdgError = tview.NewForm().
		AddTextView("Message", "", 50, 5, true, true).
		AddButton("Back", nil)
	r.wdgError.
		SetBorder(true).
		SetTitle("Error")

	flex := uiFlexWithCenteredWidget("GophKeeper", r.wdgError)
	r.uiPages.AddPage(_pageError, flex, true, false)
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
	r.uiPages.SwitchToPage(_pageUserSecrets)
}

func (r *Client) uiSwitchToError(errMsg string, fnCallback func()) {
	r.wdgError.GetFormItemByLabel("Message").(*tview.TextView).SetText(errMsg)
	r.wdgError.SetFocus(r.wdgError.GetFormItemIndex("Error"))

	btnBackInd := r.wdgError.GetButtonIndex("Back")
	r.wdgError.GetButton(btnBackInd).SetSelectedFunc(fnCallback)
	r.uiPages.SwitchToPage(_pageError)
}

func uiFlexWithCenteredWidget(title string, wdg tview.Primitive) *tview.Flex {
	flexCommon := tview.NewFlex()
	flexCommon.SetBorder(true).SetTitle(title)

	flexWdg := tview.NewFlex().SetDirection(tview.FlexRow)
	flexWdg.AddItem(nil, 0, 1, false)
	flexWdg.AddItem(wdg, 0, 1, true)
	flexWdg.AddItem(nil, 0, 1, false)

	flexCommon.AddItem(nil, 0, 1, false)
	flexCommon.AddItem(flexWdg, 0, 1, true)
	flexCommon.AddItem(nil, 0, 1, false)

	return flexCommon
}
