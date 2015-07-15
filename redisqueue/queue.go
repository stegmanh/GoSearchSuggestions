package redisqueue

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

func QueueLength(pool *redis.Pool, setName string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("LLEN", setName))
}

//Returns 1 if added to queue, 0 if not
func AddUniqueToQueue(pool *redis.Pool, hashName, queueName, toAdd string) (int, error) {
	exists, err := HashAdd(pool, hashName, toAdd, "true")
	if err != nil {
		return 0, err
	}
	if exists == 1 {
		QueuePush(pool, queueName, toAdd)
		return 1, nil
	}
	//Exists == 0 so the field already exists and we didnt add to the queue
	return 0, nil
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

func HashAdd(pool *redis.Pool, setName string, key string, value string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HSET", setName, key, value))
}

func HashLength(pool *redis.Pool, setName string) (int, error) {
	conn := pool.Get()
	defer conn.Close()
	return redis.Int(conn.Do("HLEN", setName))
}
