package main

import "fmt"

func main() { // BB #0 ending.

	if true { // BB #1 ending.

		fmt.Printf("True")
		if false { // BB #2 ending.
			fmt.Println("False")
		} else { // BB #3 ending.
			i := 10 // BB #4 ending.
		}
	} else { // BB #5 ending.
		y := 5 // BB #6 ending.
	}

	x := 2 + 5 // BB #7 ending.
}
