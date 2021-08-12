package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strconv"
	"github.com/router"
	//"github.com/gomodule/redigo/redis"
	"github.com/valyala/fasthttp"
	//"github.com/gorilla/mux"
	"../operation"
	"time"
	"net/url"
)

type response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"shortUrl"`
  }

type handler struct {
	schema  string
	host    string
	operation operation.Service
	
}

func New(schema string, host string, operation operation.Service) *router.Router {
	router := router.New()
  
	h := handler{schema, host, operation}
	router.POST("/encode/", responseHandler(h.encode))
	router.GET("/{shortLink}", h.redirect)
	router.GET("/{shortLink}/info", responseHandler(h.decode))
	return router
  }

  func responseHandler(h func(ctx *fasthttp.RequestCtx) (interface{}, int, error)) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
	   data, status, err := h(ctx)
	   fmt.Println(data)
	   
	   if err != nil {
		  data = err.Error()
	   }
	   ctx.Response.Header.Set("Content-Type", "application/json")
	   ctx.Response.SetStatusCode(status)
	   err = json.NewEncoder(ctx.Response.BodyWriter()).Encode(response{Data: data, Success: err == nil})
	   if err != nil {
		  log.Printf("could not encode response to output: %v", err)
	   }
	}
  }
  
  func (h handler) encode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	var input struct {
	   URL     string `json:"url"`
	   Expires string `json:"expires"`
	}
  
	if err := json.Unmarshal(ctx.PostBody(), &input); err != nil {
	   return nil, http.StatusBadRequest, fmt.Errorf("Unable to decode JSON request body: %v", err)
	}
  
	uri, err := url.ParseRequestURI(input.URL)
  
	if err != nil {
	   return nil, http.StatusBadRequest, fmt.Errorf("Invalid url")
	}
  
	layoutISO := "2006-01-02 15:04:05"
	expires, err := time.Parse(layoutISO, input.Expires)
	if err != nil {
	   return nil, http.StatusBadRequest, fmt.Errorf("Invalid expiration date")
	}
  
	c, err := h.operation.Store(uri.String(), expires)
	if err != nil {
	   return nil, http.StatusInternalServerError, fmt.Errorf("Could not store in database: %v", err)
	}
  
	u := url.URL{
	   Scheme: h.schema,
	   Host:   h.host,
	   Path:   c}
  
	fmt.Printf("Generated link: %v \n", u.String())
  
	return u.String(), http.StatusCreated, nil
  }
  
  func (h handler) decode(ctx *fasthttp.RequestCtx) (interface{}, int, error) {
	code := ctx.UserValue("shortLink").(string)
  
	model, err := h.operation.LoadInfo(code)
	if err != nil {
	   return nil, http.StatusNotFound, fmt.Errorf("URL not found")
	}
  
	return model, http.StatusOK, nil
  }
  
  func (h handler) redirect(ctx *fasthttp.RequestCtx) {
	code := ctx.UserValue("shortLink").(string)
  
	uri, err := h.operation.Getlink(code)
	if err != nil {
	   ctx.Response.Header.Set("Content-Type", "application/json")
	   ctx.Response.SetStatusCode(http.StatusNotFound)
	   return
	}
  
	ctx.Redirect(uri, http.StatusMovedPermanently)
  }
   

  