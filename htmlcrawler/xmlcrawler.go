package htmlcrawler

import (
	"encoding/xml"
	"errors"
	//"fmt" //Debugging
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

type SiteMapIndex struct {
	SiteMaps []SiteMap `xml:",any"`
}

type SiteMap struct {
	Location string `xml:"loc"`
	Lastmod  string `xml:"lastmod"`
}

func GetXmlUrls(url string) ([]string, error) {
	now := time.Now()
	urls := make([]string, 0)
	if path.Ext(url) != ".xml" {
		return urls, errors.New("Not an XML file")
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return urls, errors.New("Error reaching the url" + url)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return urls, errors.New("Error reading the file")
	}
	var sm SiteMapIndex
	err = xml.Unmarshal(body, &sm)
	if err != nil {
		return urls, errors.New("Error unmarshing the XML")
	}
	for _, smo := range sm.SiteMaps {
		if len(smo.Lastmod) != 0 {
			layouts := []string{"2006-01-02T15:04:05-05:00", "2006-01-02T15:04:05Z", "2006-01-02T15:04:05z", "2006-01-02T15:04:05z"}
			t, err := parseAgainstLayouts(layouts, smo.Lastmod)
			if err != nil || t.AddDate(0, 2, 0).Before(now) {
				continue
			}
		}
		urls = append(urls, smo.Location)
	}
	return urls, nil
}

//Pass in an array of layouts and try parsing the string against each
//Returns when layout parse doesn't return error or out of layouts
func parseAgainstLayouts(layouts []string, value string) (time.Time, error) {
	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.Parse(layout, value)
		if err != nil {
			continue
		} else {
			return t, nil
		}
	}
	return t, err
}
