package main

import (
	"./redisqueue"
	"bufio"
	"encoding/xml"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"sync"
	"time"
)

var pool redis.Pool

var robots = make(map[string]bool)
var mutex = &sync.Mutex{}

func main() {
	wg := new(sync.WaitGroup)
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
				exists, err := queue.HashAdd(pool, "urlexists", strings.TrimSpace(information[1]), "true")
				if err != nil {
					fmt.Println(err)
				}
				if exists != 1 {
					fmt.Println("Robots were already added to ")
					continue
				}
				queue.QueuePush(pool, "messagequeue", strings.TrimSpace(information[1]))
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
		url, err := queue.QueuePop(pool, "messagequeue")
		if err != nil || len(url) == 0 {
			fmt.Println(url, err)
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
		exists, err := queue.HashExists(pool, "urlexists", smo.Location)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if !exists {
			queue.HashAdd(pool, "urlexists", smo.Location, "true")
			queue.QueuePush(pool, "messagequeue", smo.Location)
		} else {
			fmt.Println("Already crawled: ", smo.Location)
		}
	}
}
