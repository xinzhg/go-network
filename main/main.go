package main

func main() {
	client := &Client{URL: "google.com:80"}
	client.Do()
}
