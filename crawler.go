package main

import (
	"./crawlerinformation"
	"./htmlcrawler"
	"./redisqueue"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"
)

var db *sql.DB

var disallowedUrls = make(map[string]bool)
var allowed []string
var mutex = &sync.Mutex{}

var now = time.Now()

//TODO: Consider polling CNN homepage or main sitemap constantly for up to date information!
//TODO: Move dispatcher to a seperate file/location -- Maybe html crawler?
//TODO: Make sure we only crawl CNN
//TODO: Figure out why we can only crawl 5000 ulrs
//TODO: Make sure we append base to start of URL in cases where we have relative links
var info *crawlerinformation.CrawlerInformation

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
	var err error

	//Init the database
	//TODO: Own module just like we did redis
	connString := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v", c.Db.User, c.Db.Password, c.Db.Db_name, c.Db.Ssl_mode)
	db, err = sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	//Initilize redis information
	redisqueue.Init()

	//Fresh Start?
	freshStart := true
	if freshStart {
		_, err := redisqueue.ClearMultipleStorage([]string{"urlexists", "crawlerstatus", "messagequeue"})
		if err != nil {
			panic(err)
		}
	}

	disallowedUrls, allowed, err = htmlcrawler.LoadRobots("http://cnn.com")
	if err != nil {
		panic(err)
	}
	for _, url := range allowed {
		_, err := redisqueue.AddUniqueToQueue("urlexists", "messagequeue", url)
		if err != nil {
			fmt.Println(err)
		}
	}

	//Load Crawler Information
	info = &crawlerinformation.CrawlerInformation{}
	info.New()
	info.StoreSelf("crawlerstatus")

	go updateCrawlerStatus()
	Dispatch(10)
}

func loadConfig(path string, c *config) {
	content, err := ioutil.ReadFile(path)
	//Panic because reading config is going to be useful
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(content, c)
}

//Worker
var WorkerQueue chan chan string

type Worker struct {
	Url         chan string
	WorkerQueue chan chan string
	Commands    chan string
}

func NewWorker(workerQueue chan chan string) Worker {
	worker := Worker{
		Url:         make(chan string),
		WorkerQueue: workerQueue,
		Commands:    make(chan string),
	}
	return worker
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerQueue <- w.Url

			select {
			case url := <-w.Url:
				handleUrl(url)
				//Update Relevant Information
				//Todo maybe use a map... Concurrency or something
				mutex.Lock()
				info.UrlsCrawled++
				info.AppendArray(url)
				mutex.Unlock()
				//Sleep to slow things down...
				time.Sleep(time.Millisecond * 50)
			case command := <-w.Commands:
				fmt.Println(command)
			}
		}
	}()
}

func (w Worker) SendCommand(c string) {
	go func() {
		w.Commands <- c
	}()
}

func Dispatch(toDispatch int) {
	WorkerQueue := make(chan chan string, toDispatch)

	for i := 0; i < toDispatch; i++ {
		worker := NewWorker(WorkerQueue)
		worker.Start()
	}

	for {
		url, err := redisqueue.QueuePop("messagequeue")
		if _, ok := disallowedUrls[url]; ok || err != nil || len(url) == 0 {
			fmt.Println(len(url), err)
			fmt.Println("Sleeping..")
			time.Sleep(3 * time.Second)
		} else {
			worker := <-WorkerQueue
			worker <- url
		}
	}
}

func DispatchMessage(WorkerQueue chan chan string, message string) {
	for worker := range WorkerQueue {
		worker <- message
	}
}

//End Worker

func handleUrl(url string) {
	var urlsToAdd []string
	var err error
	switch path.Ext(strings.ToLower(url)) {
	case ".html":
		pi, err := htmlcrawler.CrawlHTML(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		urlsToAdd = pi.Urls
		err = pi.StorePage(db)
		if err != nil {
			fmt.Println(err)
		}
	case ".xml":
		urlsToAdd, err = htmlcrawler.GetXmlUrls(url)
		if err != nil {
			fmt.Println(err)
			return
		}
	default:
		return
	}
	for _, url := range urlsToAdd {
		_, err := redisqueue.AddUniqueToQueue("urlexists", "messagequeue", url)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func updateCrawlerStatus() {
	for {
		mutex.Lock()
		size, err := redisqueue.QueueLength("messagequeue")
		if err != nil {
			fmt.Println(err)
			continue
		}
		info.QueueSize = size
		indexSize, err := GetArticleCount()
		if err != nil {
			fmt.Println(err)
		}
		info.IndexSize = indexSize
		mutex.Unlock()
		info.StoreSelf("crawlerstatus")
		time.Sleep(time.Second * 15)
	}
}

func GetArticleCount() (int, error) {
	var count int
	err := db.QueryRow("SELECT count(*) FROM articles").Scan(&count)
	if err != nil {
		return -1, err
	}
	return count, nil
}
