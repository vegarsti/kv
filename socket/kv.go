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

func reader(r io.Reader, messages chan []byte) {
	buf := make([]byte, 1024)
	n, err := r.Read(buf[:])
	if err != nil {
		return
	}
	messages <- buf[:n]
}

func (k *KV) Get(key string) ([]byte, error) {
	getJSON, err := json.Marshal(Request{Kind: Get, Key: key})
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	if _, err := k.c.Write(getJSON); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	messages := make(chan []byte)
	go reader(k.c, messages)
	message := <-messages
	var response Response
	if err := json.Unmarshal(message, &response); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}
	if !response.OK {
		return nil, fmt.Errorf("got not OK response")
	}
	return []byte(response.Value), nil
}

func (k *KV) Set(key string, value []byte) error {
	setJSON, err := json.Marshal(Request{Kind: Set, Key: key, Value: string(value)})
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	if _, err := k.c.Write(setJSON); err != nil {
		return fmt.Errorf("write: %w", err)
	}
	messages := make(chan []byte)
	go reader(k.c, messages)
	message := <-messages
	var response Response
	if err := json.Unmarshal(message, &response); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	if !response.OK {
		return fmt.Errorf("got not OK response")
	}
	return nil
}
