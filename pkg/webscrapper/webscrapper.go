package webscrapper

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strings"
)

func New(httpClient *http.Client) *WebScrapper {
	return &WebScrapper{
		httpClient: httpClient,
	}
}

type WebScrapper struct {
	httpClient *http.Client
}

type WebsiteData struct {
	HTMLVersion            string
	PageTitle              string
	HeadingsCountLevel1    int
	HeadingsCountLevel2    int
	HeadingsCountLevel3    int
	HeadingsCountLevel4    int
	HeadingsCountLevel5    int
	HeadingsCountLevel6    int
	InternalLinkCount      int
	ExternalLinkCount      int
	InaccessibleLinksCount int
	HasLoginForm           bool
}

type Link struct {
	Href     string
	Valid    bool
	External bool
}

func (webscrapper *WebScrapper) Extract(initURL string) (*WebsiteData, error) {

	// get the HTML document
	res, err := webscrapper.httpClient.Get(initURL)
	if err != nil {
		return nil, fmt.Errorf("could not get page: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	wd := new(WebsiteData)

	// todo extract data concurrently
	wd = getDocType(doc, wd)
	wd = getHxStats(doc, wd)
	wd = getTitle(doc, wd)
	wd = extractLinksStats(doc, wd, initURL)
	wd = hasLogin(doc, wd)

	return wd, nil
}

func getDocType(doc *html.Node, wd *WebsiteData) *WebsiteData {
	if len(doc.FirstChild.Attr) != 0 {
		wd.HTMLVersion = doc.FirstChild.Attr[0].Val
	} else {
		wd.HTMLVersion = "HTML5 and beyond"
	}

	return wd
}

func getTitle(doc *html.Node, wd *WebsiteData) *WebsiteData {
	var title string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "title" {
			title = n.FirstChild.Data
			return
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	wd.PageTitle = title

	return wd
}

func getHxStats(doc *html.Node, websiteData *WebsiteData)  *WebsiteData {
	var countH1 int
	var countH2 int
	var countH3 int
	var countH4 int
	var countH5 int
	var countH6 int
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "h1" {
			countH1++
		} else if n.Type == html.ElementNode && n.Data == "h2" {
			countH2++
		} else if n.Type == html.ElementNode && n.Data == "h3" {
			countH3++
		} else if n.Type == html.ElementNode && n.Data == "h4" {
			countH4++
		} else if n.Type == html.ElementNode && n.Data == "h5" {
			countH5++
		} else if n.Type == html.ElementNode && n.Data == "h6" {
			countH6++
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	websiteData.HeadingsCountLevel1 = countH1
	websiteData.HeadingsCountLevel2 = countH2
	websiteData.HeadingsCountLevel3 = countH3
	websiteData.HeadingsCountLevel4 = countH4
	websiteData.HeadingsCountLevel5 = countH5
	websiteData.HeadingsCountLevel6 = countH6
	return websiteData
}

func extractLinksStats(doc *html.Node, websiteData *WebsiteData, initURL string) *WebsiteData {

	var internalLinkCount int
	var externalLinkCount int
	var inaccessibleLinksCount int
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			v, validLink := getAttr(n.Attr, "href")

			if v == "" || !validLink {
				inaccessibleLinksCount++
				return
			}

			if len(v) != 0 && v[0] == 'h' {
				_, err := url.Parse(v)
				if err != nil {
					validLink = false
				} else {
					validLink = true
				}

				externalLinkCount++
			} else if len(v) != 0 && v[0] == '/' {

				// parse input url
				urlParsed, err := url.Parse(initURL)
				if err != nil {
					validLink = false
				} else {
					validLink = true

					host := urlParsed.Scheme + "://" + urlParsed.Hostname()
					newURL := host + v

					_, err = url.Parse(newURL)
					if err != nil {
						validLink = false
					} else {
						validLink = true
						internalLinkCount++
					}
				}
			}

			if !validLink {
				inaccessibleLinksCount++
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	websiteData.ExternalLinkCount = externalLinkCount
	websiteData.InternalLinkCount = internalLinkCount
	websiteData.InaccessibleLinksCount = inaccessibleLinksCount

	return websiteData
}

func hasLogin(doc *html.Node, wd *WebsiteData) *WebsiteData {
	var hasLogin bool
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			class, ok := getAttr(n.Attr, "class")
			if ok {
				if containsLoginKeyWord(strings.ToLower(class)) {
					hasLogin = true
					return
				}
			}

			id, ok := getAttr(n.Attr, "id")
			if ok {
				if containsLoginKeyWord(strings.ToLower(id)) {
					hasLogin = true
					return
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	wd.HasLoginForm = hasLogin

	return wd
}

func containsLoginKeyWord(s string) bool {
	if strings.Contains(s, "auth") {
		return true
	}
	if strings.Contains(s, "signin") {
		return true
	}
	if strings.Contains(s, "login") {
		return true
	}
	return false
}

func getAttr(attrs []html.Attribute, attrName string) (string, bool) {
	var attrVal string
	var found bool
	for _, attr := range attrs {
		if attr.Key == attrName {
			attrVal = attr.Val
			found = true
		}
	}
	return attrVal, found
}
