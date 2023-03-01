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
	"sync"
	"sync/atomic"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
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

	// 禁用按钮, 避免重复点击事件
	f.BtnRequest.SetEnabled(false)
	f.EditRequestUrl.SetEnabled(false)

	// 必须在新协程执行完回调恢复按钮状态, 因为事件是独立消息
	go func() {
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
			timeout = 10000
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

		// Headers 解析
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

			var memoText string
			if f.CheckRequestClean.Checked() {
				doc, _ := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
				doc.Find(spider.DefaultDocRemoveTags).Remove()
				body, _ := doc.Html()
				body = html.UnescapeString(body)
				memoText = fun.NormaliseLine(body)
			} else {
				memoText = fun.String(resp.Body)
			}

			// 主协程绘制
			vcl.ThreadSync(func() {
				f.debug("Request Success : " + urlStr + ", use " + fun.ToString(use) + "ms")
				f.debug("\tCharset : " + fun.ToString(resp.Charset))
				f.debug("\tContent-Length : " + fun.ToString(resp.ContentLength))

				f.MemoRequest.SetText(memoText)
			})

		} else {
			// 主协程绘制
			vcl.ThreadSync(func() {
				f.debug("Request Failed : " + err.Error())
			})
		}

		vcl.ThreadSync(func() {
			f.BtnRequest.SetEnabled(true)
			f.EditRequestUrl.SetEnabled(true)
		})
	}()

}

// btnLinkRequestClick 链接提取功能
func (f *TFormMain) btnLinkRequestClick() {
	urlStr := f.EditLinkUrl.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Link Failed : url is empty")
		return
	}

	f.BtnLinkRequest.SetEnabled(false)
	f.EditLinkUrl.SetEnabled(false)

	// 清空数据
	f.clearGridLink()

	go func() {
		// 超时时间
		timeout := fun.ToInt(f.EditLinkTimeout.Text())
		if timeout < 0 {
			timeout = 10000
		}

		// 限制域名
		strictDomain := true
		if !f.CheckLinkStrictDomain.Checked() {
			strictDomain = false
		}

		// 最大重试次数
		maxRetry := fun.ToInt(f.EditLinkRetry.Text())

		var err error
		start := fun.Timestamp(true)
		if linkData, err = spider.GetLinkData(urlStr, strictDomain, timeout, maxRetry); err == nil {
			use := fun.Timestamp(true) - start

			// 主线程绘制
			vcl.ThreadSync(func() {
				f.debug("Request Link Success : " + urlStr + ", use " + fun.ToString(use) + "ms")

				f.renderGridLink()
			})

		} else {
			// 主线程绘制
			vcl.ThreadSync(func() {
				f.debug("Request Link Failed : " + err.Error())
			})
		}

		vcl.ThreadSync(func() {
			f.BtnLinkRequest.SetEnabled(true)
			f.EditLinkUrl.SetEnabled(true)
		})
	}()

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

	result := fmt.Sprintf("\tRender : Content(%d), List(%d), Unknown(%d), None(%d), Filters(%d), Subdomains(%d)",
		len(linkData.LinkRes.Content),
		len(linkData.LinkRes.List),
		len(linkData.LinkRes.Unknown),
		len(linkData.LinkRes.None),
		len(linkData.Filters),
		len(linkData.SubDomains),
	)
	f.debug(result)
}

func (f *TFormMain) fillGridLink(grid *vcl.TStringGrid, datas map[string]string) {
	for key, value := range datas {
		i := grid.RowCount()
		grid.InsertColRow(false, i)
		grid.SetCells(1, i, key)
		grid.SetCells(2, i, fun.ToString(value))
	}
}

func (f *TFormMain) clearStringGrid(grid *vcl.TStringGrid, fixRow bool) {
	var i, start int32

	if fixRow {
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
			f.clearStringGrid(grid, false)
			f.fillGridLink(grid, searchData)

			f.debug("Link Search Result : " + fun.ToString(searchCount))
		} else {
			f.clearStringGrid(grid, false)
			f.debug("Link Search Failed : not found")
		}
	} else {
		f.debug("Link Search Failed : data is empty")
	}
}

