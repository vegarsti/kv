package main

import (
	"fmt"
	"kv/bolt"
	"os"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Fprintf(os.Stderr, "the kv CLI is not implemented yet\n")
		os.Exit(1)
	}
	if len(os.Args[1:]) == 1 || len(os.Args[1:]) > 3 {
		fmt.Fprintf(os.Stderr, "usage: `kv [put] [key] [value]` or `kv [get] [key]` or `kv [delete] [key]`\n")
		os.Exit(1)
	}

	// k := socket.New("/tmp/kv.sock")
	// k := embedded.New()
	k := bolt.New()
	if err := k.Open(); err != nil {
		fmt.Fprintf(os.Stderr, "open: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := k.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "close: %v\n", err)
			os.Exit(1)
		}
	}()
	key := os.Args[2]
	switch os.Args[1] {
	case "get":
		value, err := k.Get(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "get %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", value)
	case "put":
		if len(os.Args[1:]) != 3 {
			fmt.Fprintf(os.Stderr, "usage: `kv [put] [key] [value]`\n")
			os.Exit(1)
		}
		value := os.Args[3]
		if err := k.Put(key, value); err != nil {
			fmt.Fprintf(os.Stderr, "put %s: %s: %v\n", key, value, err)
			os.Exit(1)
		}
		fmt.Println("OK")
	case "delete":
		if err := k.Delete(key); err != nil {
			fmt.Fprintf(os.Stderr, "delete %s: %v\n", key, err)
			os.Exit(1)
		}
		fmt.Println("OK")
	default:
		fmt.Fprintf(os.Stderr, "usage: `kv [put] [key] [value]` or `kv [get] [key]` or `kv [delete] [key]`\n")
		os.Exit(1)
	}
}
