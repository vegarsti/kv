package main

import (
	"kv/socket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	s := socket.NewServer("/tmp/kv.sock")
	if err := s.Open(); err != nil {
		log.Printf("%v", err)
		os.Exit(1)
	}
	defer s.Close()
	log.Printf("Listening on %s", s.Address())
	var interruptChannel = make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
	sig := <-interruptChannel
	log.Println("Got interrupt signal: ", sig.String())
	log.Println("Shutting down server")
}