func (f *TFormMain) btnDomainRequestClick() {
	domain := f.EditDomain.Text()

	if fun.Blank(domain) {
		f.debug("Request Domain Failed : domain is empty")
		return
	}

	f.BtnDomainRequest.SetEnabled(false)
	f.EditDomain.SetEnabled(false)

	// 清空数据
	f.clearDomainContent()
	f.clearStringGrid(f.GridDomainSubdomain, false)

	// 主域名请求
	go func() {
		// 超时时间
		timeout := fun.ToInt(f.EditDomainTimeout.Text())
		if timeout < 0 {
			timeout = 10000
		}

		// 最大重试次数
		maxRetry := fun.ToInt(f.EditDomainRetry.Text())

		start := fun.Timestamp(true)
		if domainRes, err := spider.DetectDomain(domain, timeout, maxRetry); err == nil {
			use := fun.Timestamp(true) - start

			// 填充基本信息
			charset := domainRes.Charset.Charset
			charsetPos := domainRes.Charset.CharsetPos
			lang := spider.LangEnZhMap[domainRes.Lang.Lang]
			langPos := domainRes.Lang.LangPos

			vcl.ThreadSync(func() {
				f.debug("Request Domain Success : " + domain + ", use " + fun.ToString(use) + "ms")
				f.GridDomainData.SetCells(1, 1, domainRes.Title)
				f.GridDomainData.SetCells(1, 2, domainRes.TitleClean)
				f.GridDomainData.SetCells(1, 3, domainRes.Description)
				f.GridDomainData.SetCells(1, 4, domainRes.Scheme)
				f.GridDomainData.SetCells(1, 5, charset+", "+charsetPos)
				f.GridDomainData.SetCells(1, 6, lang+", "+langPos)
				f.GridDomainData.SetCells(1, 7, domainRes.Country)
				f.GridDomainData.SetCells(1, 8, domainRes.Province)
				f.GridDomainData.SetCells(1, 9, fun.ToString(domainRes.State))
				f.GridDomainData.SetCells(1, 10, domainRes.Icp)
				f.GridDomainData.SetCells(1, 11, domainRes.HomeDomain)
				f.GridDomainData.SetCells(1, 12, fun.ToString(domainRes.ContentCount))
				f.GridDomainData.SetCells(1, 13, fun.ToString(domainRes.ListCount))
				f.GridDomainData.SetCells(1, 14, fun.ToString(len(domainRes.SubDomains)))
			})

			// 是否请求子域名
			if f.CheckDomainSubdomain.Checked() && len(domainRes.SubDomains) > 0 {
				f.domainSubdomainRequest(domainRes, timeout, maxRetry)
			}

		} else {
			vcl.ThreadSync(func() {
				f.debug("Request Domain Failed : " + err.Error())
			})
		}

		vcl.ThreadSync(func() {
			if !f.CheckDomainSubdomain.Checked() {
				f.BtnDomainRequest.SetEnabled(true)
				f.EditDomain.SetEnabled(true)
			}
		})
	}()
}

func (f *TFormMain) domainSubdomainRequest(domainRes *spider.DomainRes, timeout int, maxRetry int) {
	// 重置进度条
	var process int32
	total := len(domainRes.SubDomains)
	vcl.ThreadSync(func() {
		f.ProgressBarDomain.SetMax(int32(total))
	})

	// 数据分片，使用 5 个协程请求
	var wg sync.WaitGroup
	subdomainList := make([]string, 0)
	for s := range domainRes.SubDomains {
		subdomainList = append(subdomainList, s)
	}
	subdomainMutiParts := fun.SliceSplit(subdomainList, 5)
	for _, subdomainParts := range subdomainMutiParts {
		wg.Add(1)

		subdomainParts := subdomainParts
		go func() {

			for _, subdomain := range subdomainParts {
				atomic.AddInt32(&process, 1)

				subDomainRes, e := spider.DetectSubDomain(subdomain, timeout, maxRetry)

				vcl.ThreadSync(func() {

					i := f.GridDomainSubdomain.RowCount()
					f.GridDomainSubdomain.InsertColRow(false, i)
					f.GridDomainSubdomain.SetCells(1, i, subdomain)

					if e == nil {
						f.debug("\tRequest Domain Subdomain Success : " + subdomain)

						lang := ""
						if _, exist := spider.LangEnZhMap[subDomainRes.Lang.Lang]; exist {
							lang = spider.LangEnZhMap[subDomainRes.Lang.Lang]
						}

						f.GridDomainSubdomain.SetCells(2, i, subDomainRes.Title)
						f.GridDomainSubdomain.SetCells(3, i, subDomainRes.Charset.Charset)
						f.GridDomainSubdomain.SetCells(4, i, lang)
						f.GridDomainSubdomain.SetCells(5, i, fun.ToString(subDomainRes.State))
						f.GridDomainSubdomain.SetCells(6, i, fun.ToString(subDomainRes.ContentCount))
						f.GridDomainSubdomain.SetCells(7, i, fun.ToString(subDomainRes.ListCount))
					} else {
						f.debug("\tRequest Domain Subdomain Error : " + subdomain)

						f.GridDomainSubdomain.SetCells(2, i, e.Error())
						f.GridDomainSubdomain.SetCells(3, i, "")
						f.GridDomainSubdomain.SetCells(4, i, "")
						f.GridDomainSubdomain.SetCells(5, i, "false")
						f.GridDomainSubdomain.SetCells(6, i, "")
						f.GridDomainSubdomain.SetCells(7, i, "")
					}

					// 更新进度条
					f.ProgressBarDomain.SetPosition(atomic.LoadInt32(&process))
				})
			}

			wg.Done()
		}()
	}

	wg.Wait()

	// 所有协程执行完毕
	vcl.ThreadSync(func() {
		f.debug("\tRequest Domain Subdomain Done")
		f.BtnDomainRequest.SetEnabled(true)
		f.EditDomain.SetEnabled(true)
	})
}

