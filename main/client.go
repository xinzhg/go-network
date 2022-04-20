package main

import (
	"io/ioutil"
	"log"
	"net"
	"time"
)

type Client struct {
	URL string
}

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
	log.Println("cnt in client:", cnt)
	res, err := ioutil.ReadAll(connBackUp)
	if err != nil {
		panic(err)
	}
	log.Println(string(res))
	connBackUp.SetDeadline(time.Now().Add(-1 * time.Second))
	connBackUp.Close()
}
