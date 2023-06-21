package ui

import "github.com/rivo/tview"

func (r *App) createErrorPage() {
	r.frmError = tview.NewForm().
		AddTextView("Message", "", 0, 1, true, true).
		AddButton("Ok", nil)
	r.frmError.
		SetBorder(true)

	flex := uiCenteredWidget(r.frmError, 10, 1)
	r.uiPages.AddPage(_pageError, flex, true, false)
}

func (r *App) showError(errMsg string, fnCallback func()) {
	r.frmError.GetFormItemByLabel("Message").(*tview.TextView).SetText(errMsg)

	btnBackInd := r.frmError.GetButtonIndex("Ok")
	r.frmError.GetButton(btnBackInd).SetSelectedFunc(fnCallback)
	r.uiPages.SwitchToPage(_pageError)
}
