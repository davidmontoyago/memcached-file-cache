package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/filecache", upload).Methods("PUT")
	router.HandleFunc("/filecache/{key}", download).Methods("GET")

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
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("file len", len(file))
}

func download(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	fileKey := vars["key"]
	log.Println(fileKey)
}
