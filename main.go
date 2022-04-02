package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type handler struct {
}
type Age struct {
	age int64
	mu  sync.RWMutex
}

var age *Age
var router map[string]map[string]http.HandlerFunc

func main() {
	router = map[string]map[string]http.HandlerFunc{
		"age": map[string]http.HandlerFunc{
			"get":  get_age,
			"post": post_age,
		},
	}

	age = &Age{
		mu: sync.RWMutex{},
	}
	s := &http.Server{
		Addr:         ":8081",
		Handler:      new(handler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatalln(s.ListenAndServe())
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	log.Println(r.Form)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	uri := r.RequestURI
	method := strings.ToLower(r.Method)
	uri = strings.TrimLeft(uri, "/")
	if fns, ok := router[uri]; ok {
		if fn, ok := fns[method]; ok {
			fn(w, r)
			return
		}
	}
	res := fmt.Sprintf("handle not found,req is  %s, method %s", uri, method)
	w.Write([]byte(res))
}
func post_age(w http.ResponseWriter, r *http.Request) {
	age := r.FormValue("age")
	log.Println(r.Form)
	w.Write([]byte(age))
}
func get_age(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("get age"))
}

