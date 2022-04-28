package main

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"time"

	u "golang.org/x/sys/unix"
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

	epfd, err := u.EpollCreate1(0)
	if err != nil {
		panic("error in creating epoll instance" + err.Error())
	}
	var event u.EpollEvent
	event.Events = syscall.EPOLLIN | syscall.EPOLLOUT
	event.Fd = int32(f.Fd())
	err = u.EpollCtl(epfd, u.EPOLL_CTL_ADD, int(event.Fd), &event)
	if err != nil {
		panic("error in epoll ctl" + err.Error())
	}

	fmt.Println(SERVER, "conn's fd:", f.Fd())
	if err != nil {
		log.Println(SERVER, "error:", err)
		return
	}
	var events [512]u.EpollEvent
	for {
		log.Println(SERVER, "looping")
		select {
		case <-s.Done:
			log.Println(SERVER, "terminating server")
			return
		default:
			nevents, err := u.EpollWait(epfd, events[:], -1)
			if err != nil {
				panic("error in waiting" + err.Error())
			}
			log.Println(SERVER, "number of events from epoll", nevents)
			for i := 0; i < nevents; i++ {
				event := events[i]
				bs := [512]byte{}
				n, err := u.Read(int(event.Fd), bs[:])
				if err != nil {
					panic("error in epoll reading" + err.Error())
				}
				log.Println(SERVER, "details:", n, string(bs[:]), event.Events, event.Fd, event.Pad)
			}
			//go func() {
			daytime := time.Now().String() + EOF
			n, err := u.Write(int(event.Fd), []byte(daytime))
			if err != nil {
				panic("error in epoll sending" + err.Error())
			}
			log.Println(SERVER, "sent", n)
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
