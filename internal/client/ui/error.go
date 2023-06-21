package ui

import "github.com/rivo/tview"

func (r *UIApp) createErrorPage() {
	r.frmError = tview.NewForm().
		AddTextView("Message", "", 0, 1, true, true).
		AddButton("Back", nil)
	r.frmError.
		SetBorder(true).
		SetTitle("Error")

	flex := uiCenteredWidget(r.frmError, 10, 1)
	r.uiPages.AddPage(_pageError, flex, true, false)
}

func (r *UIApp) showError(errMsg string, fnCallback func()) {
	r.frmError.GetFormItemByLabel("Message").(*tview.TextView).SetText(errMsg)
	r.frmError.SetFocus(r.frmError.GetFormItemIndex("Error"))

	btnBackInd := r.frmError.GetButtonIndex("Back")
	r.frmError.GetButton(btnBackInd).SetSelectedFunc(fnCallback)
	r.uiPages.SwitchToPage(_pageError)
}
