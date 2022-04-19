package main

import (
	"io/ioutil"
	"log"
	"net"
)

type Client struct {
	ServerAddr string
}

func (c *Client) Do() {
	if c.ServerAddr == "" {
		panic("missing domain")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", c.ServerAddr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		panic(err)
	}
	cnt, err := conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	if err != nil {
		panic(err)
	}
	log.Println("cnt:", cnt)
	res, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	log.Println(string(res))
}
