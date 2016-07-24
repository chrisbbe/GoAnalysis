// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package main // BB #0 starting.

import "fmt"

func main() {
	// BB #1 starting.
	fmt.Println("Hello World")
}

// Cyclomatic Complexity M = E - N + 2P.
// E = Number of edges in control flow graph.
// N = Number of nodes in control flow graph.
// P = Number of connected components in graph.
// 		 File level: 	M = 1 - 2 + 2 * 2 = 3
// Function level: 	M = 0 - 1 + 2 * 1 = 1
