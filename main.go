package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args[1:]) < 2 {
		fmt.Fprintf(os.Stderr, "usage: `kv [set] [key] [value]` or `kv [get] [key]`\n")
		os.Exit(1)
	}
	key := os.Args[2]
	switch os.Args[1] {
	case "get":
		fmt.Printf("get %s\n", key)
	case "set":
		if len(os.Args[1:]) != 3 {
			fmt.Fprintf(os.Stderr, "usage: `kv [set] [key] [value]`\n")
			os.Exit(1)
		}
		value := os.Args[3]
		fmt.Printf("set %s: %s\n", key, value)
	default:
		fmt.Fprintf(os.Stderr, "usage: `kv [set] [key] [value]` or `kv [get] [key]`\n")
		os.Exit(1)
	}
}
