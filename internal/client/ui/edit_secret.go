package ui

import (
	"errors"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	"github.com/rivo/tview"
)

func (r *App) createEditUserSecretPage() {
	r.frmEditUserSecret = tview.NewForm()
	r.frmEditUserSecret.
		AddInputField("Type", "", 0, nil, nil).
		AddInputField("Name", "", 0, nil, nil).
		AddTextArea("Meta", "", 0, 3, 0, nil).
		AddButton("Save", r.doSaveSecret).
		AddButton("Delete", r.doDeleteUserSecret).
		AddButton("Back to list", r.doReloadUserSecrets)
	r.frmCreateUserSecret.
		SetBorder(true).
		SetTitle("Edit secret")

	r.frmEditUserSecret.GetFormItemByLabel("Type").SetDisabled(true)
	r.frmEditUserSecret.GetFormItemByLabel("Name").SetDisabled(true)

	r.uiPages.AddPage(_pageEditUserSecret, uiCenteredWidget(r.frmEditUserSecret, 0, 10), true, false)
}

func (r *App) showEditUserSecretPage(secret *model.Secret) {
	payload, err := secret.GetPayload()
	if err != nil {
		r.showError(_msgClientError, r.doReloadUserSecrets)
	}

	metaIndex := r.frmEditUserSecret.GetFormItemIndex("Meta")
	i := r.frmEditUserSecret.GetFormItemCount() - 1
	for i != metaIndex {
		r.frmEditUserSecret.RemoveFormItem(i)
		i--
	}

	r.frmEditUserSecret.GetFormItemByLabel("Type").(*tview.InputField).SetText(secret.Type.String())
	r.frmEditUserSecret.GetFormItemByLabel("Name").(*tview.InputField).SetText(secret.Name)
	r.frmEditUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).SetText(secret.Meta, false)

	switch secret.Type {
	case model.CredsSecret:
		creds, _ := payload.(*model.CredsPayload)
		r.frmEditUserSecret.AddInputField("Login", creds.Login, 0, nil, nil)
		r.frmEditUserSecret.AddInputField("Password", creds.Password, 0, nil, nil)
	case model.TextSecret:
		text, _ := payload.(*model.TextPayload)
		r.frmEditUserSecret.AddTextArea("Text", text.Data, 0, 3, 0, nil)
	case model.BinarySecret:
		bin, _ := payload.(*model.BinaryPayload)
		r.frmEditUserSecret.AddTextArea("Binary", string(bin.Data), 0, 3, 0, nil)
	case model.CardSecret:
		card, _ := payload.(*model.CardPayload)
		r.frmEditUserSecret.
			AddInputField("Card number", card.CardNum, 0, nil, nil).
			AddInputField("Card holder", card.CardHolder, 0, nil, nil).
			AddInputField("Valid thru", card.ValidThru, 0, nil, nil).
			AddInputField("CVV", card.CVV, 0, nil, nil)
	}

	r.uiPages.SwitchToPage(_pageEditUserSecret)
}

func (r *App) doEditUserSecret() {
	secretName := r.lstSecrets[r.wdgLstSecrets.GetCurrentItem()].Name
	secret, err := r.tr.SecretGet(r.cltToken, secretName)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.showError(_msgInternalServerError, r.doReloadUserSecrets)
		case errors.Is(err, transport.ErrSecretNotFound):
			r.showError(_msgSecretNotFound, r.doReloadUserSecrets)
		}
		return
	}

	r.showEditUserSecretPage(secret)
}

func (r *App) doDeleteUserSecret() {
	secretName := r.frmEditUserSecret.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	if err := r.tr.SecretDelete(r.cltToken, secretName); err != nil {
		r.showError(_msgInternalServerError, r.doReloadUserSecrets)
		return
	}

	r.doReloadUserSecrets()
}

func (r *App) doSaveSecret() {
	curSecret := r.lstSecrets[r.wdgLstSecrets.GetCurrentItem()]

	updSecret := &model.SecretUpdate{
		Meta:    r.frmEditUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).GetText(),
		Version: curSecret.Version + 1,
	}

	var payload model.Payload
	switch curSecret.Type {
	case model.CredsSecret:
		payload = model.NewCredsPayload(
			r.frmEditUserSecret.GetFormItemByLabel("Login").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("Password").(*tview.InputField).GetText(),
		)
	case model.TextSecret:
		payload = model.NewTextPayload(
			r.frmEditUserSecret.GetFormItemByLabel("Text").(*tview.TextArea).GetText(),
		)
	case model.BinarySecret:
		payload = model.NewBinaryPayload([]byte(
			r.frmEditUserSecret.GetFormItemByLabel("Binary").(*tview.TextArea).GetText(),
		))
	case model.CardSecret:
		payload = model.NewCardPayload(
			r.frmEditUserSecret.GetFormItemByLabel("Card number").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("Card holder").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("Valid thru").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("CVV").(*tview.InputField).GetText(),
		)
	}

	err := r.tr.SecretUpdate(r.cltToken, curSecret.Name, updSecret, payload)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrInternalServerError):
			r.showError(_msgInternalServerError, r.doReloadUserSecrets)
		case errors.Is(err, transport.ErrSecretNotFound):
			r.showError(_msgSecretNotFound, r.doReloadUserSecrets)
		case errors.Is(err, transport.ErrSecretOutdated):
			r.showError(_msgSecretOutdated, r.doReloadUserSecrets)
		}
		return
	}

	r.doReloadUserSecrets()
}
