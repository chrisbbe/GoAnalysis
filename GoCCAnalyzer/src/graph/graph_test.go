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
