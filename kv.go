package kv

type KV interface {
	Get(key string) (string, error)
	Put(key string, value string) error
	Delete(key string) error
}
