package main

import (
	"log"
	"net"
	"time"
)

type Server struct {
	Done chan struct{}
}

func (s *Server) Do() {
	service := ":1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		panic(err)
	}
	log.Println(tcpAddr.String())
	listener, err := net.Listen("tcp", tcpAddr.String())
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-s.Done:
			log.Println("terminating server")
			return
		default:
			conn, err := listener.Accept()
			if err != nil {
				log.Println("error:", err)
				continue
			}
			daytime := time.Now().String()
			cnt, err := conn.Write([]byte(daytime))
			if err != nil {
				log.Println("error", err)
				continue
			}
			log.Println("cnt:", cnt)
			conn.Close()
		}
	}
}
