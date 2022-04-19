package main

import "time"

func main() {
	server := &Server{}
	go func() {
		server.Do()
	}()
	client := &Client{URL: ":1200"}
	for i := 0; i < 100; i++ {
		client.Do()
	}
	server.Done <- struct{}{}
	time.Sleep(5 * time.Second)
}
