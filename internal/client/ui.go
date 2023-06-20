package client

import (
	"context"
	"time"

	pb "github.com/devldavydov/gophkeeper/internal/grpc"
	"github.com/rivo/tview"
)

const (
	_pageLogin       = "login"
	_pageCreateUser  = "create user"
	_pageUserSecrets = "user secrets"
)

func (r *Client) createUIApplication() {
	r.uiApp = tview.NewApplication()
	r.uiPages = tview.NewPages()

	r.uiCreateLoginPage()
	r.uiCreateUserPage()
	r.uiCreateUserSecretsPage()

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

func (r *Client) doCreateUser() {
	userLogin := r.wdgCreateUser.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.wdgCreateUser.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pbToken, _ := r.gClt.UserCreate(ctx, &pb.User{Login: userLogin, Password: userPassword})
	r.cltToken = pbToken.Token
	r.uiSwitchToUserSecrets()
}

func (r *Client) doLogin() {
	userLogin := r.wdgLogin.GetFormItemByLabel("Login").(*tview.InputField).GetText()
	userPassword := r.wdgLogin.GetFormItemByLabel("Password").(*tview.InputField).GetText()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	pbToken, _ := r.gClt.UserLogin(ctx, &pb.User{Login: userLogin, Password: userPassword})
	r.cltToken = pbToken.Token
	r.uiSwitchToUserSecrets()
}
