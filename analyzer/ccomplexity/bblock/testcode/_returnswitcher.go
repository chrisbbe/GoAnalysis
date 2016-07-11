// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main

import "fmt"

func main() { // 0.
	month := 10
	fmt.Printf("Month %d is %s\n", month, monthNumberToString(month))
} // 1.

func monthNumberToString(month int) string { // 2
	switch month { // 3
	case 1:
		return "January"
	case 2:
		return "Febrary"
	case 3:
		return "March"
	case 4:
		return "April"
	case 5:
		return "May"
	case 6:
		return "June"
	case 7:
		return "Juni"
	case 8:
		return "August"
	case 9:
		return "September"
	case 10:
		return "October"
	case 11:
		return "November"
	case 12:
		return "Desember"
	}
}
