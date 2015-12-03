package main

import "fmt"

func main() {
	// BB #0
	if true { // BB #1
		fmt.Printf("Sant")
		fmt.Printf("True") // BB #2
	} else { // BB #3
		fmt.Printf("Usant")
		fmt.Printf("False") // BB #4
	}
}
