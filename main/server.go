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
	for {
		log.Println(SERVER + "looping")
		select {
		case <-s.Done:
			log.Println(SERVER + "terminating server")
			return
		default:
			log.Println(SERVER + "before accept")
			conn, err := listener.AcceptTCP()
			log.Println(SERVER + "after accept")
			if err != nil {
				log.Println(SERVER+"error:", err)
				continue
			}
			go func() {
				//conn.SetDeadline(time.Now().Add(3 * time.Second))
				//defer conn.Close()
				daytime := time.Now().String() + EOF
				recv := [512]byte{}
				cnt, err := conn.Read(recv[:])
				if err != nil {
					log.Println(SERVER+"error", err)
					//conn.Close()
					return
				}
				fmt.Println(SERVER, "what server recev:", string(recv[:]), len(recv[:]), len(recv))
				log.Println(SERVER+"cnt in read server", cnt)
				//conn.SetNoDelay(true)
				log.Println(SERVER, "before io.copy")
				//io.Copy(ioutil.Discard, conn)
				log.Println(SERVER, "after io.copy")
				log.Println(SERVER, "before write")
				log.Println(SERVER, "before write")
				fd, err := conn.File()
				cnt, err = fd.Write([]byte(daytime))
				err = fd.Sync()
				//cnt, err = conn.Write([]byte(daytime))
				log.Println(SERVER, "after write")
				if err != nil {
					log.Println(SERVER+"error", err)
					//conn.Close()
					return
				}
				//conn.Close()
				log.Println(SERVER+"cnt in server:", cnt)
			}()
		}
	}
}
