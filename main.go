package main

import (
	"fmt"
	"log"
	"main/cache"
	"main/server"
	"net/http"
)

type String string

func (s String) Len() int {
	return len(s)
}

var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}
//
//
//func main() {
//	loadCounts := make(map[string]int, len(db))
//	gee := cache.NewGroup("scores", 2<<10, cache.GetterFunc(
//		func(key string) ([]byte, error) {
//			log.Println("[SlowDB] search key", key)
//			if v, ok := db[key]; ok {
//				if _, ok := loadCounts[key]; !ok {
//					loadCounts[key] = 0
//				}
//				loadCounts[key] += 1
//				return []byte(v), nil
//			}
//			return nil, fmt.Errorf("%s not exist", key)
//		}))
//
//	for k, v := range db {
//		if view, err := gee.Get(k); err != nil || view.String() != v {
//			log.Fatalln("failed to get value of Tom")
//		} // load from callback function
//
//		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
//			log.Fatalf("cache %s miss\n", k)
//		} // cache hit
//	}
//}

func main() {
	cache.NewGroup("scores", 2<<10, cache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := server.New(addr)
	log.Println("bincache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
