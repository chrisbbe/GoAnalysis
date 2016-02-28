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

// CAUTION: This file in used by a unit test in file basicblock_test.go,
// changing the structure or logic in the program will with high possibility
// break the test and give false positive errors. Please DO NOT change this
// file unless you know what you are doing!
package main //BB #0 starting.

import "fmt"

func main() { // BB #1 starting.
	if true { // BB #2 starting.
		fmt.Printf("Sant")
		fmt.Printf("True")
	} else { // BB #3 starting.
		fmt.Printf("Usant")
		fmt.Printf("False")
	}

	fmt.Println("In the middle")

	if true { // BB #4 starting.
		fmt.Printf("Sant")
		fmt.Printf("True")
	} else { // BB #5 starting.
		fmt.Printf("Usant")
		fmt.Printf("False")
	}
}
