package ui

import (
	"errors"
	"fmt"
	"os"

	"github.com/devldavydov/gophkeeper/internal/client/transport"
	"github.com/devldavydov/gophkeeper/internal/common/model"
	gkMsgp "github.com/devldavydov/gophkeeper/internal/common/msgp"
	"github.com/rivo/tview"
	"github.com/tinylib/msgp/msgp"
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
	r.frmEditUserSecret.
		SetBorder(true).
		SetTitle("Edit secret")

	r.frmEditUserSecret.GetFormItemByLabel("Type").SetDisabled(true)
	r.frmEditUserSecret.GetFormItemByLabel("Name").SetDisabled(true)

	r.uiPages.AddPage(_pageEditUserSecret, uiCenteredWidget(r.frmEditUserSecret, 0, 10), true, false)
}

func (r *App) showEditUserSecretPageClear(secret *model.Secret) {
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
	if btnIndex := r.frmEditUserSecret.GetButtonIndex("Copy file"); btnIndex != -1 {
		r.frmEditUserSecret.RemoveButton(btnIndex)
	}
	r.frmEditUserSecret.SetFocus(metaIndex)

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
		r.frmEditUserSecret.AddInputField("File path", "", 0, nil, nil)

		bin, _ := payload.(*model.BinaryPayload)
		saveName := fmt.Sprintf("/tmp/GophKeeper_%s", secret.Name)
		r.frmEditUserSecret.AddButton("Copy file", func() {
			var msg string
			if err = os.WriteFile(saveName, bin.Data, 0600); err != nil {
				msg = fmt.Sprintf("File copy error: %v", err)
			} else {
				msg = fmt.Sprintf("File copied: %s", saveName)
			}
			r.showError(msg, r.showEditUserSecretPage)
		})
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

func (r *App) showEditUserSecretPage() {
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

	r.showEditUserSecretPageClear(secret)
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
		Meta:          r.frmEditUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).GetText(),
		Version:       curSecret.Version + 1,
		UpdatePayload: true,
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
		filePath := r.frmEditUserSecret.GetFormItemByLabel("File path").(*tview.InputField).GetText()
		// Save new file payload, only if path set
		if filePath != "" {
			fileData, err := os.ReadFile(filePath)
			if err != nil {
				r.showError(fmt.Sprintf("File read error: %v", err), r.showEditUserSecretPage)
				return
			}
			payload = model.NewBinaryPayload(fileData)
		}
	case model.CardSecret:
		payload = model.NewCardPayload(
			r.frmEditUserSecret.GetFormItemByLabel("Card number").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("Card holder").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("Valid thru").(*tview.InputField).GetText(),
			r.frmEditUserSecret.GetFormItemByLabel("CVV").(*tview.InputField).GetText(),
		)
	}

	var err error
	if payload != nil {
		var payloadRaw []byte
		payloadRaw, err = gkMsgp.Serialize(payload.(msgp.Encodable))
		if err != nil {
			r.showError(_msgClientError, r.showEditUserSecretPage)
			return
		}
		updSecret.PayloadRaw = payloadRaw
	} else {
		updSecret.UpdatePayload = false
	}

	err = r.tr.SecretUpdate(r.cltToken, curSecret.Name, updSecret)
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
