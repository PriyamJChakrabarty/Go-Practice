package main

import "fmt"

func main() {
	var n int
	fmt.Print("Enter a number: ")
	fmt.Scan(&n)

	isPrime := true

	if n <= 1 {
		isPrime = false
	}

	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			isPrime = false
			break
		}
	}

	if isPrime {
		fmt.Println("Prime number")
	} else {
		fmt.Println("Not prime")
	}
}
