package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	num := rand.Intn(100)
	fmt.Println("Random number:", num)
}
