package models

import (
	"github.com/garyburd/redigo/redis"
)

func HashGet(queueName, key string) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("HGET", queueName, key))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}
	return resp, nil
}
