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

func (r *App) createCreateUserSecretPage() {
	r.frmCreateUserSecret = tview.NewForm()
	r.frmCreateUserSecret.
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
		AddButton("Create", r.doCreateUserSecret).
		AddButton("Back to list", r.doReloadUserSecrets)
	r.frmCreateUserSecret.
		SetBorder(true).
		SetTitle("Create secret")

	r.uiPages.AddPage(_pageCreateUserSecret, uiCenteredWidget(r.frmCreateUserSecret, 0, 10), true, false)
}

func (r *App) showCreateUserSecretCleared() {
	r.frmCreateUserSecret.GetFormItemByLabel("Type").(*tview.DropDown).SetCurrentOption(0)
	r.frmCreateUserSecret.SetFocus(r.frmCreateUser.GetFormItemIndex("Type"))
	r.frmCreateUserSecret.GetFormItemByLabel("Name").(*tview.InputField).SetText("")
	r.frmCreateUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).SetText("", true)
	r.uiPages.SwitchToPage(_pageCreateUserSecret)
}

func (r *App) showCreateUserSecret() {
	r.uiPages.SwitchToPage(_pageCreateUserSecret)
}

func (r *App) doCreateUserSecret() {
	secret := &model.Secret{
		Name: r.frmCreateUserSecret.GetFormItemByLabel("Name").(*tview.InputField).GetText(),
		Meta: r.frmCreateUserSecret.GetFormItemByLabel("Meta").(*tview.TextArea).GetText(),
	}
	_, fType := r.frmCreateUserSecret.GetFormItemByLabel("Type").(*tview.DropDown).GetCurrentOption()

	if secret.Name == "" {
		r.showError("Empty name", func() {
			r.uiPages.SwitchToPage(_pageCreateUserSecret)
		})
		return
	}
	if fType == "" {
		r.showError("Empty type", func() {
			r.uiPages.SwitchToPage(_pageCreateUserSecret)
		})
		return
	}

	var payload model.Payload

	switch fType {
	case model.CredsSecret.String():
		secret.Type = model.CredsSecret
		payload = model.NewCredsPayload(
			r.frmCreateUserSecret.GetFormItemByLabel("Login").(*tview.InputField).GetText(),
			r.frmCreateUserSecret.GetFormItemByLabel("Password").(*tview.InputField).GetText(),
		)
	case model.TextSecret.String():
		secret.Type = model.TextSecret
		payload = model.NewTextPayload(
			r.frmCreateUserSecret.GetFormItemByLabel("Text").(*tview.TextArea).GetText(),
		)
	case model.BinarySecret.String():
		secret.Type = model.BinarySecret
		fileData, err := os.ReadFile(r.frmCreateUserSecret.GetFormItemByLabel("File path").(*tview.InputField).GetText())
		if err != nil {
			r.showError(fmt.Sprintf("File read error: %v", err), r.showCreateUserSecret)
			return
		}
		payload = model.NewBinaryPayload(fileData)
	case model.CardSecret.String():
		secret.Type = model.CardSecret
		payload = model.NewCardPayload(
			r.frmCreateUserSecret.GetFormItemByLabel("Card number").(*tview.InputField).GetText(),
			r.frmCreateUserSecret.GetFormItemByLabel("Card holder").(*tview.InputField).GetText(),
			r.frmCreateUserSecret.GetFormItemByLabel("Valid thru").(*tview.InputField).GetText(),
			r.frmCreateUserSecret.GetFormItemByLabel("CVV").(*tview.InputField).GetText(),
		)
	}

	payloadRaw, err := gkMsgp.Serialize(payload.(msgp.Encodable))
	if err != nil {
		r.logger.Errorf("Create user secret Msgp serialize error: %v", err)
		r.showError(_msgClientError, r.showCreateUserSecretCleared)
		return
	}
	secret.PayloadRaw = payloadRaw

	err = r.tr.SecretCreate(r.cltToken, secret)
	if err != nil {
		switch {
		case errors.Is(err, transport.ErrSecretAlreadyExists):
			r.showError(_msgSecretAlreadyExists, r.showCreateUserSecret)
		case errors.Is(err, transport.ErrInternalError):
			r.showError(_msgClientError, r.showCreateUserSecret)
		case errors.Is(err, transport.ErrInternalServerError):
			r.showError(_msgInternalServerError, r.showCreateUserSecret)
		case errors.Is(err, transport.ErrSecretPayloadSizeExceeded):
			r.showError(_msgSecretPayloadSizeExceeded, r.showCreateUserSecret)
		case errors.Is(err, transport.ErrSecretInvalid):
			r.showError(_msgSecretInvalid, r.showCreateUserSecret)
		}
		return
	}

	r.doReloadUserSecrets()
}

func (r *App) uiChangeSecretPayloadFields(choosenType string, _ int) {
	metaIndex := r.frmCreateUserSecret.GetFormItemIndex("Meta")
	i := r.frmCreateUserSecret.GetFormItemCount() - 1
	for i != metaIndex {
		r.frmCreateUserSecret.RemoveFormItem(i)
		i--
	}

	switch choosenType {
	case model.CredsSecret.String():
		r.frmCreateUserSecret.
			AddInputField("Login", "", 0, nil, nil).
			AddInputField("Password", "", 0, nil, nil)
	case model.TextSecret.String():
		r.frmCreateUserSecret.
			AddTextArea("Text", "", 0, 3, 0, nil)
	case model.BinarySecret.String():
		r.frmCreateUserSecret.
			AddInputField("File path", "", 0, nil, nil)
	case model.CardSecret.String():
		r.frmCreateUserSecret.
			AddInputField("Card number", "", 0, nil, nil).
			AddInputField("Card holder", "", 0, nil, nil).
			AddInputField("Valid thru", "", 0, nil, nil).
			AddInputField("CVV", "", 0, nil, nil)
	}
}
