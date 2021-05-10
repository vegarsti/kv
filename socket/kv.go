package socket

type KV struct {
	addr string
}

func New() *KV {
	return &KV{addr: "/tmp/echo.sock"}
}

func (k *KV) Address() string {
	return k.addr
}

func (k *KV) Get(key string) ([]byte, error) {
	return []byte(`hello world`), nil
}

func (k *KV) Set(key string, value []byte) error {
	return nil
}
