package redisqueue

import (
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
)

type config struct {
	Port     string
	Protocol string
}

var pool *redis.Pool

func Init() {
	var err error
	pool, err = MakePool()
	if err != nil {
		panic(err)
	}
}

func MakePool() (*redis.Pool, error) {
	content, err := ioutil.ReadFile("./redisqueue/config.json")
	if err != nil {
		return &redis.Pool{}, err
	}
	var c config
	err = json.Unmarshal(content, &c)
	if err != nil {
		return &redis.Pool{}, err
	}
	pool := redis.Pool{
		MaxIdle:   50,
		MaxActive: 0,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(c.Protocol, c.Port)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
	return &pool, nil
}
