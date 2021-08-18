package operation

type Service interface {
	Store(string) (string, error)
	Getlink(string) (string, error)
	LoadInfo(string) (*Item, error)
	Close() error
}

type Item struct {
	URL      string `json:"url" redis:"url"`
	ShortURl string `json:"ShortUrl" redis:"ShortUrl"`
}
