package client

import (
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/rivo/tview"
)

const (
	_pageLogin            = "login"
	_pageCreateUser       = "create user"
	_pageUserSecrets      = "user secrets"
	_pageCreateUserSecret = "create user secret" //nolint:gosec // OK
	_pageError            = "error"
)

func (r *Client) createUIApplication() {
	r.uiApp = tview.NewApplication()
	r.uiPages = tview.NewPages()

	r.uiCreateLoginPage()
	r.uiCreateUserPage()
	r.uiCreateUserSecretsPage()
	r.uiCreateUserCreateSecretPage()
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

	flex := uiCenteredWidget(r.wdgLogin, 10, 1)
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

	flex := uiCenteredWidget(r.wdgCreateUser, 10, 1)
	r.uiPages.AddPage(_pageCreateUser, flex, true, false)
}

func (r *Client) uiCreateUserSecretsPage() {
	r.wdgLstSecrets = tview.NewList().ShowSecondaryText(false)
	r.wdgLstSecrets.SetBorder(true).SetTitle("Secrets")

	flexSecrets := tview.NewFlex().SetDirection(tview.FlexRow)
	flexSecrets.AddItem(r.wdgLstSecrets, 0, 1, true)

	formActions := tview.NewForm().
		AddTextView("User", "", 30, 1, false, false).
		AddButton("Create secret", r.uiSwitchToCreateUserSecret).
		AddButton("Reload", r.doReload).
		AddButton("Logout", r.uiSwitchToLogin)
	r.wdgUser, _ = formActions.GetFormItemByLabel("User").(*tview.TextView)

	flexSecrets.AddItem(formActions, 5, 1, false)

	r.uiPages.AddPage(_pageUserSecrets, uiCenteredWidget(flexSecrets, 0, 10), true, false)
}

func (r *Client) uiCreateUserCreateSecretPage() {
	r.wdgCreateUserSecret = tview.NewForm()
	r.wdgCreateUserSecret.
		AddDropDown(
			"Type",
			[]string{
				"",
				model.CredsSecret.String(),
				model.TextSecret.String(),
				model.BinarySecret.String(),
				model.CardSecret.String(),
			},
			0,
			r.uiChangeSecretPayloadFields,
		).
		AddInputField("Name", "", 0, nil, nil).
		AddTextArea("Meta", "", 0, 3, 0, nil).
		AddButton("Create", r.doCreateSecret).
		AddButton("Back to list", r.doReload)
	r.wdgCreateUserSecret.
		SetBorder(true).
		SetTitle("Create secret")

	r.uiPages.AddPage(_pageCreateUserSecret, uiCenteredWidget(r.wdgCreateUserSecret, 0, 10), true, false)
}

func (r *Client) uiCreateErrorPage() {
	r.wdgError = tview.NewForm().
		AddTextView("Message", "", 0, 1, true, true).
		AddButton("Back", nil)
	r.wdgError.
		SetBorder(true).
		SetTitle("Error")

	flex := uiCenteredWidget(r.wdgError, 10, 1)
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

func (r *Client) uiSwitchToCreateUserSecret() {
	r.wdgCreateUserSecret.GetFormItemByLabel("Type").(*tview.DropDown).SetCurrentOption(0)
	r.wdgCreateUserSecret.GetFormItemByLabel("Name").(*tview.InputField).SetText("")
	r.wdgCreateUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).SetText("", true)
	r.uiPages.SwitchToPage(_pageCreateUserSecret)
}

func (r *Client) uiSwitchToError(errMsg string, fnCallback func()) {
	r.wdgError.GetFormItemByLabel("Message").(*tview.TextView).SetText(errMsg)
	r.wdgError.SetFocus(r.wdgError.GetFormItemIndex("Error"))

	btnBackInd := r.wdgError.GetButtonIndex("Back")
	r.wdgError.GetButton(btnBackInd).SetSelectedFunc(fnCallback)
	r.uiPages.SwitchToPage(_pageError)
}

func (r *Client) uiChangeSecretPayloadFields(choosenType string, _ int) {
	metaIndex := r.wdgCreateUserSecret.GetFormItemIndex("Meta")
	i := r.wdgCreateUserSecret.GetFormItemCount() - 1
	for i != metaIndex {
		r.wdgCreateUserSecret.RemoveFormItem(i)
		i--
	}

	switch choosenType {
	case model.CredsSecret.String():
		r.wdgCreateUserSecret.
			AddInputField("Login", "", 0, nil, nil).
			AddInputField("Password", "", 0, nil, nil)
	case model.TextSecret.String():
		r.wdgCreateUserSecret.
			AddTextArea("Text", "", 0, 3, 0, nil)
	case model.BinarySecret.String():
		r.wdgCreateUserSecret.
			AddTextArea("Binary", "", 0, 3, 0, nil)
	case model.CardSecret.String():
		r.wdgCreateUserSecret.
			AddInputField("Card number", "", 0, nil, nil).
			AddInputField("Card holder", "", 0, nil, nil).
			AddInputField("Valid thru", "", 0, nil, nil).
			AddInputField("CVV", "", 0, nil, nil)
	}
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
