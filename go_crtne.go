package main

import (
	"fmt"
	"time"
)

func worker(id int, ch chan string) {
	time.Sleep(time.Millisecond * 500)
	ch <- fmt.Sprintf("done %d", id)
}

func main() {
	ch := make(chan string)

	for i := 1; i <= 3; i++ {
		go worker(i, ch)
	}

	for i := 1; i <= 3; i++ {
		fmt.Println(<-ch)
	}
}