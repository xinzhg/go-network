package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

import _ "net/http/pprof"

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	server := &Server{Done: make(chan struct{}, 1)}
	go func() {
		server.Do()
	}()
	time.Sleep(1 * time.Second)
	client := &Client{URL: ":1200"}
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			time.Sleep(1 * time.Second)
			client.Do()
		}()
	}
	//client = &Client{URL: ":1200"}
	//wg := sync.WaitGroup{}
	//for i := 0; i < 10; i++ {
	//	wg.Add(1)
	//	go func() {
	//		defer func() {
	//			wg.Done()
	//			if r := recover(); r != nil {
	//				fmt.Println("Recovered in f", r)
	//			}
	//		}()
	//		time.Sleep(1 * time.Second)
	//		client.Do()
	//	}()
	//}
	log.Println("waiting")
	wg.Wait()
	//time.Sleep(1 * time.Minute)
	log.Println("before server Done")
	server.Done <- struct{}{}
	server.Shutdown()
	log.Println("after server Done")
	time.Sleep(5 * time.Second)
}
