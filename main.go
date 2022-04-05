package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type handler struct {
}
type Age struct {
	Age int64
	mu  sync.RWMutex
}

var (
	age              *Age
	slideWindow      map[string]int64
	slideWindowMutex sync.RWMutex
	router           map[string]map[string]http.HandlerFunc
	layout           string = "2006-01-02-15-04"
	x                int64  = 2
)

func init() {
	slideWindow = make(map[string]int64)
	age = &Age{
		mu: sync.RWMutex{},
	}
}

func main() {
	router = map[string]map[string]http.HandlerFunc{
		"age": map[string]http.HandlerFunc{
			"get":  get_age,
			"post": post_age,
		},
		"car": map[string]http.HandlerFunc{
			"get": get_car,
		},
		"rate": map[string]http.HandlerFunc{
			"get": get_rate,
		},
		"buffer": map[string]http.HandlerFunc{
			"get": get_buffer,
		},
	}

	s := &http.Server{
		Addr:         ":8081",
		Handler:      new(handler),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Fatalln(s.ListenAndServe())
}

//Implement the serverHttp interface
//Routes different requests to the corresponding methods
func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uri := r.RequestURI
	method := strings.ToLower(r.Method)
	uri = strings.TrimLeft(uri, "/")
	if fns, ok := router[uri]; ok {
		if fn, ok := fns[method]; ok {
			w.Header().Add("Content-Type", "application/json")
			fn(w, r)
			return
		}
	}
	res := fmt.Sprintf("handle not found,req is  %s, method %s", uri, method)
	w.Write([]byte(res))
}
func post_age(w http.ResponseWriter, r *http.Request) {
	req := &Age{}
	err := deserialize(r.Body, req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	age.mu.Lock()
	defer age.mu.Unlock()
	age.Age = req.Age
	fmt.Println("DEBUG", "post_age", r.RemoteAddr)

	response := map[string]bool{
		"ok": true,
	}

	returnJson(w, response)
}
func get_age(w http.ResponseWriter, r *http.Request) {
	response := map[string]int64{
		"age": age.Age,
	}
	fmt.Println("DEBUG", "get_age", r.RemoteAddr)

	returnJson(w, response)
}
func get_car(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	response := map[string]int64{
		"car": now.UnixNano(),
	}
	fmt.Println("DEBUG", "get_car", r.RemoteAddr)
	add(now)
	returnJson(w, response)
}
func get_rate(w http.ResponseWriter, r *http.Request) {
	response := map[string]int64{
		"rate": getR(time.Now()),
	}
	fmt.Println("DEBUG", "get_rate", r.RemoteAddr)

	returnJson(w, response)
}
func get_buffer(w http.ResponseWriter, r *http.Request) {
	response := map[string]int64{
		"buffer": x * getR(time.Now()),
	}
	fmt.Println("DEBUG", "get_buffer", r.RemoteAddr)

	returnJson(w, response)
}

//Parsing the parameters of a json request
func deserialize(body io.Reader, age *Age) error {
	bodyByte, err := ioutil.ReadAll(body)
	if err != nil {

		return err
	}
	err = json.Unmarshal(bodyByte, age)
	return err
}

//Serialize the result to json and return
func returnJson(w http.ResponseWriter, data interface{}) {
	res, err := json.Marshal(data)
	if err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusBadRequest)
	}
	w.Write(res)
}

// Increase car sales data
func add(now time.Time) {
	slideWindowMutex.Lock()
	defer slideWindowMutex.Unlock()
	nowStr := now.Format(layout)
	fmt.Println("DEBUG", "add", fmt.Sprintf("time is %s", nowStr))
	if _, ok := slideWindow[nowStr]; ok {
		slideWindow[nowStr]++
	} else {
		slideWindow[nowStr] = 1
	}

}

// Get the total number of cars sold within 60 minutes
func getR(now time.Time) int64 {
	slideWindowMutex.RLock()
	defer slideWindowMutex.RUnlock()
	var r int64
	for i := 0; i < 60; i++ {
		timeStr := now.Add(time.Duration(0-i*60) * time.Second).Format(layout)
		if num, ok := slideWindow[timeStr]; ok {
			r = r + num
		}
	}
	return r

}
