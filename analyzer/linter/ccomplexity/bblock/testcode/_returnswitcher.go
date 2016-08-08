// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() {
	// BB #0 ending.
	month := 10
	fmt.Printf("Month %d is %s\n", month, monthNumberToString(month))
} // BB #1 ending.

func monthNumberToString(month int) string {
	// BB #2 ending
	switch month { // BB #3 ending.
	case 1:
		return "January" // BB #4 ending.
	case 2:
		return "Febrary" // BB #5 ending.
	case 3:
		return "March" // BB #6 ending.
	case 4:
		return "April" // BB #7 ending.
	case 5:
		return "May" // BB #8 ending.
	case 6:
		return "June" // BB #9 ending.
	case 7:
		return "Juni" // BB #10 ending.
	case 8:
		return "August" // BB #11 ending.
	case 9:
		return "September" // BB #12 ending.
	case 10:
		return "October" // BB #13 ending.
	case 11:
		return "November" // BB #14 ending.
	case 12:
		return "Desember" // BB #15 ending.
	}
	return "Unknown month"
} // BB #16 ending.
