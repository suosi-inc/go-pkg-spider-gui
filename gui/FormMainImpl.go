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
	"github.com/ying32/govcl/vcl/types"
	"github.com/ying32/govcl/vcl/types/colors"
)

var (
	linkSearching = false
)

// ::private::
type TFormMainFields struct {
}

func (f *TFormMain) OnFormCreate(sender vcl.IObject) {
	f.PageControl.SetActivePageIndex(0)
}

// OnBtnRequestClick 请求测试功能
func (f *TFormMain) OnBtnRequestClick(sender vcl.IObject) {
	f.MemoRequest.SetText("")

	// Request Url
	urlStr := f.EditRequestUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Failed : url is empty")
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

		f.debug("Request Success : " + urlStr + ", use " + fun.ToString(use) + "ms")
		f.debug("\tCharset : " + fun.ToString(resp.Charset))
		f.debug("\tContent-Length : " + fun.ToString(resp.ContentLength))

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
		f.debug("Request Failed : " + err.Error())
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
		f.debug("DomainTop Failed : domain is empty")
		return
	}

	f.debug("DomainTop : " + domain)
	var top string
	if strings.HasPrefix(domain, "http") {
		top = extract.DomainTopFromUrl(domain)
	} else {
		top = extract.DomainTop(domain)
	}

	f.EditToolDomainResult.SetText(top)
	f.debug("\tResult : " + top)
}

func (f *TFormMain) OnBtnLinkRequestClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Link Failed : url is empty")
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

	f.clearStringGrid(f.GridLinkContent, false)
	f.clearStringGrid(f.GridLinkList, false)
	f.clearStringGrid(f.GridLinkUnknow, false)
	f.clearStringGrid(f.GridLinkNone, false)
	f.clearStringGrid(f.GridLinkFilter, false)
	f.clearStringGrid(f.GridLinkDomain, false)

	start := fun.Timestamp(true)
	if linkData, err := spider.GetLinkData(urlStr, strictDomain, timeout, maxRetry); err == nil {
		use := fun.Timestamp(true) - start

		result := fmt.Sprintf("\tResult : Content(%d), List(%d), Unknown(%d), None(%d), Filters(%d), Subdomains(%d)",
			len(linkData.LinkRes.Content),
			len(linkData.LinkRes.List),
			len(linkData.LinkRes.Unknown),
			len(linkData.LinkRes.None),
			len(linkData.Filters),
			len(linkData.SubDomains),
		)
		f.debug("Request Link Success : " + urlStr + ", use " + fun.ToString(use) + "ms")
		f.debug(result)

		f.renderGridLink(f.GridLinkContent, linkData.LinkRes.Content)
		f.renderGridLink(f.GridLinkList, linkData.LinkRes.List)
		f.renderGridLink(f.GridLinkUnknow, linkData.LinkRes.Unknown)
		f.renderGridLink(f.GridLinkNone, linkData.LinkRes.None)
		f.renderGridLink(f.GridLinkFilter, linkData.Filters)

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

func (f *TFormMain) OnBtnLinkSearchClick(sender vcl.IObject) {
	keyword := f.EditLinkSearch.Text()
	if fun.Blank(keyword) {
		f.debug("Link Search Failed : keyword is empty")
		return
	}

	// 根据当前活动的 StringGrid 进行搜索
	var gridLink *vcl.TStringGrid
	activePage := f.PageControlLink.ActivePageIndex()
	switch activePage {
	case 0:
		gridLink = f.GridLinkContent
	case 1:
		gridLink = f.GridLinkList
	case 2:
		gridLink = f.GridLinkUnknow
	case 3:
		gridLink = f.GridLinkNone
	case 4:
		gridLink = f.GridLinkFilter
	case 5:
		gridLink = f.GridLinkDomain
	}

	rowCount := gridLink.RowCount()
	if rowCount > 1 {
		var i int32
		var found bool
		for i = 1; i < gridLink.RowCount(); i++ {
			cell := gridLink.Cells(1, i)
			if strings.Contains(cell, keyword) {

				// 重绘单元格，触发 DrawCell 事件，进行高亮
				linkSearching = true
				gridLink.InvalidateCell(1, i)
				gridLink.SetRow(i)
				f.debug(fmt.Sprintf("Link Search Result : page(%d), row(%d)", activePage, i))
				return
			}
		}
		if !found {
			f.debug("Link Search Failed : not found")
		}
	} else {
		f.debug("Link Search Failed : data is empty")
	}
}

func (f *TFormMain) OnGridLinkContentDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkContent, aCol, aRow, aRect, aState)
}

func (f *TFormMain) OnGridLinkListDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkList, aCol, aRow, aRect, aState)
}

func (f *TFormMain) OnGridLinkUnknowDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkUnknow, aCol, aRow, aRect, aState)
}

func (f *TFormMain) OnGridLinkNoneDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkNone, aCol, aRow, aRect, aState)
}

func (f *TFormMain) OnGridLinkFilterDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkFilter, aCol, aRow, aRect, aState)
}

func (f *TFormMain) OnGridLinkDomainDrawCell(sender vcl.IObject, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	f.searchGridLinkDrawCell(f.GridLinkDomain, aCol, aRow, aRect, aState)
}

func (f *TFormMain) searchGridLinkDrawCell(grid *vcl.TStringGrid, aCol int32, aRow int32, aRect types.TRect, aState types.TGridDrawState) {
	if linkSearching {
		grid.Canvas().Pen().SetColor(colors.ClRed)
		grid.Canvas().Pen().SetWidth(1)
		grid.Canvas().Brush().SetStyle(types.BsClear)
		f.debug(fmt.Sprintf("%d, %d", aRow, aCol))
		grid.Canvas().Rectangle(aRect.Left+1, aRect.Top+1, aRect.Right-1, aRect.Bottom-1)
		linkSearching = false
	}
}

func (f *TFormMain) renderGridLink(grid *vcl.TStringGrid, datas map[string]string) {
	var i int32
	i = 1
	for key, value := range datas {
		grid.InsertColRow(false, i)
		grid.SetCells(1, i, key)
		grid.SetCells(2, i, fun.ToString(value))
		i++
	}
}

func (f *TFormMain) clearStringGrid(grid *vcl.TStringGrid, title bool) {
	var i, start int32

	if title {
		start = 0
	} else {
		start = 1
	}

	// 需要从最后一行向前删除
	for i = grid.RowCount() - 1; i >= start; i-- {
		grid.DeleteRow(i)
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
		f.debug("Request Link Failed : url is empty")
	}

	f.OpenBrowser(urlStr)
}

func (f *TFormMain) OnBtnRequestLinkClick(sender vcl.IObject) {
	urlStr := f.EditRequestUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Failed : url is empty")
	}

	f.OpenBrowser(urlStr)
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

func (f *TFormMain) removeToolBtnDown() {
	f.ToolBtnRequest.SetDown(false)
	f.ToolBtnDomain.SetDown(false)
	f.ToolBtnLink.SetDown(false)
	f.ToolBtnContent.SetDown(false)
	f.ToolBtnTool.SetDown(false)
}

func (f *TFormMain) OnBtnLinkOpenClick(sender vcl.IObject) {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Link Failed : url is empty")
	}

	f.OpenBrowser(urlStr)
}
