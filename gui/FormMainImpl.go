package gui

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"html"
	"net/http"
	"net/textproto"
	"net/url"
	"os/exec"
	"runtime"
	"strings"

	"github.com/PuerkitoBio/goquery"
	spider "github.com/suosi-inc/go-pkg-spider"
	"github.com/suosi-inc/go-pkg-spider/extract"
	"github.com/x-funs/go-fun"
	"github.com/ying32/govcl/vcl"
)

// ::private::
type TFormMainFields struct {
}

// OnBtnRequestClick 请求测试功能运行
func (f *TFormMain) OnBtnRequestClick(sender vcl.IObject) {
	f.MemoRequest.SetText("")

	// Request Url
	urlStr := f.EditRequestUrl.Text()
	if fun.Blank(urlStr) {
		f.Debug("Request Failed : url is empty")
		return
	}

	req := &spider.HttpReq{
		HttpReq: &fun.HttpReq{
			DisableRedirect: true,
			Headers:         make(map[string]string),
		},
		ForceTextContentType: true,
		DisableCharset:       false,
	}

	// 超时时间
	timeout := fun.ToInt(f.EditRequestTimeout.Text())
	if timeout < 0 {
		timeout = 30000
	}

	// UserAgent
	if !fun.Blank(f.EditRequestUa.Text()) {
		req.HttpReq.UserAgent = f.EditRequestUa.Text()
	}

	// ContentType
	if !f.CheckRequestType.Checked() {
		if !fun.Blank(f.EditRequestType.Text()) {
			req.ForceTextContentType = false
			req.AllowedContentTypes = []string{f.EditRequestType.Text()}
		}
	}

	// MaxContentLength
	contentLength := fun.ToInt64(f.EditRequestLength.Text())
	if contentLength > 0 {
		req.HttpReq.MaxContentLength = contentLength
	}

	// MaxRedirect
	if !f.CheckRequestRedirect.Checked() {
		maxRedirect := fun.ToInt(f.EditRequestRedirect.Text())
		req.DisableRedirect = false
		req.HttpReq.MaxRedirect = maxRedirect
	}

	// Headers
	if !fun.Blank(f.MemoRequestHeader.Text()) {
		headerText := f.MemoRequestHeader.Text()

		buffer := bytes.NewBufferString(headerText)
		scanner := bufio.NewScanner(buffer)
		mimeHeader := textproto.MIMEHeader{}
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				mimeHeader.Add(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}

		headerMap := make(map[string]string, len(mimeHeader))
		for key, values := range mimeHeader {
			if len(values) > 0 {
				headerMap[key] = values[0]
			}
		}
		req.HttpReq.Headers = headerMap
	}

	// Proxy
	if !fun.Blank(f.EditRequestProxy.Text()) {
		transport := &http.Transport{
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			DisableKeepAlives: true,
		}
		proxyString := f.EditRequestProxy.Text()
		proxy, _ := url.Parse(proxyString)
		transport.Proxy = http.ProxyURL(proxy)

		req.HttpReq.Transport = transport
	}

	// Charset
	if !f.CheckRequestCharset.Checked() {
		req.DisableCharset = true
	}

	start := fun.Timestamp(true)
	if resp, err := spider.HttpGetResp(urlStr, req, timeout); err == nil {
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
	if f.PanelDebug.Visible() {
		f.MemoDebug.Append(str)
	}
}

// OnBtnRequestDefaultClick 请求测试功能默认参数
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
	domain := f.EditToolDomain.Text()
	if fun.Blank(domain) {
		f.Debug("DomainTop Failed : domain is empty")
		return
	}

	f.Debug("DomainTop : " + domain)
	var top string
	if strings.HasPrefix(domain, "http") {
		top = extract.DomainTopFromUrl(domain)
	} else {
		top = extract.DomainTop(domain)
	}

	f.EditToolDomainResult.SetText(top)
	f.Debug("\tResult : " + top)
}

