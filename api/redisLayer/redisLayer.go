package redisLayer

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var pool *redis.Pool
var connection redis.Conn

func Initialize() {
	pool = newPool()
}

func SetKeyBytes(key string, value []byte) error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func GetKeyBytes(key string) ([]byte, error) {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Printf(err.Error())
	}

	return []byte(value), nil
}

func SetKeyString(key string, value string) error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", key, value)
	if err != nil {
		return err
	}

	return nil
}

func GetKeyString(key string) (string, error) {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	value, err := redis.String(conn.Do("GET", key))
	if err != nil {
		fmt.Printf(err.Error())
	}

	return value, nil
}

func Exists(key string) (bool, error) {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	ok, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return ok, fmt.Errorf("Error checking if key %s exists: %v", key, err)
	}

	return ok, err
}

func FlushDb() error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("FLUSHDB")
	return err
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   3,
		MaxActive: 2000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
