// Copyright (c) 2015-2016 The GoAnalysis Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can
// be found in the LICENSE file.
package graph

import (
	"bytes"
	"fmt"
	"github.com/chrisbbe/GoAnalysis/analyzer/globalvars"
	"github.com/chrisbbe/GoAnalysis/analyzer/linter/ccomplexity/graph/stack"
	"io"
	"os"
	"os/exec"
	"time"
)

// All node-values in a graph implements the Value interface.
type Value interface {
	UID() string
	String() string
}

// Graph is the main data-structure representing the graph.
// Holding references to the root node and all nodes in the graph.
type Graph struct {
	Root     *Node                         //First node added to the graph.
	Nodes    map[string]*Node              //All nodes in the graph.
	scc      []*StronglyConnectedComponent //Holds strongly-connected components sets.
	preCount int                           // Used by SCC.
	stack    stack.Stack                   // Used by SCC.
}

// Node represents a single node in the graph, holding the node value.
type Node struct {
	Value            //Holds node interface value.
	visited  bool    //Used by DFS.
	inEdges  []*Node //Holds in-going edges.
	outEdges []*Node //Holds out-going edges.
	low      int     // Used by SCC.
}

// StronglyConnectedComponent holds a set of nodes in a strongly connected components.
type StronglyConnectedComponent struct {
	Nodes []*Node
}

// NewGraph initialize and returns a pointer to a new graph object.
func NewGraph() *Graph {
	return &Graph{Nodes: map[string]*Node{}}
}

func (graph *Graph) InsertNode(node *Node) {
	if graph.Root == nil {
		graph.Root = node
		graph.Nodes[graph.Root.Value.UID()] = graph.Root
	} else if graph.Nodes[node.Value.UID()] == nil {
		graph.Nodes[node.Value.UID()] = node
	}
}

// insertOutNode inserts outgoing directed edge to 'node'.
func (n *Node) insertOutEdge(node *Node) {
	n.outEdges = append(n.outEdges, node)
}

// insertInEdge inserts ingoing directed edge to 'node'.
func (n *Node) insertInEdge(node *Node) {
	n.inEdges = append(n.inEdges, node)
}

// Returning back the interface function from NodeValue.
func (n *Node) String() string {
	return n.Value.String()
}

// InsertEdge inserts directed edge between leftNode and rightNode. It also inserts the node in the graph correctly
// if the node does not already exist in the graph.
func (graph *Graph) InsertEdge(leftNode *Node, rightNode *Node) {
	if len(graph.Nodes) == 0 {
		graph.Root = leftNode
		graph.Root.insertOutEdge(rightNode)
		rightNode.insertInEdge(graph.Root)
		graph.Nodes[graph.Root.Value.UID()] = graph.Root
		graph.Nodes[rightNode.Value.UID()] = rightNode
	} else {
		//Get left and right node if they already exist.
		if graph.Nodes[leftNode.Value.UID()] != nil {
			leftNode = graph.Nodes[leftNode.Value.UID()]
		}
		if graph.Nodes[rightNode.Value.UID()] != nil {
			rightNode = graph.Nodes[rightNode.Value.UID()]
		}

		if graph.Nodes[leftNode.Value.UID()] == nil {
			leftNode.insertOutEdge(rightNode)
			rightNode.insertInEdge(leftNode)
			graph.Nodes[leftNode.Value.UID()] = leftNode
			graph.Nodes[rightNode.Value.UID()] = rightNode
		} else {
			leftNode.insertOutEdge(rightNode)
			rightNode.insertInEdge(leftNode)
			graph.Nodes[rightNode.Value.UID()] = rightNode
		}
	}
}

// getDFS is an internal helper method for GetDFS() to perform depth-first-search on the graph.
func (node *Node) getDFS() (nodes []*Node) {
	if !node.visited {
		nodes = append(nodes, node)
	}
	node.visited = true
	for _, n := range node.outEdges {
		if !n.visited {
			nodes = append(nodes, n.getDFS()...)
		}
	}
	return nodes
}

// GetDFS performs depth first search in the 'graph' and returns the result in 'nodes'.
func (graph *Graph) GetDFS() (nodes []*Node) {
	nodes = graph.Root.getDFS()
	//Clean up nodes by setting visited = false
	for _, node := range graph.Nodes {
		node.visited = false
	}
	return nodes
}

