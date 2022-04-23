package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const EOF = "\000"

const SERVER = "service side:"

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
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		panic(err)
	}
	s.listener = listener
	conn, err := listener.AcceptTCP()
	f, _ := conn.File()
	fmt.Println(SERVER, "conn's fd:", f.Fd())
	if err != nil {
		log.Println(SERVER, "error:", err)
		return
	}
	for {
		log.Println(SERVER, "looping")
		select {
		case <-s.Done:
			log.Println(SERVER, "terminating server")
			return
		default:
			//go func() {
			daytime := time.Now().String() + EOF
			recv := [512]byte{}
			_, err := conn.Read(recv[:])
			if err != nil {
				log.Println(SERVER, "error", err)
				//conn.Close()
				return
			}
			fmt.Println(SERVER, "what server recevs:", string(recv[:]), len(recv[:]), len(recv))
			_, err = conn.Write([]byte(daytime))
			log.Println(SERVER, "sent")
			if err != nil {
				log.Println(SERVER, "error", err)
				//conn.Close()
				return
			}
			//conn.Close()
			//}()
		}
	}
}
