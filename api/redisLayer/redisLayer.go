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

func SetSeatgeekEvents(postCode string, seatGeekEvents []byte) error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", postCode, seatGeekEvents)
	if err != nil {
		return err
	}

	return nil
}

func GetSeatgeekEvents(postCode string) ([]byte, error) {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	seatGeekEvents, err := redis.String(conn.Do("GET", postCode))
	if err != nil {
		fmt.Printf(err.Error())
	}

	return []byte(seatGeekEvents), nil
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

func Ping() {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	pong, err := conn.Do("PING")

	if err != nil {
		fmt.Printf(err.Error())
	}

	s, err := redis.String(pong, err)
	if err != nil {
		fmt.Printf(err.Error())
	}

	fmt.Printf("PING resp: %s\n", s)
}

func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   3,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
