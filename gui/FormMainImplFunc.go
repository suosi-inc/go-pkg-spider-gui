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

var linkData *spider.LinkData

// OnBtnRequestClick 请求测试功能
func (f *TFormMain) btnRequestClick() {
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

// btnLinkRequestClick 链接提取功能
func (f *TFormMain) btnLinkRequestClick() {
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

	f.clearGridLink()

	var err error
	start := fun.Timestamp(true)
	if linkData, err = spider.GetLinkData(urlStr, strictDomain, timeout, maxRetry); err == nil {
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

		// 渲染所有表格
		f.renderGridLink()
	}
}

// clearGridLink 清空 GridLink
func (f *TFormMain) clearGridLink() {
	f.clearStringGrid(f.GridLinkContent, false)
	f.clearStringGrid(f.GridLinkList, false)
	f.clearStringGrid(f.GridLinkUnknown, false)
	f.clearStringGrid(f.GridLinkNone, false)
	f.clearStringGrid(f.GridLinkFilter, false)
	f.clearStringGrid(f.GridLinkDomain, false)
}

// renderGridLink 绘制 GridLink
func (f *TFormMain) renderGridLink() {
	// 处理 linkData.SubDomains 数据格式
	subdomains := make(map[string]string, 0)
	for subdomain, v := range linkData.SubDomains {
		subdomains[subdomain] = fun.ToString(v)
	}

	// 填充每个
	f.fillGridLink(f.GridLinkContent, linkData.LinkRes.Content)
	f.fillGridLink(f.GridLinkList, linkData.LinkRes.List)
	f.fillGridLink(f.GridLinkUnknown, linkData.LinkRes.Unknown)
	f.fillGridLink(f.GridLinkNone, linkData.LinkRes.None)
	f.fillGridLink(f.GridLinkFilter, linkData.Filters)
	f.fillGridLink(f.GridLinkDomain, subdomains)
}

func (f *TFormMain) fillGridLink(grid *vcl.TStringGrid, datas map[string]string) {
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

// btnLinkSearchClick 链接提取结果搜索功能
func (f *TFormMain) btnLinkSearchClick() {
	keyword := f.EditLinkSearch.Text()

	// 已经加载数据
	if linkData != nil {
		if !fun.Blank(keyword) {
			// 根据当前活动的 StringGrid 进行搜索
			var gridLink *vcl.TStringGrid
			activePage := f.PageControlLink.ActivePageIndex()
			switch activePage {
			case 0:
				gridLink = f.GridLinkContent
				f.searchGridLink(gridLink, keyword, linkData.LinkRes.Content)
			case 1:
				gridLink = f.GridLinkList
				f.searchGridLink(gridLink, keyword, linkData.LinkRes.List)
			case 2:
				gridLink = f.GridLinkUnknown
				f.searchGridLink(gridLink, keyword, linkData.LinkRes.Unknown)
			case 3:
				gridLink = f.GridLinkNone
				f.searchGridLink(gridLink, keyword, linkData.LinkRes.None)
			case 4:
				gridLink = f.GridLinkFilter
				f.searchGridLink(gridLink, keyword, linkData.Filters)
			case 5:
				gridLink = f.GridLinkDomain
				// 处理 linkData.SubDomains 数据格式
				subdomains := make(map[string]string, 0)
				for subdomain, v := range linkData.SubDomains {
					subdomains[subdomain] = fun.ToString(v)
				}
				f.searchGridLink(gridLink, keyword, subdomains)
			}
		} else {
			f.clearGridLink()
			f.renderGridLink()
		}

	} else {
		f.debug("Link Search Failed : linkData is empty")
	}
}

// searchGridLink
func (f *TFormMain) searchGridLink(grid *vcl.TStringGrid, keyword string, datas map[string]string) {
	count := len(datas)
	searchData := make(map[string]string, 0)

	if count > 0 {
		for key, value := range datas {
			if strings.Contains(key, keyword) {
				searchData[key] = value
			}
		}

		searchCount := len(searchData)
		if searchCount > 0 {
			f.debug("Link Search Result : " + fun.ToString(searchCount))

			f.clearStringGrid(grid, false)
			f.fillGridLink(grid, searchData)
		} else {
			f.debug("Link Search Failed : not found")
		}
	} else {
		f.debug("Link Search Failed : data is empty")
	}
}

// btnToolDomainRequestClick 辅助工具域名提取
func (f *TFormMain) btnToolDomainRequestClick() {
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

// openBrowser 打开系统浏览器
func (f *TFormMain) openBrowser(urlStr string) {
	if fun.Blank(urlStr) {
		f.debug("Request Link Failed : url is empty")
		return
	}

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

// removeToolBtnDown 取消功能栏按钮状态
func (f *TFormMain) removeToolBtnDown() {
	f.ToolBtnRequest.SetDown(false)
	f.ToolBtnDomain.SetDown(false)
	f.ToolBtnLink.SetDown(false)
	f.ToolBtnContent.SetDown(false)
	f.ToolBtnTool.SetDown(false)
}
