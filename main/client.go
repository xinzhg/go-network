package main

import (
	"log"
	"net"
	"strconv"
	"sync"
)

type Client struct {
	URL string
}

const CLIENT = "client side: "

var connBackUp *net.TCPConn

var seq = 0
var m = sync.Mutex{}

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
	m.Lock()
	seq++
	m.Unlock()
	_, err = connBackUp.Write([]byte("info " + strconv.Itoa(seq)))
	log.Println(CLIENT, "sent msg")
	if err != nil {
		panic(err)
	}
	//log.Println(CLIENT, "cnt in client:", cnt)
	//log.Println(CLIENT, "before readAll in client")
	//connBackUp.SetNoDelay(true)
	//res, err := ioutil.ReadAll(connBackUp)
	res := [512]byte{}
	log.Println(CLIENT, "before reading")
	_, err = connBackUp.Read(res[:])
	log.Println(CLIENT, "after reading")
	if err != nil {
		panic(err)
		return
	}
	if err != nil {
		panic(err)
	}
	log.Println(CLIENT, string(res[:]))

	//log.Println(CLIENT, "after readAll in client")
	//connBackUp.SetDeadline(time.Now().Add(-1 * time.Second))
	//log.Println(CLIENT, "before close in client")
	//connBackUp.Close()
	//log.Println(CLIENT, "after close in client")
}
