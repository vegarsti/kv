package bolt

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

type KV struct {
	file       string
	db         *bolt.DB
	bucketName []byte
}

func New() *KV {
	return &KV{file: "/tmp/kv-bolt.db"}
}

func (k *KV) Open() error {
	db, err := bolt.Open(k.file, 0600, nil)
	if err != nil {
		return fmt.Errorf("bolt open: %w", err)
	}
	k.bucketName = []byte("keys")
	// Create the bucket if it doesn't exist
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists(k.bucketName); err != nil {
			return fmt.Errorf("create bucket if not exists: %w", err)
		}
		return nil
	}); err != nil {
		log.Printf("update: %v", err.Error())
	}
	// set db on struct
	k.db = db
	return nil
}

func (k *KV) Close() error {
	return k.db.Close()
}

func (k *KV) Get(key string) (string, error) {
	var value string
	if err := k.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(k.bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket was nil")
		}
		val := bucket.Get([]byte(key))
		if val == nil {
			return fmt.Errorf("value was nil")
		}
		value = string(val)
		return nil
	}); err != nil {
		return "", fmt.Errorf("update: %w", err)
	}
	return value, nil
}

func (k *KV) Put(key string, value string) error {
	if err := k.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(k.bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket was nil")
		}
		if err := bucket.Put([]byte(key), []byte(value)); err != nil {
			return fmt.Errorf("put: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}

func (k *KV) Delete(key string) error {
	if err := k.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(k.bucketName)
		if bucket == nil {
			return fmt.Errorf("bucket was nil")
		}
		if err := bucket.Delete([]byte(key)); err != nil {
			return fmt.Errorf("delete: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("update: %w", err)
	}
	return nil
}
