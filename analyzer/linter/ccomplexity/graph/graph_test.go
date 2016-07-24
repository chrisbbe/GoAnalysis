// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package graph_test

import (
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/graph"
	"testing"
)

// Letter is the type to store in the graph.
type Letter struct {
	letter string
}

// Satisfies the Value interface.
func (l Letter) UID() string {
	return l.letter
}

// Letter must satisfies the Value interface in graph.
func (l Letter) String() string {
	return l.letter
}

func sccExists(correctSCCList, actualSCCList []*graph.StronglyConnectedComponent) bool {
	existCounter := 0

	for _, correctScc := range correctSCCList {
		for _, actualScc := range actualSCCList {
			if existInList(correctScc.Nodes, actualScc.Nodes) {
				existCounter++
			}
		}
	}

	return existCounter == len(correctSCCList)
}

// Compare two lists.
func existInList(elementList, list []*graph.Node) bool {
	existCounter := 0

	for _, e := range elementList {
		for _, r := range list {
			if e.Value == r.Value {
				existCounter++
			}
		}
	}
	return existCounter == len(elementList)
}

func TestDirectedGraph(test *testing.T) {
	//Create some nodes.
	a := graph.Node{Value: Letter{"A"}}
	b := graph.Node{Value: Letter{"B"}}
	c := graph.Node{Value: Letter{"C"}}
	d := graph.Node{Value: Letter{"D"}}
	e := graph.Node{Value: Letter{"E"}}
	f := graph.Node{Value: Letter{"F"}}
	g := graph.Node{Value: Letter{"G"}}
	h := graph.Node{Value: Letter{"H"}}

	graph := graph.NewGraph()

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
	a := graph.Node{Value: Letter{"A"}}
	b := graph.Node{Value: Letter{"B"}}
	c := graph.Node{Value: Letter{"C"}}
	d := graph.Node{Value: Letter{"D"}}
	e := graph.Node{Value: Letter{"E"}}
	f := graph.Node{Value: Letter{"F"}}
	g := graph.Node{Value: Letter{"G"}}
	h := graph.Node{Value: Letter{"H"}}

	graph := graph.NewGraph()

	//Add directed node-pairs to graph.
	graph.InsertEdge(&a, &b)
	graph.InsertEdge(&a, &d)
	graph.InsertEdge(&b, &d)
	graph.InsertEdge(&b, &c)
	graph.InsertEdge(&c, &e)
	graph.InsertEdge(&e, &g)
	graph.InsertEdge(&e, &f)
	graph.InsertEdge(&f, &h)

	expectedDepthFirstSearch := graph.GetDFS()
	correctDfs := []Letter{{"A"}, {"B"}, {"D"}, {"C"}, {"E"}, {"G"}, {"F"}, {"H"}}

	//Equal length.
	if len(correctDfs) != len(expectedDepthFirstSearch) {
		t.Errorf("Length of DFS (%d) is not equal length of correct DFS (%d)!\n", len(expectedDepthFirstSearch), len(correctDfs))
	}

	//Check nodes in depth first search.
	for index, node := range expectedDepthFirstSearch {
		if node.String() != correctDfs[index].String() {
			t.Error("Faen i helvete")
		}
	}

}

func TestDepthFirstSearchInCycleGraph(t *testing.T) {
	//Create some nodes.
	a := graph.Node{Value: Letter{"A"}}
	b := graph.Node{Value: Letter{"B"}}
	c := graph.Node{Value: Letter{"C"}}
	d := graph.Node{Value: Letter{"D"}}
	e := graph.Node{Value: Letter{"E"}}
	f := graph.Node{Value: Letter{"F"}}
	g := graph.Node{Value: Letter{"G"}}

	directGraph := graph.NewGraph()

	directGraph.InsertEdge(&a, &b)
	directGraph.InsertEdge(&a, &c)
	directGraph.InsertEdge(&a, &e)
	directGraph.InsertEdge(&b, &d)
	directGraph.InsertEdge(&b, &f)
	directGraph.InsertEdge(&c, &g)
	directGraph.InsertEdge(&f, &e)

	expectedDepthFirstSearch := directGraph.GetDFS()
	correctDepthFirstSearch := []Letter{
		{"A"}, {"B"}, {"D"}, {"F"}, {"E"}, {"C"}, {"G"},
	}

	//Compare DFS with correct DFS.
	for index := 0; index < len(expectedDepthFirstSearch); index++ {
		if expectedDepthFirstSearch[index].Value.String() != correctDepthFirstSearch[index].letter {
			t.Errorf("Element nr. %d in DFS should be %s, not %s\n", index, correctDepthFirstSearch[index].letter,
				expectedDepthFirstSearch[index].Value.String())
		}
	}

}

