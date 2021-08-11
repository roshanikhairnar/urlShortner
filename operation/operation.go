package operation

import "time"

type Service interface {
  Store(string, time.Time) (string, error)
  Getlink(string) (string, error)
  //LoadInfo(string) (*Item, error)
  Close() error
}

type Item struct {
  Id      uint64 `json:"id" redis:"id"`
  URL     string `json:"url" redis:"url"`
  Expires string `json:"expires" redis:"expires"`
  Visits  int    `json:"visits" redis:"visits"`
}