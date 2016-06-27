package main

import "fmt"

func main() {
	number := 3
	fmt.Printf("Number %d is %s\n", number, integerToString(month))
}

func integerToString(number int) string {
	switch number {
	case 1:
		return "1"
	case 2:
		return "2"
	case 3:
		return "3"
	default:
		return "Invalid number"
	}
}
