//This little piece of code is used to demonstrate
//one of the common mistakes done i Go, namely to
//use goroutines on loop-iterator variables.
//
//Author: Christian Bergum Bergersen (chrisbbe@ifi.uio.no)
package main

import (
	"fmt"
	"log"
)

// This small piece of code illustrates
// one of the common mistakes new Go
// people do when using Goroutines inside
// for-loops, referencing the loop-variable.
//
func main() {
	//Not-thread safe.
	for val := 0; val < 5; val++ {
		go func() {
			log.Println(val)
		}()
	}

	//Safe for-loop, no Go routine!
	for val := 0; val < 5; val++ {
		log.Println(val)
	}

	//Thread safe for-loop.
	for val := 0; val < 5; val++ {
		go func(val int) {
			log.Println(val)
		}(val)
	}
}
