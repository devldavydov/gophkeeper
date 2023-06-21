package ui

import (
	"fmt"

	"github.com/rivo/tview"
)

func (r *UIApp) createUserSecretsPage() {
	r.wdgLstSecrets = tview.NewList().ShowSecondaryText(false)
	r.wdgLstSecrets.SetBorder(true).SetTitle("Secrets")

	flexSecrets := tview.NewFlex().SetDirection(tview.FlexRow)
	flexSecrets.AddItem(r.wdgLstSecrets, 0, 1, true)

	formActions := tview.NewForm().
		AddTextView("User", "", 30, 1, false, false).
		AddButton("Create secret", r.showCreateUserSecret).
		AddButton("Reload", r.doReloadUserSecrets).
		AddButton("Logout", r.showLogin)
	r.wdgUser, _ = formActions.GetFormItemByLabel("User").(*tview.TextView)

	flexSecrets.AddItem(formActions, 5, 1, false)

	r.uiPages.AddPage(_pageUserSecrets, uiCenteredWidget(flexSecrets, 0, 10), true, false)
}

func (r *UIApp) doReloadUserSecrets() {
	r.wdgLstSecrets.Clear()

	// Load secrets from server
	lstSecrets, err := r.tr.SecretGetList(r.cltToken)
	if err != nil {
		r.showError(_msgInternalServerError, r.showCreateUser)
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
				r.doEditUserSecret(scrt.Name)
			})
	}

	r.uiPages.SwitchToPage(_pageUserSecrets)
}
