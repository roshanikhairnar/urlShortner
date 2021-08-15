package main

import (
	"log"
	"urlShortner/config"
	"urlShortner/handler"
	redis "urlShortner/redisdb"

	"github.com/valyala/fasthttp"
)

func main() {
	configuration, err := config.ExtractConfiguration("config/configuration.json")
	if err != nil {
		log.Fatal(err)
	}

	service, err := redis.New(configuration.Redis.Host, configuration.Redis.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer service.Close()
	router := handler.New(configuration.Options.Schema, configuration.Options.Prefix, service)

	log.Fatal(fasthttp.ListenAndServe(":"+configuration.Server.Port, router.Handler))
	//defer service.Close()
}
