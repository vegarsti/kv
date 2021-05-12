package socket

import (
	"fmt"
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

func (k *KV) Get(key string) ([]byte, error) {
	if _, err := k.c.Write([]byte(key)); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}
	return []byte(`hello world`), nil
}

func (k *KV) Set(key string, value []byte) error {
	return nil
}
