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

func TestDirectedGraph(test *testing.T) {
	//Create some nodes.
	a := Node{Value: "A"}
	b := Node{Value: "B"}
	c := Node{Value: "C"}
	d := Node{Value: "D"}
	e := Node{Value: "E"}
	f := Node{Value: "F"}
	g := Node{Value: "G"}
	h := Node{Value: "H"}

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

	//Test number of nodes in graph.
	if graph.GetNumberOfNodes() != 8 {
		test.Fatalf("Graph should contain 8 nodes, not %d!\n", graph.GetNumberOfNodes())
	}

	//Node A should be root node.
	if graph.Root.Value != a.Value {
		test.Errorf("Node A should be root node, not node %v\n", graph.Root.Value)
	}

	//Test node A.
	if a.GetInDegree() != 0 {
		test.Errorf("Node A in-degree should be 0, not %d\n", a.GetInDegree())
	}
	if a.GetOutDegree() != 2 {
		test.Errorf("Node A out-degree should be 2, not %d\n", a.GetInDegree())
	}

	//Test node B.
	if b.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", b.GetInDegree())
	}
	if b.GetOutDegree() != 2 {
		test.Errorf("Node A out-degree should be 2, not %d\n", b.GetInDegree())
	}

	//Test node C.
	if c.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", c.GetInDegree())
	}
	if c.GetOutDegree() != 1 {
		test.Errorf("Node A out-degree should be 1, not %d\n", c.GetInDegree())
	}

	//Test node D.
	if d.GetInDegree() != 2 {
		test.Errorf("Node A in-degree should be 2, not %d\n", d.GetInDegree())
	}
	if d.GetOutDegree() != 0 {
		test.Errorf("Node A out-degree should be 0, not %d\n", d.GetInDegree())
	}

	//Test node E.
	if e.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", e.GetInDegree())
	}
	if e.GetOutDegree() != 2 {
		test.Errorf("Node A out-degree should be 2, not %d\n", e.GetInDegree())
	}

	//Test node F.
	if f.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", f.GetInDegree())
	}
	if f.GetOutDegree() != 1 {
		test.Errorf("Node A out-degree should be 1, not %d\n", f.GetInDegree())
	}

	//Test node G.
	if g.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", g.GetInDegree())
	}
	if g.GetOutDegree() != 0 {
		test.Errorf("Node A out-degree should be 0, not %d\n", g.GetInDegree())
	}

	//Test node H.
	if h.GetInDegree() != 1 {
		test.Errorf("Node A in-degree should be 1, not %d\n", h.GetInDegree())
	}
	if h.GetOutDegree() != 0 {
		test.Errorf("Node A out-degree should be 0, not %d\n", h.GetInDegree())
	}

}

func TestDepthFirstSearchInGraph(t *testing.T) {
	//Create some nodes.
	a := Node{Value: "A"}
	b := Node{Value: "B"}
	c := Node{Value: "C"}
	d := Node{Value: "D"}
	e := Node{Value: "E"}
	f := Node{Value: "F"}
	g := Node{Value: "G"}
	h := Node{Value: "H"}

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

func TestDepthFirstSearchInCycleGraph(test *testing.T) {
	//Create some nodes.
	a := Node{Value: "A"}
	b := Node{Value: "B"}
	c := Node{Value: "C"}
	d := Node{Value: "D"}
	e := Node{Value: "E"}
	f := Node{Value: "F"}
	g := Node{Value: "G"}

	graph := NewGraph()

	graph.InsertEdge(&a, &b)
	graph.InsertEdge(&a, &c)
	graph.InsertEdge(&a, &e)
	graph.InsertEdge(&b, &d)
	graph.InsertEdge(&b, &f)
	graph.InsertEdge(&c, &g)
	graph.InsertEdge(&f, &e)

	expectedDepthFirstSearch := graph.GetDFS()
	actualDepthFirstSearch := []string{"A", "B", "D", "F", "E", "C", "G"}
	//Compare DFS with correct DFS.
	for i := 0; i < len(expectedDepthFirstSearch); i++ {
		if expectedDepthFirstSearch[i].Value != actualDepthFirstSearch[i] {
			test.Errorf("Index %d in DFS-list is %s, should be %s!\n", i, expectedDepthFirstSearch[i].Value,
				actualDepthFirstSearch[i])
		}
	}
}

func TestStronglyConnectedComponentsInGraph(t *testing.T) {
	graph := NewGraph()

	a := &Node{Value: 0}
	b := &Node{Value: 1}
	c := &Node{Value: 2}
	d := &Node{Value: 3}
	e := &Node{Value: 4}
	f := &Node{Value: 5}
	g := &Node{Value: 6}
	h := &Node{Value: 7}

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
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: 5}, &Node{Value: 6}}},
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: 7}, &Node{Value: 2}, &Node{Value: 3}}},
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: 1}, &Node{Value: 0}, &Node{Value: 4}}},
	}

	if len(expectedStronglyConnectedComponents) != len(actualStronglyConnectedComponents) {
		t.Fatalf("Number of strongly connected components be %d, but are %d!\n", len(actualStronglyConnectedComponents),
			len(expectedStronglyConnectedComponents))
	}

	for i, e := range actualStronglyConnectedComponents {
		for _, f := range e.Nodes {
			//Hver eneste variabel i i
			if !containsNodeValue(f.Value, expectedStronglyConnectedComponents[i].Nodes) {
				t.Fatalf("Strongly connected component nr. %d should contain value %v!\n", i, f.Value)
			}
		}
	}
}

func TestStronglyConnectedComponentsInGraph2(t *testing.T) {
	graph := NewGraph()

	a := &Node{Value: "A"}
	b := &Node{Value: "B"}
	c := &Node{Value: "C"}
	d := &Node{Value: "D"}
	e := &Node{Value: "E"}
	f := &Node{Value: "F"}
	g := &Node{Value: "G"}
	h := &Node{Value: "H"}

	graph.InsertEdge(a, b)
	graph.InsertEdge(a, f)

	graph.InsertEdge(b, f)
	graph.InsertEdge(b, c)

	graph.InsertEdge(c, d)
	graph.InsertEdge(c, g)

	graph.InsertEdge(e, a)

	graph.InsertEdge(f, e)
	graph.InsertEdge(f, g)

	graph.InsertEdge(g, c)

	graph.InsertEdge(h, g)

	expectedStronglyConnectedComponents := graph.GetSCComponents()
	actualStronglyConnectedComponents := []StronglyConnectedComponent{
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: "D"}}},
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: "C"}, &Node{Value: "G"}}},
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: "A"}, &Node{Value: "B"}, &Node{Value: "E"}, &Node{Value: "F"}}},
		StronglyConnectedComponent{Nodes: []*Node{&Node{Value: "H"}}},
	}

	if len(expectedStronglyConnectedComponents) != len(actualStronglyConnectedComponents) {
		t.Fatalf("Number of strongly connected components should be %d, but are %d!\n", len(actualStronglyConnectedComponents),
			len(expectedStronglyConnectedComponents))
	}

	for i, e := range actualStronglyConnectedComponents {
		for _, f := range e.Nodes {
			//TODO
			//Hver eneste variabel i i
			if !containsNodeValue(f.Value, expectedStronglyConnectedComponents[i].Nodes) {
				t.Fatalf("Strongly connected component nr. %d should contain value %v!\n", i, f.Value)
			}
		}
	}
}

//Helper function to test if a list of Nodes contains a value.
func containsNodeValue(value interface{}, nodes []*Node) bool {
	for _, v := range nodes {
		if v.Value == value {
			return true
		}
	}
	return false
}
