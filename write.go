package main

import (
	"os"
)

func main() {
	content := []byte("Hello from Go!")

	os.WriteFile("output.txt", content, 0644)
}
