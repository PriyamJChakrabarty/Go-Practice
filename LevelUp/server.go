package main

import (
	"fmt"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":8080")
	for {
		conn, _ := ln.Accept()
		go func(c net.Conn) {
			buf := make([]byte, 1024)
			n, _ := c.Read(buf)
			fmt.Println(string(buf[:n]))
			c.Write([]byte("Message received"))
			c.Close()
		}(conn)
	}
}