package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	server := &Server{Done: make(chan struct{}, 1)}
	go func() {
		server.Do()
	}()
	client := &Client{URL: ":1200"}
	for i := 0; i < 10; i++ {
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			time.Sleep(1 * time.Second)
			client.Do()
		}()
	}
	log.Println("before server Done")
	server.Done <- struct{}{}
	server.Shutdown()
	log.Println("after server Done")
	time.Sleep(5 * time.Second)
}
