# urlShortner

Description - url shortner API using golang redis
Installation
1. Install go https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-18-04
2. Set $GOROOT and $GOPATH
3. Install redis https://www.digitalocean.com/community/tutorials/how-to-install-and-secure-redis-on-ubuntu-20-04

Description of package
1. conversion - includes conversion of string to encoded url
2. redisdb - includes redis related functions
3. config - includes setting up configuration and extraction of same
4. handler - includes http handler function
5. main - includes main.go

Steps to run
1. start redis-server on another terminal
3. To run use following command --> go run main.go  
4. run command on other terminal -->  curl -L -X POST 'localhost:8080/encode' \
                                      -H 'Content-Type: application/json' \
                                      --data-raw '{
                                          "url": "https://www.google.com",
                                          "expires": "year-month-date hr:min:sec"
                                      }' 

Dependent packages - use go get <github packages>
1.  "github.com/gomodule/redigo/redis"
2.  "github.com/fasthttp/router"
3.  "github.com/valyala/fasthttp"



   



**docker image link** - https://hub.docker.com/r/roshanidocker/roshani_docker_repo/tags?page=1&ordering=last_updated
  
**Pull REDIS Docker image**
  
  docker pull redis
  
**Run a REDIS Docker container**
  
  docker run --name my-redis -d redis
  
  #build Docker file 
  sudo docker build -t my-urlshortner-golang .
  docker run  -p 8080:8080 --link my-redis1:redis --name my-urlshortner-golang my-urlshortner-golang
  OR
  
**Pull docker image of urlShortner**
  
  docker pull roshanidocker/roshani_docker_repo:my-urlshortner-golang
  
  docker run -it --rm -p 8080:8080 --link my-redis:redis --name urlshortener roshani_docker_repo:my-urlshortner-golang

