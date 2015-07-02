package htmlcrawler

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
)

type PageInformation struct {
	Title, Date string
	Urls        []string
}

//Helper function to traverse a HTML node and update the page information
func TraverseNode(n *html.Node, pi *PageInformation) {
	//Get and print the text inside the title element
	// Change this to check if Title node and then look at child -- Might only work for title node
	if pi.Title == "" && n.Type == html.ElementNode && n.Data == "title" {
		pi.Title = n.FirstChild.Data
	}
	// Parse all urls and append to array
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				pi.Urls = append(pi.Urls, attr.Val)
			}
		}
	}
	//property="og:pubdate
	//Might be kind of gross but it works! Date stored in content
	if pi.Date == "" && n.Type == html.ElementNode && n.Data == "meta" {
		content := ""
		for _, attr := range n.Attr {
			if attr.Key == "content" {
				content = attr.Val
			}
			if attr.Key == "property" && attr.Val == "og:pubdate" {
				pi.Date = content
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		TraverseNode(c, pi)
	}
}

//Takes a URL in string format and starts the recursive calls to TraverseNode that populate page information
func CrawlHTML(url string) (PageInformation, error) {
	resp, err := http.Get(url)
	pi := &PageInformation{Title: "", Date: "", Urls: make([]string, 0)}
	if err != nil {
		return *pi, err
	}
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return *pi, err
	}
	TraverseNode(doc, pi)
	return *pi, nil
}
