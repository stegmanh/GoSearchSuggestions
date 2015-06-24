package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var robots = make(map[string]bool)
var alreadyCrawled = make(map[string]bool)
var mutex = &sync.Mutex{}
var urls = make(chan string, 100000)

func main() {
	wg := new(sync.WaitGroup)
	loadRobots("http://cnn.com")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go worker(wg)
	}
	wg.Wait()
}

//TODO: Use path to join paths instead of concat
func loadRobots(root string) {
	resp, err := http.Get(root + "/robots.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		information := strings.SplitN(scanner.Text(), ":", 2)
		if len(information) == 2 {
			switch information[0] {
			case "Sitemap":
				urls <- strings.TrimSpace(information[1])
			case "Disallow":
				disallowedUrl := root + strings.TrimSpace(information[1])
				robots[disallowedUrl] = true
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("We got as err ", err)
	}
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case url := <-urls:
			go handleUrl(url)
		default:
			fmt.Println("Sleeping..")
			time.Sleep(3 * time.Second)
		}
	}
}

type SiteMapIndex struct {
	SiteMaps []SiteMap `xml:",any"`
}

type SiteMap struct {
	Location string `xml:"loc"`
	Lastmod  string `xml:"lastmod"`
}

func handleUrl(url string) {
	//Temp to handle extensions
	if path.Ext(url) != ".xml" {
		return
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var sm SiteMapIndex
	err = xml.Unmarshal(body, &sm)
	if err != nil {
		fmt.Println(err)
	}
	for _, smo := range sm.SiteMaps {
		mutex.Lock()
		_, exists := alreadyCrawled[smo.Location]
		if !exists {
			alreadyCrawled[smo.Location] = true
			urls <- smo.Location
		} else {
			fmt.Println("Already crawled: ", smo.Location+", ", len(alreadyCrawled))
		}
		mutex.Unlock()
	}
}
