package main

import "fmt"

func main() {
	student := make(map[string]int)

	student["Math"] = 90
	student["Science"] = 85

	fmt.Println("Math marks:", student["Math"])
	fmt.Println("Science marks:", student["Science"])
}
