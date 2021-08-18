package redis

import (
	"time"

	"fmt"
	"strconv"

	"urlShortner/conversion"
	"urlShortner/operation"

	"github.com/gomodule/redigo/redis"
)

type redisdb struct{ pool *redis.Pool }

const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

func New(host string, port string) (*redisdb, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		},
	}

	return &redisdb{pool}, nil
}

func (r *redisdb) doesExist(id uint64) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return exists
}

func (r *redisdb) Store(url string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()

	//var id =""
	/*
		if r.doesExist(url) {
			fmt.Println("url exists")
		} */

	shortLink := operation.Item{url, conversion.GenerateShortLink(url, UserId)}
	fmt.Println("shortlink", shortLink)

	link1, err := conn.Do("HMSET", redis.Args{"Shortener:" + conversion.GenerateShortLink(url, UserId)}.AddFlat(shortLink)...)
	if err != nil {
		fmt.Println("erro in storing")
		return "", err
	}

	fmt.Println("HMSET return", link1)

	return conversion.GenerateShortLink(url, UserId), nil
}

func (r *redisdb) Getlink(code string) (string, error) {
	conn := r.pool.Get()
	defer conn.Close()
	fmt.Println("code:", code)
	/* decodedId, err := conversion.Decode(code)
	if err != nil {
		return "", err
	}
	fmt.Println("decoded link:", decodedId) */
	urlString, err := redis.String(conn.Do("HGET", "Shortener:"+code, "url"))
	fmt.Println(urlString)
	if err != nil {
		return "", err
	} else if len(urlString) == 0 {
		return "", err
	}

	//_, err = conn.Do("HINCRBY", "Shortener:"+strconv.FormatUint(code, 10), "visits", 1)

	return urlString, nil
}

func (r *redisdb) isAvailable(id uint64) bool {
	conn := r.pool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
	if err != nil {
		return false
	}
	return !exists
}

func (r *redisdb) LoadInfo(code string) (*operation.Item, error) {
	conn := r.pool.Get()
	defer conn.Close()

	decodedId, err := conversion.Decode(code)
	if err != nil {
		return nil, err
	}

	values, err := redis.Values(conn.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10)))
	fmt.Println(values)
	if err != nil {
		return nil, err
	} else if len(values) == 0 {
		return nil, err
	}
	var shortLink operation.Item
	err = redis.ScanStruct(values, &shortLink)
	if err != nil {
		return nil, err
	}

	return &shortLink, nil
}

func (r *redisdb) Close() error {
	return r.pool.Close()
}
