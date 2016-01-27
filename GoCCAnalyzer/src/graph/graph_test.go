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

	graph := NewGraph()

	//Add directed node-pairs to graph.
	graph.InsertEdge(&a, &b)
	graph.InsertEdge(&a, &d)
	graph.InsertEdge(&b, &d)
	graph.InsertEdge(&b, &c)
	graph.InsertEdge(&c, &e)
	graph.InsertEdge(&e, &g)
	graph.InsertEdge(&e, &f)
	graph.InsertEdge(&f, &h)

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

func TestStronglyConnectedComponentsInGraph(t *testing.T) {
	graph := NewGraph()

	a := &Node{Value:0}
	b := &Node{Value:1}
	c := &Node{Value:2}
	d := &Node{Value:3}
	e := &Node{Value:4}
	f := &Node{Value:5}
	g := &Node{Value:6}
	h := &Node{Value:7}

	graph.InsertEdge(a, b)
	graph.InsertEdge(b, c)
	graph.InsertEdge(c, d)
	graph.InsertEdge(d, c)
	graph.InsertEdge(d, h)
	graph.InsertEdge(h, d)
	graph.InsertEdge(c, g)
	graph.InsertEdge(h, g)
	graph.InsertEdge(f, g)
	graph.InsertEdge(g, f)
	graph.InsertEdge(b, f)
	graph.InsertEdge(e, f)
	graph.InsertEdge(e, a)
	graph.InsertEdge(b, e)

	expectedStronglyConnectedComponents := graph.GetSCComponents()
	actualStronglyConnectedComponents := []StronglyConnectedComponent{
		StronglyConnectedComponent{Nodes:[]*Node{&Node{Value:5}, &Node{Value:6}}},
		StronglyConnectedComponent{Nodes:[]*Node{&Node{Value:7}, &Node{Value:2}, &Node{Value:3}}},
		StronglyConnectedComponent{Nodes:[]*Node{&Node{Value:1}, &Node{Value:0}, &Node{Value:4}}},
	}

	if len(expectedStronglyConnectedComponents) != len(actualStronglyConnectedComponents) {
		t.Fatalf("Number of strongly connected components be %d, but are %d!\n", len(actualStronglyConnectedComponents),
			len(expectedStronglyConnectedComponents))
	}

	for i, e := range actualStronglyConnectedComponents {
		for _, f := range e.Nodes { //Hver eneste variabel i i
			if !containsNodeValue(f.Value, expectedStronglyConnectedComponents[i].Nodes) {
				t.Fatalf("Strongly connected component nr. %d should contain value %v!\n", i, f.Value)
			}
		}
	}
}

//Helper function to test if a list of Nodes contains a value.
func containsNodeValue(value interface{}, nodes []*Node) (bool) {
	for _, v := range nodes {
		if v.Value == value {
			return true
		}
	}
	return false
}