// dfs is internal modified depth first search method for GetSCComponents to find strongly connected components.
func (graph *Graph) dfs(v *Node) {
	v.low = graph.preCount
	graph.preCount++
	v.visited = true
	graph.stack.Push(v)

	min := v.low
	for _, w := range v.outEdges {
		if !w.visited {
			graph.dfs(w)
		}
		if w.low < min {
			min = w.low
		}
	}

	if min < v.low {
		v.low = min
		return
	}

	component := []*Node{}
	var w interface{}
	var err error
	for ok := true; ok; ok = w.(*Node) != v {
		w, err = graph.stack.Pop()
		if err != nil {
			panic(err)
		}
		component = append(component, w.(*Node))
		w.(*Node).low = len(graph.Nodes) - 1
	}

	graph.scc = append(graph.scc, &StronglyConnectedComponent{Nodes: component})
}

// GetSCComponents performs Tarjans algorithm to detect
// strongly connected components in 'graph', returns
// a list of lists containing nodes in each strongly
// connected component.
func (graph *Graph) GetSCComponents() []*StronglyConnectedComponent {
	graph.scc = []*StronglyConnectedComponent{} //Init list.
	for _, node := range graph.Nodes {
		if !node.visited {
			graph.dfs(node)
		}
	}
	//Clean up nodes by setting visited = false
	for _, node := range graph.Nodes {
		node.visited = false
		node.low = 0
	}

	//Clean up graph.
	graph.stack = stack.Stack{}
	graph.preCount = 0
	return graph.scc
}

// GetNumberOfNodes returns number of
// nodes in 'graph'.
func (graph *Graph) GetNumberOfNodes() int {
	return len(graph.Nodes)
}

// GetNumberOfEdges returns number of
// edges in 'graph'.
func (graph *Graph) GetNumberOfEdges() (numberOfEdges int) {
	for _, node := range graph.Nodes {
		numberOfEdges += node.GetOutDegree()
	}
	return numberOfEdges
}

// GetNumberOfSCComponents return number of
// strongly connected components in 'graph'.
func (graph *Graph) GetNumberOfSCComponents() int {
	//We don't want to run Tarjan's once again if algorithm is already executed.
	if graph.scc == nil {
		return len(graph.GetSCComponents())
	}
	return len(graph.scc)
}

// GetInDegree returns number of ingoing edges to 'node'.
func (node *Node) GetInDegree() int {
	return len(node.inEdges)
}

// GetOutDegree returns number of outgoing edges from 'node'.
func (node *Node) GetOutDegree() int {
	return len(node.outEdges)
}

func (node *Node) GetOutNodes() []*Node {
	return node.outEdges
}

func (node *Node) GetInNodes() []*Node {
	return node.inEdges
}

// Draw writes the graph to file according to the Graphviz (www.graphviz.org) format
// and tries to compile the graph to PDF by using dot.
func (graph *Graph) Draw(name string) (err error) {
	dottyFile, err := os.Create(name + ".dot")
	if err != nil {
		return err
	}
	defer func() {
		err = dottyFile.Close()
	}()

	var content bytes.Buffer

	// Write header information.
	content.WriteString("/* --------------------------------------------------- */\n")
	content.WriteString(fmt.Sprintf("/* Generated by %s\n", globalvars.PROGRAM_NAME))
	content.WriteString(fmt.Sprintf("/* Version: %s\n", globalvars.VERSION))
	content.WriteString(fmt.Sprintf("/* Website: %s\n", globalvars.WEBSITE))
	content.WriteString(fmt.Sprintf("/* Date: %s\n", time.Now().String()))
	content.WriteString("/* --------------------------------------------------- */\n")

	// Start writing the graph.
	if _, err := content.WriteString("digraph AST {\n"); err != nil {
		return err
	}

	for _, node := range graph.Nodes {
		for _, outNode := range node.GetOutNodes() {
			if _, err = content.WriteString(fmt.Sprintf("\t\"%s\" -> \"%s\";\n", node, outNode)); err != nil {
				return err
			}
		}
	}
	if _, err := content.WriteString("}\n"); err != nil {
		return err
	}

	if _, err := io.WriteString(dottyFile, content.String()); err != nil {
		return err
	}
	cmd := exec.Command("dot", "-Tpdf", dottyFile.Name(), "-o", name+".pdf")
	return cmd.Run()
}
