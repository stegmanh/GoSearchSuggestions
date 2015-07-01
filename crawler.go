package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"path"
	"queue"
	"strings"
	"sync"
	"time"
)

var pool = redis.Pool

var robots = make(map[string]bool)
var mutex = &sync.Mutex{}

func main() {
	wg := new(sync.WaitGroup)
	loadRobots("http://cnn.com")

	//Initilize redis information
	pool = CreatePool()

	//Number of goroutines to create to process urls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		worker(wg)
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
		panic(err)
	}
}

func worker(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case url := <-urls:
			//Consider running in goroutine sometime
			handleUrl(url)
			//Sleep to slow things down...
			time.Sleep(time.Millisecond * 50)
		default:
			fmt.Println("Sleeping..")
			time.Sleep(3 * time.Second)
		}
	}
}

//XML structures
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
		exists, err := HashExists(&pool, "urlexists", smo.Location)
		if err {
			fmt.Println(err)
			continue
		}
		if !exists {
			alreadyCrawled[smo.Location] = true
			QueuePush(&pool, "messagequeue", value)
		} else {
			fmt.Println("Already crawled: ", smo.Location+", ", len(alreadyCrawled))
		}
	}
}
