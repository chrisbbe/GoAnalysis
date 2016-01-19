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
package graph

import (
	"testing"
)


func TestDepthFirstSearchInGraph(t *testing.T) {
	//Create some nodes.
	a := Node{Value:"A"}
	b := Node{Value:"B"}
	c := Node{Value:"C"}
	d := Node{Value:"D"}
	e := Node{Value:"E"}
	f := Node{Value:"F"}
	g := Node{Value:"G"}
	h := Node{Value:"H"}

	graph := New()

	//Add directed node-pairs to graph.
	graph.InsertNode(&a, &b)
	graph.InsertNode(&a, &d)
	graph.InsertNode(&b, &d)
	graph.InsertNode(&b, &c)
	graph.InsertNode(&c, &e)
	graph.InsertNode(&e, &g)
	graph.InsertNode(&e, &f)
	graph.InsertNode(&f, &h)

	dfs := graph.GetDFS()
	correctDfs := []string{"A", "B", "D", "C", "E", "G", "F", "H"}

	//First test length of both DFS-lists.
	if len(correctDfs) != len(dfs) {
		t.Errorf("Length of DFS (%d) is not equal length of correct DFS (%d)!\n", len(dfs), len(correctDfs))
	}
	//Compare DFS with correct DFS.
	for i := 0; i < len(dfs); i++ {
		if dfs[i].Value != correctDfs[i] {
			t.Errorf("Index %d in DFS-list is %s, should be %s!\n", i, dfs[i].Value, correctDfs[i])
		}
	}
}