func TestStronglyConnectedComponentsInGraph(t *testing.T) {
	directedGraph := graph.NewGraph()

	a := graph.Node{Value: Letter{"A"}}
	b := graph.Node{Value: Letter{"B"}}
	c := graph.Node{Value: Letter{"C"}}
	d := graph.Node{Value: Letter{"D"}}
	e := graph.Node{Value: Letter{"E"}}
	f := graph.Node{Value: Letter{"F"}}
	g := graph.Node{Value: Letter{"G"}}
	h := graph.Node{Value: Letter{"H"}}

	directedGraph.InsertEdge(&a, &b)
	directedGraph.InsertEdge(&b, &c)
	directedGraph.InsertEdge(&c, &d)
	directedGraph.InsertEdge(&d, &c)
	directedGraph.InsertEdge(&d, &h)
	directedGraph.InsertEdge(&h, &d)
	directedGraph.InsertEdge(&c, &g)
	directedGraph.InsertEdge(&h, &g)
	directedGraph.InsertEdge(&f, &g)
	directedGraph.InsertEdge(&g, &f)
	directedGraph.InsertEdge(&b, &f)
	directedGraph.InsertEdge(&e, &f)
	directedGraph.InsertEdge(&e, &a)
	directedGraph.InsertEdge(&b, &e)

	expectedStronglyConnectedComponents := directedGraph.GetSCComponents()
	correctStronglyConnectedComponents := []*graph.StronglyConnectedComponent{
		{
			Nodes: []*graph.Node{
				{Value: Letter{"F"}},
				{Value: Letter{"G"}},
			}},
		{
			Nodes: []*graph.Node{
				{Value: Letter{"H"}},
				{Value: Letter{"C"}},
				{Value: Letter{"D"}},
			}},
		{
			Nodes: []*graph.Node{
				{Value: Letter{"B"}},
				{Value: Letter{"A"}},
				{Value: Letter{"E"}},
			}},
	}

	// Check the number of SCC sets.
	if len(expectedStronglyConnectedComponents) != len(correctStronglyConnectedComponents) {
		t.Fatalf("Number of strongly connected components be %d, but are %d!\n", len(correctStronglyConnectedComponents),
			len(expectedStronglyConnectedComponents))
	}

	if !sccExists(correctStronglyConnectedComponents, expectedStronglyConnectedComponents) {
		t.Error("Not all SCC exists")
	}

}

func TestStronglyConnectedComponentsInGraph2(t *testing.T) {
	directedGraph := graph.NewGraph()

	a := graph.Node{Value: Letter{"A"}}
	b := graph.Node{Value: Letter{"B"}}
	c := graph.Node{Value: Letter{"C"}}
	d := graph.Node{Value: Letter{"D"}}
	e := graph.Node{Value: Letter{"E"}}
	f := graph.Node{Value: Letter{"F"}}
	g := graph.Node{Value: Letter{"G"}}
	h := graph.Node{Value: Letter{"H"}}

	directedGraph.InsertEdge(&a, &b)
	directedGraph.InsertEdge(&a, &f)
	directedGraph.InsertEdge(&b, &f)
	directedGraph.InsertEdge(&b, &c)
	directedGraph.InsertEdge(&c, &d)
	directedGraph.InsertEdge(&c, &g)
	directedGraph.InsertEdge(&e, &a)
	directedGraph.InsertEdge(&f, &e)
	directedGraph.InsertEdge(&f, &g)
	directedGraph.InsertEdge(&g, &c)
	directedGraph.InsertEdge(&h, &g)

	expectedStronglyConnectedComponents := directedGraph.GetSCComponents()
	correctStronglyConnectedComponents := []*graph.StronglyConnectedComponent{
		{
			Nodes: []*graph.Node{
				{Value: Letter{"D"}},
			}},
		{
			Nodes: []*graph.Node{
				{Value: Letter{"C"}},
				{Value: Letter{"G"}},
			}},
		{
			Nodes: []*graph.Node{
				{Value: Letter{"A"}},
				{Value: Letter{"B"}},
				{Value: Letter{"E"}},
				{Value: Letter{"F"}},
			}},
		{
			Nodes: []*graph.Node{
				{Value: Letter{"H"}},
			}},
	}

	// Check the number of SCC sets.
	if len(expectedStronglyConnectedComponents) != len(correctStronglyConnectedComponents) {
		t.Fatalf("Number of strongly connected components be %d, but are %d!\n", len(correctStronglyConnectedComponents),
			len(expectedStronglyConnectedComponents))
	}

	if !sccExists(correctStronglyConnectedComponents, expectedStronglyConnectedComponents) {
		t.Error("Not all SCC exists")
	}
}
