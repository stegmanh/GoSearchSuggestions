package htmlcrawler

import (
	"encoding/xml"
	"errors"
	"fmt" //Debugging
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
		//Check so we dont have super old site maps
		if len(smo.Lastmod) != 0 {
			//Simple layout will change
			layout := "2006-01-02T15:04:05-05:00"
			t, err := time.Parse(layout, smo.Lastmod)
			if err != nil {
				fmt.Println(smo.Lastmod)
				//Do something with error but we will just go on
			} else {
				//If 2 months old continue, else it will add to the pool
				if t.AddDate(0, 2, 0).Before(now) {
					continue
				}
			}
		}
		urls = append(urls, smo.Location)
	}
	return urls, nil
}
