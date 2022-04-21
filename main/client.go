package main

import (
	"io/ioutil"
	"log"
	"net"
)

type Client struct {
	URL string
}

const CLIENT = "client side: "

var connBackUp *net.TCPConn

func (c *Client) Do() {
	if c.URL == "" {
		panic("missing domain")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.URL)
	if err != nil {
		panic(err)
	}
	if connBackUp == nil {
		connBackUp, err = net.DialTCP("tcp", nil, tcpAddr)
	}
	if err != nil {
		panic(err)
	}
	cnt, err := connBackUp.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	if err != nil {
		panic(err)
	}
	log.Println(CLIENT, "cnt in client:", cnt)
	log.Println(CLIENT, "before readAll in client")
	//connBackUp.SetNoDelay(true)
	res, err := ioutil.ReadAll(connBackUp)
	if err != nil {
		panic(err)
	}
	log.Println(CLIENT, "after readAll in client")
	log.Println(CLIENT, string(res))
	//connBackUp.SetDeadline(time.Now().Add(-1 * time.Second))
	log.Println(CLIENT, "before close in client")
	//connBackUp.Close()
	log.Println(CLIENT, "after close in client")
}
