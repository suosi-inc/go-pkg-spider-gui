package gui

import (
	"bytes"
	"html"

	"github.com/PuerkitoBio/goquery"
	spider "github.com/suosi-inc/go-pkg-spider"
	"github.com/x-funs/go-fun"
	"github.com/ying32/govcl/vcl"
)

// ::private::
type TFormMainFields struct {
}

// OnToolBtnDebugClick 调试窗口按钮切换
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

// OnMenuDebugCopyClick 调试窗口复制
func (f *TFormMain) OnMenuDebugCopyClick(sender vcl.IObject) {
	f.MemoDebug.CopyToClipboard()
}

// OnMenuDebugClearClick 调试窗口清除
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

func (f *TFormMain) OnBtnRequestClick(sender vcl.IObject) {
	f.MemoRequest.SetText("")
	urlStr := f.EditRequestUrl.Text()

	timeout := fun.ToInt(f.EditRequestTimeout.Text())
	if timeout < 0 {
		timeout = 30000
	}

	start := fun.Timestamp(true)
	if resp, err := spider.HttpGetResp(urlStr, nil, timeout); err == nil {
		use := fun.Timestamp(true) - start

		f.Debug("Request Success : " + urlStr + ", use " + fun.ToString(use) + "ms")
		f.Debug("\tCharset : " + fun.ToString(resp.Charset))
		f.Debug("\tContent-Length : " + fun.ToString(resp.ContentLength))

		if f.CheckRequestClean.Checked() {
			doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
			doc.Find(spider.DefaultDocRemoveTags).Remove()
			body, _ := doc.Html()
			body = html.UnescapeString(body)
			body = fun.NormaliseLine(body)
			f.MemoRequest.SetText(body)
		} else {
			f.MemoRequest.SetText(fun.String(resp.Body))
		}

	} else {
		f.Debug("Request Failed : " + err.Error())
	}
}

func (f *TFormMain) Debug(str string) {
	f.MemoDebug.Append(str)
}
