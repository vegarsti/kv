package embedded

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type KV struct {
	file string
	m    map[string]string
}

func New() *KV {
	return &KV{file: "/tmp/kv.db", m: make(map[string]string)}
}

func (k *KV) Open() error {
	if _, err := os.OpenFile(k.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	if err := k.readFile(); err != nil {
		return err
	}
	return nil
}

func (k *KV) Close() error {
	k.flush()
	return nil
}

func (k *KV) Get(key string) (string, error) {
	v, ok := k.m[key]
	if !ok {
		return "", fmt.Errorf("not found")
	}
	return v, nil
}

func (k *KV) Put(key string, value string) error {
	return nil
}

func (k *KV) Delete(key string) error {
	return nil
}

func (k *KV) flush() error {
	return nil
}

func (k *KV) readFile() error {
	file, err := os.Open(k.file)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}
	defer file.Close()
	s := bufio.NewScanner(file)
	i := 0
	for s.Scan() {
		line := s.Text()
		words := strings.Split(line, " ")
		if len(words) != 2 {
			return fmt.Errorf("line %d: expected 2 words, got %d", i, len(words))
		}
		key := words[0]
		value := words[1]
		k.m[key] = value
		i++
	}
	return nil
}
