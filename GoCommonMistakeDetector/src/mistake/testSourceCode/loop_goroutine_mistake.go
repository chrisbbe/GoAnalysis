//This little piece of code is used to demonstrate
//one of the common mistakes done i Go, namely to
//use goroutines on loop-iterator variables.
//
//Author: Christian Bergum Bergersen (chrisbbe@ifi.uio.no)
package main

import "fmt"

// This small piece of code illustrates
// one of the common mistakes new Go
// people do when using Goroutines inside
// for-loops, referencing the loop-variable.
//
func main() {
	//Not-thread safe.
	for val := 0; val < 5; val++ {
		go func() {
			fmt.Println(val)
		}()
	}

	//Safe foor-loop, no goroutine!
	for val := 0; val < 5; val++ {
		fmt.Println(val)
	}

	//Thread safe foor-loop.
	for val := 0; val < 5; val++ {
		go func(val int) {
			fmt.Println(val)
		}(val)
	}
}
