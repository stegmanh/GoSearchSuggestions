package main

import (
	"./htmlcrawler"
	"./redisqueue"
	"bufio"
	"database/sql"
	"encoding/json"
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

var now = time.Now()

//Consider moving this + redis to its own repo -- This way we have consistent package accross all things
type crawlerInformation struct {
	Status string //Running, stopped, idle
	// Cpu         string            //Cpu usage -- Not doing yet
	// Ram         string            //Ram usage -- Not doing yet
	UrlsCrawled int               //# Of total urls crawled
	LastCrawled []string          //Last of some arbitrary number
	QueueSize   int               //Size of the queue
	IndexSize   int               //Size of DB
	Errors      map[string]string //Map of all errors (Maybe we dont use)
}

func (c *crawlerInformation) new() {
	c.Status = "Running"
	c.UrlsCrawled = 0
	c.LastCrawled = make([]string, 0)
	c.QueueSize = 0
	c.IndexSize = 0
	c.Errors = make(map[string]string)
}

//I dont think this should exist. Crawler should be about current crawl... Maybe store at shutdown?
func (c *crawlerInformation) storeSelf(pool *redis.Pool, hashName string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = redisqueue.HashAdd(pool, hashName, "data", string(data))
	return err
}

func (c *crawlerInformation) appendArray(value string) {
	//Hard coded here
	if len(c.LastCrawled) > 9 {
		temp := []string{value}
		c.LastCrawled = append(temp, c.LastCrawled[1:]...)
	} else {
		c.LastCrawled = append(c.LastCrawled, value)
	}
}

var info *crawlerInformation

type dbConfig struct {
	User     string
	Password string
	Db_name  string
	Ssl_mode string
}

type config struct {
	Db dbConfig
}

func main() {
	var c config
	loadConfig("config.json", &c)
	wg := new(sync.WaitGroup)
	var err error
	//Init the database
	//TODO: Own module just like we did redis
	connString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v", c.Db.User, c.Db.Password, c.Db.Db_name, c.Db.Ssl_mode)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}
	//Initilize redis information
	pool, err = redisqueue.MakePool()
	if err != nil {
		panic(err)
	}

	//loadRobots("http://cnn.com")

	//Load Crawler Information
	info = &crawlerInformation{}
	info.new()
	info.storeSelf(&pool, "crawlerstatus")

	//Number of goroutines to create to process urls
	for i := 0; i < 10; i++ {
		wg.Add(1)
		worker(wg)
	}
	wg.Wait()
}

func loadConfig(path string, c *config) {
	content, err := ioutil.ReadFile(path)
	//Panic because reading config is going to be useful
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, c)
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
				_, err := addUniqueToQueue(&pool, "urlexists", "messagequeue", strings.TrimSpace(information[1]))
				if err != nil {
					fmt.Println(err)
				}
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
		url, err := redisqueue.QueuePop(&pool, "messagequeue")
		if err != nil || len(url) == 0 {
			fmt.Println(len(url), err)
			fmt.Println("Sleeping..")
			time.Sleep(3 * time.Second)
		} else {
			//Consider running in goroutine sometime
			handleUrl(url)
			//Update Relevant Information
			//Todo maybe use a map... Concurrency or something
			mutex.Lock()
			info.UrlsCrawled++
			info.appendArray(url)
			info.storeSelf(&pool, "crawlerstatus")
			mutex.Unlock()
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
		//Check so we dont have super old site maps
		if len(smo.Lastmod) != 0 {
			//Simple layout will change
			layout := "2006-01-02T15:04:05-05:00"
			t, err := time.Parse(layout, smo.Lastmod)
			if err != nil {
				//Do something with error but we will just go on
			} else {
				if t.AddDate(0, 2, 0).Before(now) {
					fmt.Println("Too old!  ", t)
					continue
				}
			}
		}
		_, err = addUniqueToQueue(&pool, "urlexists", "messagequeue", smo.Location)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleHTML(url string) {
	pi, err := htmlcrawler.CrawlHTML(url)
	if err != nil {
		fmt.Println(err)
	}
	for _, url := range pi.Urls {
		_, err := addUniqueToQueue(&pool, "urlexists", "messagequeue", url)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = pi.StorePage(db)
	if err != nil {
		fmt.Println(err)
	}
}

//Returns 1 if added to queue, 0 if not
func addUniqueToQueue(pool *redis.Pool, hashName, queueName, toAdd string) (int, error) {
	exists, err := redisqueue.HashAdd(pool, hashName, toAdd, "true")
	if err != nil {
		return 0, err
	}
	if exists == 1 {
		redisqueue.QueuePush(pool, queueName, toAdd)
		return 1, nil
	}
	//Exists == 0 so the field already exists and we didnt add to the queue
	return 0, nil
}
