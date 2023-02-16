package gui

import (
	"github.com/ying32/govcl/vcl"
)

// ::private::
type TFormMainFields struct {
}

func (f *TFormMain) OnToolBtnDebugClick(sender vcl.IObject) {
	// 调试窗口按钮切换
	if !f.PanelDebug.Visible() {
		f.SplitterDebug.SetVisible(true)
		f.PanelDebug.SetVisible(true)

		side := vcl.AsAnchorSide(f.PageControl)
		side.SetControl(vcl.AsControl(f.SplitterDebug))
		f.MemoDebug.Append(side.ToString())

		f.PageControl.SetAnchorSideBottom(side)
	} else {
		f.SplitterDebug.SetVisible(false)
		f.PanelDebug.SetVisible(false)

		side := vcl.AsAnchorSide(f.PageControl)
		side.SetControl(f)

		f.PageControl.SetAnchorSideBottom(side)
	}
}

func (f *TFormMain) OnMenuDebugCopyClick(sender vcl.IObject) {
	f.MemoDebug.CopyToClipboard()
}

func (f *TFormMain) OnMenuDebugClearClick(sender vcl.IObject) {
	f.MemoDebug.SetText("")
}

func (f *TFormMain) OnCheckRequestTypeChange(sender vcl.IObject) {
	if f.CheckRequestType.Checked() {
		f.EditRequestType.SetEnabled(false)
	} else {
		f.EditRequestType.SetEnabled(true)
	}
}

func (f *TFormMain) OnCheckRequestRedirectChange(sender vcl.IObject) {
	if f.CheckRequestRedirect.Checked() {
		f.EditRequestRedirect.SetEnabled(false)
	} else {
		f.EditRequestRedirect.SetEnabled(true)
	}
}
