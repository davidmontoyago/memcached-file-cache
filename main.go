package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/davidmontoyago/interview-davidmontoyago-d660952eff664d8bac96c9124d7f8582/pkg/filecache"
	"github.com/gorilla/mux"
)

const memcachedAddr = "localhost:11211"

var memcachedClient *memcache.Client

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/filecache", upload).Methods("POST")
	router.HandleFunc("/filecache/{key}", download).Methods("GET")

	memcachedClient = memcache.New(memcachedAddr)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func upload(res http.ResponseWriter, req *http.Request) {
	file, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(res, toJSONError(err), http.StatusBadRequest)
		return
	}
	fc := filecache.New(memcachedClient)
	key, err := fc.Put(file)
	if err != nil {
		log.Println("failed to put file:", err)
		http.Error(res, toJSONError(err), http.StatusBadRequest)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{"ok": "true", "key": key})
	log.Printf("success! uploaded file %s\n", key)
}

func download(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileKey := vars["key"]

	fc := filecache.New(memcachedClient)
	file, err := fc.Get(fileKey)
	if err != nil {
		log.Println("failed to get file:", err)
		http.Error(res, toJSONError(err), http.StatusNotFound)
		return
	}

	http.ServeContent(res, req, fmt.Sprintf("%s.dat", fileKey), time.Now(), bytes.NewReader(file))
	log.Printf("success! downloaded file %s\n", fileKey)
}

func toJSONError(err error) string {
	errorResponse := new(bytes.Buffer)
	json.NewEncoder(errorResponse).Encode(map[string]string{"ok": "false", "error": err.Error()})
	return errorResponse.String()
}
