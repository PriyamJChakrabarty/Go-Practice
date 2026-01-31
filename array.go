package main

import "fmt"

func main() {
	var arr [3]int = [3]int{10, 20, 30}

	for i := 0; i < len(arr); i++ {
		fmt.Println(arr[i])
	}
}
