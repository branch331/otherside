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

func SetSeatgeekEvents(postCode string) error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	_, err := conn.Do("SET", postCode, "test")
	if err != nil {
		return err
	}

	return nil
}

func GetSeatgeekEvents(postCode string) error {
	if pool == nil {
		Initialize()
	}

	conn := pool.Get()
	defer conn.Close()

	s, err := redis.String(conn.Do("GET", postCode))
	if err != nil {
		return err
	}

	fmt.Printf("got em: %s\n", s)

	return nil
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
