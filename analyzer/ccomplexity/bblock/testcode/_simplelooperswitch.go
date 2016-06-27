package main

import "fmt"

func main() {

	for i := 0; ; i++ {
		switch i {
		case 0:
			fmt.Printf("Zero")
		case 1:
			fmt.Printf("One")
		}
	}
}
