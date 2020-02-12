package main

import (
	"log"

	"github.com/bradfitz/gomemcache/memcache"
)

func main() {
	log.Println("smoke testing memcached...")
	mc := memcache.New("localhost:11211")
	mc.Set(&memcache.Item{Key: "test-key", Value: []byte("test-value")})

	it, err := mc.Get("test-key")
	if err != nil {
		log.Fatalln("unable to access memcached", err)
	}
	log.Println("success!", it.Key, string(it.Value))
	mc.Delete("test-key")
	log.Println("done.")
}