func (f *TFormMain) clearDomainContent() {
	f.ProgressBarDomain.SetPosition(0)
	f.ProgressBarDomain.SetMax(100)

	for rows := 1; rows <= 14; rows++ {
		f.GridDomainData.SetCells(1, int32(rows), "")
	}
}

func (f *TFormMain) btnNewsRequestClick() {
	urlStr := f.EditNewsUrl.Text()
	title := f.EditNewsTitle.Text()
	if fun.Blank(urlStr) {
		f.debug("Request Content Failed : url is empty")
		return
	}

	f.BtnNewsRequest.SetEnabled(false)
	f.EditNewsUrl.SetEnabled(false)

	// 清空数据
	f.clearNewsContent()

	go func() {
		// 超时时间
		timeout := fun.ToInt(f.EditNewsTimeout.Text())
		if timeout < 0 {
			timeout = 10000
		}

		// 最大重试次数
		maxRetry := fun.ToInt(f.EditNewsRetry.Text())

		// 正文样式
		contentType := f.RadioNewsContentType.ItemIndex()

		if news, _, err := spider.GetNews(urlStr, title, timeout, maxRetry); err == nil {
			// Info
			contentData := news.ContentNode.Data
			contentAttr := fun.ToString(news.ContentNode.Attr)

			var htmlStr string
			switch contentType {
			case 0:
				htmlStr = strings.ReplaceAll(news.Content, fun.LF, fun.CRLF)
			case 1:
				node := goquery.NewDocumentFromNode(news.ContentNode)
				contentHtml, _ := node.Html()
				p := bluemonday.NewPolicy()
				p.AllowElements("p")
				p.AllowImages()
				htmlStr = p.Sanitize(contentHtml)
			case 2:
				node := goquery.NewDocumentFromNode(news.ContentNode)
				htmlStr, _ = node.Html()
			}

			vcl.ThreadSync(func() {
				f.GridNewsInfo.SetCells(1, 1, news.TitlePos)
				f.GridNewsInfo.SetCells(1, 2, news.TimePos)
				f.GridNewsInfo.SetCells(1, 3, contentData+contentAttr)
				f.GridNewsInfo.SetCells(1, 4, news.Time)
				f.GridNewsInfo.SetCells(1, 5, spider.LangEnZhMap[news.Lang])
				f.GridNewsInfo.SetCells(1, 6, fun.ToString(news.Spend)+"ms")
				f.EditNewsResultTitle.SetText(news.Title)
				f.EditNewsResultTime.SetText(news.TimeLocal)

				f.MemoNewsContent.Append(htmlStr)
			})
		} else {
			vcl.ThreadSync(func() {
				f.debug("Request News Failed : " + err.Error())
			})
		}

		vcl.ThreadSync(func() {
			f.BtnNewsRequest.SetEnabled(true)
			f.EditNewsUrl.SetEnabled(true)
		})
	}()
}

func (f *TFormMain) clearNewsContent() {
	f.EditNewsResultTitle.SetText("")
	f.EditNewsResultTime.SetText("")
	f.MemoNewsContent.SetText("")

	for rows := 1; rows <= 6; rows++ {
		f.GridNewsInfo.SetCells(1, int32(rows), "")
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

func (f *TFormMain) btnToolLangClick() {
	text := f.MemoToolLang.Text()
	if fun.Blank(text) {
		f.debug("Request Tool Lang Failed : text is empty")
		return
	}

	lang, _ := spider.LangText(text)

	f.EditToolLang.SetText("")

	if !fun.Blank(lang) {
		f.EditToolLang.SetText(spider.LangEnZhMap[lang])
		f.debug("\tResult : " + lang)
	}
}
