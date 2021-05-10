package kv

type KV interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}
