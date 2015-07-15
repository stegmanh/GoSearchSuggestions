package crawlerinformation

import (
	"../redisqueue"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
)

type CrawlerInformation struct {
	Status string //Running, stopped, idle
	// Cpu         string            //Cpu usage -- Not doing yet
	// Ram         string            //Ram usage -- Not doing yet
	UrlsCrawled int               //# Of total urls crawled
	LastCrawled []string          //Last of some arbitrary number
	QueueSize   int               //Size of the queue
	IndexSize   int               //Size of DB
	Errors      map[string]string //Map of all errors (Maybe we dont use)
}

func (c *CrawlerInformation) New() {
	c.Status = "Running"
	c.UrlsCrawled = 0
	c.LastCrawled = make([]string, 0)
	c.QueueSize = 0
	c.IndexSize = 0
	c.Errors = make(map[string]string)
}

//I dont think this should exist. Crawler should be about current crawl... Maybe store at shutdown?
func (c *CrawlerInformation) StoreSelf(pool *redis.Pool, hashName string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	_, err = redisqueue.HashAdd(pool, hashName, "data", string(data))
	return err
}

func (c *CrawlerInformation) AppendArray(value string) {
	//Hard coded here
	if len(c.LastCrawled) > 9 {
		temp := []string{value}
		c.LastCrawled = append(temp, c.LastCrawled[1:]...)
	} else {
		c.LastCrawled = append(c.LastCrawled, value)
	}
}
