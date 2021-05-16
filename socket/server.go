package socket

import (
	"encoding/json"
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
	kv   map[string]string
}

func NewServer() *Server {
	return &Server{
		addr: "/tmp/echo.sock",
		quit: make(chan interface{}),
		kv:   make(map[string]string),
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
	buf := make([]byte, 2048) // max size of request payload
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
			break
		}
		if n == 0 {
			break
		}
		got := buf[:n]
		log.Printf("Received %s", string(got))
		var request Request
		if err := json.Unmarshal(got, &request); err != nil {
			log.Printf("Unmarshal request error: %v", err)
			break
		}
		var response Response
		switch request.Kind {
		case Get:
			log.Printf("GET %s", request.Key)
			response.Kind = Get
			value, ok := s.kv[request.Key]
			response.OK = ok
			response.Value = value
		case Set:
			log.Printf("SET %s %s", request.Key, request.Value)
			response.Kind = Set
			_, ok := s.kv[request.Key]
			response.OK = !ok
			if !ok {
				s.kv[request.Key] = request.Value
			}
		case Delete:
			log.Printf("DELETE %s", request.Key)
			response.Kind = Delete
			_, ok := s.kv[request.Key]
			response.OK = ok
			if ok {
				delete(s.kv, request.Key)
			}
		default:
			log.Println("method not allowed")
		}
		responseJSON, err := json.Marshal(response)
		if err != nil {
			log.Printf("Marshal response error: %v", err)
			break
		}
		if _, err := conn.Write(responseJSON); err != nil {
			log.Printf("write error: %v", err)
			break
		}
		log.Printf("Sent %s", string(responseJSON))
	}
	log.Println("Connection closed")
	log.Printf("Currently %d keys", len(s.kv))
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
