package main

import "fmt"

func main() {
	stringA := "A"
	stringB := "B"

	fmt.Printf("Before swap:\t stringA = %s, stringB = %s\n", stringA, stringB)
	swap(&stringA, &stringB)
	fmt.Printf("After swap:\t stringA = %s, stringB = %s\n", stringA, stringB)
}

func swap(a *string, b *string) {
	tmp := a
	*a = *b
	*b = *tmp
}

