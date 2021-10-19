package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const version = "version"



type httpserver struct {
	ip string
	port int
	version string
	routes map[string]func(w http.ResponseWriter,r *http.Request)
}
func (s httpserver) setResponseHeader(w http.ResponseWriter, r *http.Request) {
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			w.Header().Set(k, strings.Join(v, ","))
		}
	}
	w.Header().Set(version, os.Getenv(version))
}
func (s httpserver) init() {
	log.Printf("starting server on port %d\n", s.port)
	// 设置环境变量
	err := os.Setenv(version, s.version)
	if err != nil {
		log.Printf("set environment variable %s err\n", version)
	}

	mux := http.NewServeMux()

	if s.routes == nil || len(s.routes) == 0 {
		mux.HandleFunc("/", s.IndexHandler)
	} else {
		for route, handler := range s.routes {
			mux.HandleFunc(route, handler)
		}
	}
	err = http.ListenAndServe(":"+strconv.Itoa(s.port), mux)
	if err != nil {
		log.Println("server start fail ", err.Error())
	}
}
func (s httpserver) IndexHandler(w http.ResponseWriter, r *http.Request) {
	s.setResponseHeader(w, r)
	if len(r.Header) > 0 {
		for k, v := range r.Header {
			_, err := io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
			if err != nil {
				log.Printf("The Handler Errot")
			}
		}
	}
}
func (s httpserver) CheckHealth(w http.ResponseWriter,r *http.Request) {
	s.setResponseHeader(w,r)
	_, err := io.WriteString(w, "200")
	if err != nil {
		log.Printf("health check error")
	}
}


func main() {
	server := httpserver{
		ip:      "127.0.0.1",
		port:    8080,
		version: "1.1.1",
	}
	routes := make(map[string]func(http.ResponseWriter, *http.Request))
	routes["/index"] = server.IndexHandler
	routes["/health"] = server.CheckHealth
	server.routes = routes
	server.init()
}