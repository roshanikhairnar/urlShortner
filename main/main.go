package main

import (
  "log"
  "../config"
  "../handler"
  "../redisdb"
  "github.com/valyala/fasthttp"
)

func main() {
  configuration, err := config.ExtractConfiguration("../config/configuration.json")
  if err != nil {
     log.Fatal(err)
  } 

  service,err := redis.New(configuration.Redis.Host, configuration.Redis.Port, configuration.Redis.Password)
  if err != nil {
     log.Fatal(err)
  }
  router := handler.New(configuration.Options.Schema, configuration.Options.Prefix, service)

  log.Fatal(fasthttp.ListenAndServe(":" + configuration.Server.Port, router.Handler))
  //defer service.Close()
}