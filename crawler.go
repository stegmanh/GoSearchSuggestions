package main

import (
	"./htmlcrawler"
	"./redisqueue"
	"bufio"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var pool redis.Pool
var db *sql.DB

var robots = make(map[string]bool)
var mutex = &sync.Mutex{}

func main() {
	wg := new(sync.WaitGroup)
	var err error
	//Init the database
	db, err = sql.Open("postgres", "user=postgres password=qwert12345 dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
	//Initilize redis information
	pool = redis.Pool{
		MaxIdle:   50,
		MaxActive: 0,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}

	loadRobots("http://cnn.com")

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
				exists, err := queue.HashAdd(&pool, "urlexists", strings.TrimSpace(information[1]), "true")
				if err != nil {
					fmt.Println(err)
				}
				if exists != 1 {
					fmt.Println("Robots were already added to ")
					continue
				}
				queue.QueuePush(&pool, "messagequeue", strings.TrimSpace(information[1]))
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
		url, err := queue.QueuePop(&pool, "messagequeue")
		if err != nil || len(url) == 0 {
			fmt.Println(len(url), err)
			fmt.Println("Sleeping..")
			time.Sleep(3 * time.Second)
		} else {
			//Consider running in goroutine sometime
			handleUrl(url)
			//Sleep to slow things down...
			time.Sleep(time.Millisecond * 50)
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
	//TODO Handle the HTML files
	if path.Ext(url) == ".html" {
		handleHTML(url)
		return
	}
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
		exists, err := queue.HashAdd(&pool, "urlexists", smo.Location, "true")
		if err != nil {
			fmt.Println(err)
			continue
		}
		// If it was new...
		if exists == 1 {
			queue.QueuePush(&pool, "messagequeue", smo.Location)
		}
	}
}

func handleHTML(url string) {
	pi, err := htmlcrawler.CrawlHTML(url)
	if err != nil {
		fmt.Println(err)
	}
	for _, url := range pi.Urls {
		exists, err := queue.HashAdd(&pool, "urlexists", url, "true")
		if err != nil {
			fmt.Println(err)
			continue
		}
		if exists == 1 {
			queue.QueuePush(&pool, "messagequeue", url)
		}
	}
	err = pi.StorePage(db)
	if err != nil {
		fmt.Println(err)
	}
}
