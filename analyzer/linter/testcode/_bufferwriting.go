// Copyright (c) 2015-2016 The GoAnalysis Authors. All rights reserved.
// Use of this source code is governed by the MIT license found in the
// LICENSE file.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	w := bufio.NewWriter(os.Stdout)
	fmt.Fprint(w, "Hello, World")
	closeBuffer(w)
}

func closeBuffer(buf *bufio.Writer) {
	buf.Flush()
}
