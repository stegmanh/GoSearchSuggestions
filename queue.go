package trie

import (
	"errors"
	"github.com/garyburd/redigo/redis"
)

func QueuePop(pool *redis.Pool, queueName string) (string, error) {
	conn := pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("RPOP", queueName))
	if err != nil {
		if err == redis.ErrNil {
			return "", nil
		}
		return "", err
	}
	return resp, nil
}

func QueuePush(pool *redis.Pool, queueName string, value string) error {
	conn := pool.Get()
	defer conn.Close()
	if len(value) == 0 {
		return errors.New("Cannot input an empty string")
	}
	_, err := conn.Do("LPUSH", queueName, value)
	if err != nil {
		return err
	}
	return nil
}

func HashExists(pool *redis.Pool, setName string, value string) (bool, error) {
	conn := pool.Get()
	defer conn.Close()
	resp, err := redis.Int(conn.Do("HEXISTS", setName, value))
	if err != nil {
		return false, err
	}
	if resp == 1 {
		return true, nil
	}
	return false, nil
}

func CreatePool() redis.Pool {
	pool := redis.Pool{
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
	return pool
}
