package main

import "fmt"

func main() {
	for val := 0; val < 5; val++ {
		go func() {
			fmt.Println(val)
		}()
	}

	for val := 0; val < 5; val++ {
		fmt.Println(val)
	}
}
