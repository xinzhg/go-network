package main

import (
	"fmt"
	"time"
)

func main() {
	server := &Server{}
	go func() {
		server.Do()
	}()
	client := &Client{URL: ":1200"}
	for i := 0; i < 1; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered in f", r)
				}
			}()
			time.Sleep(1 * time.Second)
			client.Do()
		}()
	}
	server.Done <- struct{}{}
	time.Sleep(5 * time.Second)
}
