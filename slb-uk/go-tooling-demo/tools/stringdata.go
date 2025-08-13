//go:build ignore
// +build ignore

package main

import (
    "os"
)

func main() {
    _ = os.WriteFile("generated_data.txt", []byte("hello from generator\n"), 0644)
}
