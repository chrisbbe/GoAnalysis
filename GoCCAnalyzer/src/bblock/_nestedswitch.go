// The MIT License (MIT)

// Copyright (c) 2015-2016 Christian Bergum Bergersen

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package main // BB #0 starting.

import "fmt"

func main() { // BB #1 starting.
	number := 1
	secondNumber := 3

	switch number { // BB #2 starting.

	case 0: // BB #3 starting.
		fmt.Println("0")
	case 1: // BB #4 starting.
		fmt.Println("2")
		switch secondNumber { // BB #5 starting.

		case 1: // BB #6 starting.
			fmt.Println("1")
		case 2: // BB #7 starting.
			fmt.Println("2")
		default: // BB #8 starting.
			fmt.Printf("No match, secondNumber is %d!\n", number)
		}
	case 3: // BB #9 starting.
		fmt.Println("3")
	case 4: // BB #10 starting.
		fmt.Println("4")
	default: // BB #11 starting.
		fmt.Printf("No match, number is %d!\n", number)
	}
}
