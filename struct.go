package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{Name: "Priyam", Age: 20}
	fmt.Println("Name:", p.Name)
	fmt.Println("Age:", p.Age)
}
