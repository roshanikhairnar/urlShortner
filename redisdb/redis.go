package redis

import (
	
	"time"

	"github.com/gomodule/redigo/redis"
   "../operation"
   "fmt"
   "strconv"
   "math/rand"
   "../conversion"
	
)



type redisdb struct{ pool *redis.Pool }

func New(host, port, password string) (*redisdb,error){
  pool := &redis.Pool{
     MaxIdle:     10,
     IdleTimeout: 240 * time.Second,
     Dial: func() (redis.Conn, error) {
        return redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
     },
  }

  return &redisdb{pool}, nil
}


func (r *redisdb) isUsed(id uint64) bool {
   conn := r.pool.Get()
   defer conn.Close()
 
   exists, err := redis.Bool(conn.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
   if err != nil {
      return false
   }
   return exists
 }
 
 
 func (r *redisdb) Store(url string, expires time.Time) (string, error) {
   conn := r.pool.Get()
   defer conn.Close()
 
   var id uint64
 
   for used := true; used; used = r.isUsed(id) {
      id = rand.Uint64()
   }
 
   shortLink := operation.Item{id, url, expires.Format("2006-01-02 15:04:05.728046 +0300 EEST"), 0}
 
   _, err := conn.Do("HMSET", redis.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)
   if err != nil {
      return "", err
   }
 
   _, err = conn.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())
   if err != nil {
      return "", err
   }
 
   return conversion.Encode(id), nil
 }

func (r *redisdb) Getlink(code string) (string, error) {
   conn := r.pool.Get()
   defer conn.Close()
 
   decodedId, err := conversion.Decode(code)
   if err != nil {
      return "", err
   }
 
   urlString, err := redis.String(conn.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10), "url"))
   if err != nil {
      return "", err
   } else if len(urlString) == 0 {
      return "", err
   }
 
   _, err = conn.Do("HINCRBY", "Shortener:"+strconv.FormatUint(decodedId, 10), "visits", 1)
 
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
 
  
 func (r *redisdb) Close() error {
   return r.pool.Close()
 }