func (f *TFormMain) OnBtnLinkRequestClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.Debug("Request Link Failed : url is empty")
		return
	}

	// 超时时间
	timeout := fun.ToInt(f.EditLinkTimeout.Text())
	if timeout < 0 {
		timeout = 30000
	}

	// 限制域名
	strictDomain := true
	if !f.CheckLinkStrictDomain.Checked() {
		strictDomain = false
	}

	// 最大重试次数
	maxRetry := fun.ToInt(f.EditLinkRetry.Text())

	f.ClearStringGrid(f.GridLinkContent, false)
	f.ClearStringGrid(f.GridLinkList, false)
	f.ClearStringGrid(f.GridLinkUnknow, false)
	f.ClearStringGrid(f.GridLinkNone, false)
	f.ClearStringGrid(f.GridLinkFilter, false)
	f.ClearStringGrid(f.GridLinkDomain, false)

	start := fun.Timestamp(true)
	if linkData, err := spider.GetLinkData(urlStr, strictDomain, timeout, maxRetry); err == nil {
		use := fun.Timestamp(true) - start

		f.Debug("Request Link Success : " + urlStr + ", use " + fun.ToString(use) + "ms")
		f.Debug("\tResult : ")

		f.RenderGridLink(f.GridLinkContent, linkData.LinkRes.Content)
		f.RenderGridLink(f.GridLinkList, linkData.LinkRes.List)
		f.RenderGridLink(f.GridLinkUnknow, linkData.LinkRes.Unknown)
		f.RenderGridLink(f.GridLinkNone, linkData.LinkRes.None)
		f.RenderGridLink(f.GridLinkFilter, linkData.Filters)

		var i int32
		i = 1
		for subdomain, t := range linkData.SubDomains {
			f.GridLinkDomain.InsertColRow(false, i)
			f.GridLinkDomain.SetCells(1, i, subdomain)
			f.GridLinkDomain.SetCells(2, i, fun.ToString(t))
			i++
		}
	}
}

func (f *TFormMain) RenderGridLink(grid *vcl.TStringGrid, datas map[string]string) {
	var i int32
	i = 1
	for key, value := range datas {
		grid.InsertColRow(false, i)
		grid.SetCells(1, i, key)
		grid.SetCells(2, i, fun.ToString(value))
		i++
	}
}

func (f *TFormMain) ClearStringGrid(grid *vcl.TStringGrid, title bool) {
	var i, start int32

	if title {
		start = 0
	} else {
		start = 1
	}
	for i = start; i < grid.RowCount(); i++ {
		grid.DeleteColRow(false, i)
	}
}

func (f *TFormMain) OnBtnRequestTipProxyClick(sender vcl.IObject) {
	f.EditRequestProxy.SetText("http://username:password@host:port")
}

func (f *TFormMain) OnBtnRequestTipHeaderClick(sender vcl.IObject) {
	f.MemoRequestHeader.SetText("")
	f.MemoRequestHeader.Append("X-Header : test-header")
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

func (f *TFormMain) OpenBrowser(urlStr string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", urlStr).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", urlStr).Start()
	case "darwin":
		err = exec.Command("open", urlStr).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
}

func (f *TFormMain) OnBtnRequestOpenClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.Debug("Request Link Failed : url is empty")
	}

	f.OpenBrowser(urlStr)
}

func (f *TFormMain) OnBtnRequestLinkClick(sender vcl.IObject) {
	urlStr := f.EditRequestUrl.Text()
	if fun.Blank(urlStr) {
		f.Debug("Request Failed : url is empty")
	}

	f.OpenBrowser(urlStr)
}

func (f *TFormMain) OnToolBtnRequestClick(sender vcl.IObject) {
	f.RemoveToolBtnDown()
	f.ToolBtnRequest.SetDown(true)
	f.PageControl.SetActivePageIndex(0)
}

func (f *TFormMain) OnToolBtnDomainClick(sender vcl.IObject) {
	f.RemoveToolBtnDown()
	f.ToolBtnDomain.SetDown(true)
	f.PageControl.SetActivePageIndex(1)
}

func (f *TFormMain) OnToolBtnLinkClick(sender vcl.IObject) {
	f.RemoveToolBtnDown()
	f.ToolBtnLink.SetDown(true)
	f.PageControl.SetActivePageIndex(2)
}

func (f *TFormMain) OnToolBtnContentClick(sender vcl.IObject) {
	f.RemoveToolBtnDown()
	f.ToolBtnContent.SetDown(true)
	f.PageControl.SetActivePageIndex(3)
}

func (f *TFormMain) OnToolBtnToolClick(sender vcl.IObject) {
	f.RemoveToolBtnDown()
	f.ToolBtnTool.SetDown(true)
	f.PageControl.SetActivePageIndex(4)
}

func (f *TFormMain) RemoveToolBtnDown() {
	f.ToolBtnRequest.SetDown(false)
	f.ToolBtnDomain.SetDown(false)
	f.ToolBtnLink.SetDown(false)
	f.ToolBtnContent.SetDown(false)
	f.ToolBtnTool.SetDown(false)
}
