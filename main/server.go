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
			log.Println("before accept")
			conn, err := listener.Accept()
			log.Println("after accept")
			if err != nil {
				log.Println("error:", err)
				continue
			}
			daytime := time.Now().String()
			recv := make([]byte, 512)
			cnt, err := conn.Read(recv)
			if err != nil {
				log.Println("error", err)
				conn.Close()
				continue
			}
			log.Println("cnt in read server", cnt)
			cnt, err = conn.Write([]byte(daytime))
			if err != nil {
				log.Println("error", err)
				conn.Close()
				continue
			}
			log.Println("cnt in server:", cnt)
			conn.Close()
		}
	}
}
