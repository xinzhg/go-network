package main

import (
	"log"
	"net"
	"time"
)

type Server struct {
	listener net.Listener
	Done     chan struct{}
}

func (s *Server) Shutdown() {
	err := s.listener.Close()
	if err != nil {
		panic("error in shutdown, " + err.Error())
	}
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
	s.listener = listener
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
			go func() {
				//conn.SetDeadline(time.Now().Add(3 * time.Second))
				defer conn.Close()
				daytime := time.Now().String()
				recv := [512]byte{}
				cnt, err := conn.Read(recv[:])
				if err != nil {
					log.Println("error", err)
					//conn.Close()
					return
				}
				log.Println("cnt in read server", cnt)
				cnt, err = conn.Write([]byte(daytime))
				if err != nil {
					log.Println("error", err)
					//conn.Close()
					return
				}
				log.Println("cnt in server:", cnt)
			}()
		}
	}
}
