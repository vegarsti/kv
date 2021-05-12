package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
)

type KV struct {
	addr string
	c    net.Conn
}

func New() *KV {
	return &KV{addr: "/tmp/echo.sock"}
}

func (k *KV) Address() string {
	return k.addr
}

func (k *KV) Open() error {
	c, err := net.Dial("unix", k.addr)
	if err != nil {
		return fmt.Errorf("dial: %w", err)
	}
	k.c = c
	return nil
}

func (k *KV) Close() error {
	return k.c.Close()
}

func receive(r io.Reader, messages chan []byte) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		messages <- []byte(err.Error())
	}
	messages <- buf[:n]
}

func (k *KV) Get(key string) ([]byte, error) {
	request := Request{Kind: Get, Key: key}
	response, err := k.Request(request)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	if !response.OK {
		return nil, fmt.Errorf("got not OK response")
	}
	return []byte(response.Value), nil
}

func (k *KV) Set(key string, value []byte) error {
	request := Request{Kind: Set, Key: key, Value: string(value)}
	response, err := k.Request(request)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	if !response.OK {
		return fmt.Errorf("got not OK response")
	}
	return nil
}

func (k *KV) Request(request Request) (*Response, error) {
	messages := make(chan []byte)
	go receive(k.c, messages)

	setJSON, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	if _, err := k.c.Write(setJSON); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	message := <-messages
	if string(message) == "error" {
		return nil, fmt.Errorf("got error message in channel: %v", string(message))
	}
	var response Response
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	return &response, nil
}
