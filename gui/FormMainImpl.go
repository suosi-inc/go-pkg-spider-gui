package gui

import (
	"github.com/x-funs/go-fun"
	"github.com/ying32/govcl/vcl"
)

// ::private::
type TFormMainFields struct {
}

func (f *TFormMain) OnFormCreate(sender vcl.IObject) {
	f.PageControl.SetActivePageIndex(0)
}

func (f *TFormMain) OnBtnRequestClick(sender vcl.IObject) {
	f.btnRequestClick()
}

func (f *TFormMain) OnBtnRequestDefaultClick(sender vcl.IObject) {
	f.EditRequestUrl.SetText("https://www.163.com")
	f.EditRequestUa.SetText("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/108.0.0.0 Safari/537.36")
	f.CheckRequestType.SetChecked(true)
	f.EditRequestType.SetText("text/html")
	f.EditRequestLength.SetText("4096000")
	f.CheckRequestRedirect.SetChecked(false)
	f.EditRequestRedirect.SetEnabled(true)
	f.EditRequestRedirect.SetValue(2)
	f.MemoRequestHeader.SetText("")
	f.EditRequestProxy.SetText("")
	f.CheckRequestCharset.SetChecked(true)
	f.CheckRequestClean.SetChecked(true)
	f.EditRequestTimeout.SetText("30000")
}

func (f *TFormMain) OnBtnToolDomainRequestClick(sender vcl.IObject) {
	f.btnToolDomainRequestClick()

}

func (f *TFormMain) OnBtnLinkRequestClick(sender vcl.IObject) {
	f.btnLinkRequestClick()
}

func (f *TFormMain) OnBtnLinkSearchClick(sender vcl.IObject) {
	f.btnLinkSearchClick()
}

func (f *TFormMain) OnBtnRequestTipProxyClick(sender vcl.IObject) {
	f.EditRequestProxy.SetText("http://username:password@host:port")
}

func (f *TFormMain) OnBtnRequestTipHeaderClick(sender vcl.IObject) {
	f.MemoRequestHeader.SetText("")
	f.MemoRequestHeader.Append("X-Header : test-header")
}

func (f *TFormMain) OnToolBtnDebugClick(sender vcl.IObject) {
	if !f.PanelDebug.Visible() {
		f.SplitterDebug.SetVisible(true)
		f.PanelDebug.SetVisible(true)
		f.PageControl.AnchorSideBottom().SetControl(f.SplitterDebug)
	} else {
		f.SplitterDebug.SetVisible(false)
		f.PanelDebug.SetVisible(false)
		f.PageControl.AnchorSideBottom().SetControl(f)
	}
}

func (f *TFormMain) debug(str string) {
	if f.PanelDebug.Visible() {
		f.MemoDebug.Append(str)
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

func (f *TFormMain) OnBtnRequestOpenClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	f.openBrowser(urlStr)
}

func (f *TFormMain) OnBtnRequestLinkClick(sender vcl.IObject) {
	urlStr := f.EditRequestUrl.Text()

	f.openBrowser(urlStr)
}

func (f *TFormMain) OnToolBtnRequestClick(sender vcl.IObject) {
	f.removeToolBtnDown()
	f.ToolBtnRequest.SetDown(true)
	f.PageControl.SetActivePageIndex(0)
}

func (f *TFormMain) OnToolBtnDomainClick(sender vcl.IObject) {
	f.removeToolBtnDown()
	f.ToolBtnDomain.SetDown(true)
	f.PageControl.SetActivePageIndex(1)
}

func (f *TFormMain) OnToolBtnLinkClick(sender vcl.IObject) {
	f.removeToolBtnDown()
	f.ToolBtnLink.SetDown(true)
	f.PageControl.SetActivePageIndex(2)
	f.PageControlLink.SetActivePageIndex(0)
}

func (f *TFormMain) OnToolBtnContentClick(sender vcl.IObject) {
	f.removeToolBtnDown()
	f.ToolBtnContent.SetDown(true)
	f.PageControl.SetActivePageIndex(3)
}

func (f *TFormMain) OnToolBtnToolClick(sender vcl.IObject) {
	f.removeToolBtnDown()
	f.ToolBtnTool.SetDown(true)
	f.PageControl.SetActivePageIndex(4)
}

func (f *TFormMain) OnBtnLinkOpenClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Link Failed : url is empty")
	}

	f.openBrowser(urlStr)
}

func (f *TFormMain) OnBtnNewsRequestClick(sender vcl.IObject) {
	f.btnNewsRequestClick()
}

func (f *TFormMain) OnBtnNewsOpenClick(sender vcl.IObject) {
	urlStr := f.EditNewsUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request News Failed : url is empty")
	}

	f.openBrowser(urlStr)
}

func (f *TFormMain) OnBtnToolLangClick(sender vcl.IObject) {
	f.btnToolLangClick()
}

func (f *TFormMain) OnLabelToolLangTipClick(sender vcl.IObject) {

}

func (f *TFormMain) OnPageControlLinkChange(sender vcl.IObject) {

}

func (f *TFormMain) OnBtnDomainRequestClick(sender vcl.IObject) {
	f.btnDomainRequestClick()
}
