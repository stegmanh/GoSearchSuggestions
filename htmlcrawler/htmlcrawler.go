package htmlcrawler

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"golang.org/x/net/html"
	"net/http"
	"regexp"
	"strings"
)

var validUrl = regexp.MustCompile(`http[s]?://.*cnn.com.*`)

type PageInformation struct {
	Title, Date, Source string
	Body                []byte
	Urls                []string
}

func (pi *PageInformation) StorePage(db *sql.DB) error {
	//Changed to query row because query would give error with mismatched "multiple-value db.Query() in single-value context"
	rows, err := db.Query("INSERT INTO articles (title, created_at, source, body) VALUES ($1, $2, $3, $4) RETURNING title", pi.Title, pi.Date, pi.Source, string(pi.Body))
	if rows != nil {
		rows.Close()
	}
	return err
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
				url := attr.Key
				if strings.HasPrefix(url, "/") {
					url = "http://cnn.com" + url
				}
				if !validUrl.MatchString(url) {
					continue
					fmt.Println("Invalid URL to crawl ", url)
				}
				pi.Urls = append(pi.Urls, url)
			}
		}
	}

	if len(pi.Body) == 0 && n.Type == html.ElementNode && n.Data == "section" {
		for _, attr := range n.Attr {
			if attr.Key == "class" {
				classNames := strings.Split(attr.Val, " ")
				for _, class := range classNames {
					if class == "zn-body-text" {
						getText(n, pi)
					}
				}
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
	pi := &PageInformation{Title: "", Date: "", Source: url, Body: make([]byte, 0), Urls: make([]string, 0)}
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

//Consider changing this to grab images as well
func getText(n *html.Node, pi *PageInformation) {
	// I hate that this statememt exists
	if n.Type == html.TextNode && n.Parent.Data != "" && (n.Parent.Data == "p" || n.Parent.Data == "cite" || n.Parent.Data == "h3" || n.Parent.Data == "h4") {
		formattedBody := fmt.Sprintf("<%v>%v</%v>", n.Parent.Data, n.Data, n.Parent.Data)
		pi.Body = append(pi.Body, formattedBody...)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getText(c, pi)
	}
}
