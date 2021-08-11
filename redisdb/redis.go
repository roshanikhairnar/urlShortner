package redisdb

import(
	"github.com/gomodule/redigo"
   "../operation"
   "../conversion"

)
//var pool *redisclient.Pool
type redis struct{ pool *redisclient.Pool }

func InitPool(host, port, password string) (operation.Service, error) {
  connectionectionType := "tcp"
  pool := &redisclient.Pool{
     MaxIdle:     10,
     IdleTimeout: 240 * time.Second,
     Dial: func() (redisclient.connection, error) {
        return redisclient.Dial(connectionectionType, fmt.Sprintf("%s:%s", host, port))
     },
  }

  return &redis{pool}, nil
}

func (r *redis) doesExists(id uint64) bool {
   connection := r.pool.Get()
   defer connection.Close()
 
   exists, err := redisclient.Bool(connection.Do("EXISTS", "Shortener:"+strconv.FormatUint(id, 10)))
   if err != nil {
      return false
   }
   return exists
 }
 
 
 func (r *redis) Store(url string, expires time.Time) (string, error) {
   connection := r.pool.Get()
   defer connection.Close()
 
   var id uint64
 
   for used := true; used; used = r.doesExists(id) {
      id = rand.Uint64()
   }
 
   shortLink := operation.Item{id, url, expires.Format("2006-01-02 15:04:05.728046 +0300 EEST"), 0}
 
   _, err := connection.Do("HMSET", redisclient.Args{"Shortener:" + strconv.FormatUint(id, 10)}.AddFlat(shortLink)...)
   if err != nil {
      return "", err
   }
 
   _, err = connection.Do("EXPIREAT", "Shortener:"+strconv.FormatUint(id, 10), expires.Unix())
   if err != nil {
      return "", err
   }
 
   return conversion.Encode(id), nil
 }
 func (r *redis) Getlink(code string) (string, error) {
   connection := r.pool.Get()
   defer connection.Close()
 
   decodedId, err := conversion.Decode(code)
   if err != nil {
      return "", err
   }
 
   urlString, err := redisclient.String(connection.Do("HGET", "Shortener:"+strconv.FormatUint(decodedId, 10), "url"))
   if err != nil {
      return "", err
   } else if len(urlString) == 0 {
      return "", storage.ErrNoLink
   }
 
   _, err = connection.Do("HINCRBY", "Shortener:"+strconv.FormatUint(decodedId, 10), "visits", 1)
 
   return urlString, nil
 }
 
 
 
 func (r *redis) Close() error {
   return r.pool.Close()
 }