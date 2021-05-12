package socket

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

type Server struct {
	addr string
	l    net.Listener
	quit chan interface{}
	wg   sync.WaitGroup
}

func NewServer() *Server {
	return &Server{
		addr: "/tmp/echo.sock",
		quit: make(chan interface{}),
	}
}

func (s *Server) Open() error {
	if err := os.RemoveAll(s.addr); err != nil {
		return fmt.Errorf("remove all: %w", err)
	}
	l, err := net.Listen("unix", s.addr)
	if err != nil {
		return fmt.Errorf("listen error: %w", err)
	}
	s.l = l
	s.wg.Add(1)
	go s.serve()
	return nil
}

func (s *Server) Close() error {
	close(s.quit)
	s.l.Close()
	log.Println("Waiting for connections to close")
	s.wg.Wait()
	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
			break
		}
		if n == 0 {
			break
		}
		log.Printf("Received %s", string(buf[:n]))
	}
	log.Println("Connection closed")
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.l.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				log.Printf("accept error: %v", err)
			}
			continue
		}
		s.wg.Add(1)
		go func() {
			log.Println("Got connection")
			s.handleConnection(conn)
			s.wg.Done()
		}()
	}
}

func (s *Server) Address() string {
	return s.addr
}
