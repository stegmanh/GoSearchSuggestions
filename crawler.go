package main

import (
	"./htmlcrawler"
	"./redisqueue"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/lib/pq"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

var pool redis.Pool
var db *sql.DB

var disallowedUrls = make(map[string]bool)
var allowed []string
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

	disallowedUrls, allowed, err = htmlcrawler.LoadRobots("http://cnn.com")
	if err != nil {
		panic(err)
	}
	for _, url := range allowed {
		_, err := addUniqueToQueue(&pool, "urlexists", "messagequeue", url)
		if err != nil {
			fmt.Println(err)
		}
	}

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

func handleUrl(url string) {
	switch path.Ext(strings.ToLower(url)) {
	case ".html":
		handleHTML(url)
	case ".xml":
		handleXML(url)
	default:
		return
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

func handleXML(url string) {
	urls, err := htmlcrawler.GetXmlUrls(url)
	if err != nil {
		fmt.Println(err)
	}
	for _, url := range urls {
		_, err = addUniqueToQueue(&pool, "urlexists", "messagequeue", url)
		if err != nil {
			fmt.Println(err)
		}
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
