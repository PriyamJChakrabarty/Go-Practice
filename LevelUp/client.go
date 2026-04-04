package main

import (
	"fmt"
	"net"
)

func main() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	conn.Write([]byte("Hello Server"))
	buf := make([]byte, 1024)
	n, _ := conn.Read(buf)
	fmt.Println(string(buf[:n]))
	conn.Close()
}