package main

import (
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
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
	epfd, err := u.EpollCreate1(u.EPOLL_CLOEXEC)
	if err != nil {
		panic("error in creating epoll instance" + err.Error())
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			log.Println(SERVER, "received connection from client")
			if err != nil {
				fmt.Println(SERVER, "any error:", err.Error())
				//if errors.Is(err, u)
				if strings.Contains(err.Error(), "use of closed network connection") {
					defer func() {
						log.Println("finishing sever")
						err := u.Close(epfd)
						if err != nil {
							panic("error in closing epfd:" + err.Error())
						}

					}()
					return
				}
				panic("error in accepting tcp" + err.Error())
			}
			f, _ := conn.File()
			var event u.EpollEvent
			event.Events = u.EPOLLIN // | u.EPOLLOUT | u.EPOLLMSG
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
			log.Println(SERVER, "next for loop")
		}
	}()

	cnt := 0
	m := sync.Mutex{}

	for {
		//log.Println(SERVER, "looping")
		select {
		case <-s.Done:
			log.Println(SERVER, "terminating server")
			return
		default:
			var events [512]u.EpollEvent
			nevents, err := u.EpollWait(epfd, events[:], 100)
			if err != nil {
				panic("error in waiting" + err.Error())
			}
			if nevents == 0 {
				continue
			}
			log.Println(SERVER, "number of events from epoll", nevents)
			efd := 0
			for i := 0; i < nevents; i++ {
				event := events[i]
				bs := [512]byte{}
				//log.Println(SERVER, "reading", event.Fd)
				n, err := u.Read(int(event.Fd), bs[:])
				efd = int(event.Fd)
				if err != nil {
					panic("error in epoll reading" + err.Error())
				}
				m.Lock()
				cnt++
				m.Unlock()
				log.Println(SERVER, "details:", string(bs[:]), "read size:", n, "type:", event.Events, "fd:", event.Fd, "pad", event.Pad, "cnt:", cnt)
			}
			//go func() {
			daytime := time.Now().String() + EOF
			_, err = u.Write(int(efd), []byte(daytime))
			if err != nil {
				panic("error in epoll sending" + err.Error())
			}
			//log.Println(SERVER, "sent", n)
			if err != nil {
				log.Println(SERVER, "error", err)
				//conn.Close()
				return
			}
			//var event u.EpollEvent
			//event.Events = syscall.EPOLLIN | syscall.EPOLLOUT | syscall.EPOLL_CTL_ADD
			//event.Fd = int32(efd)
			//u.EpollCtl(epfd, u.EPOLL_CTL_DEL, efd, &event)
			//u.Close(efd)
			//conn.Close()
			//}()
		}
	}
}
