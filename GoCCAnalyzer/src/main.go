package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"go/token"
	"go/parser"
	"github.com/chrisbbe/GoThesis/GoCCAnalyzer/src/graph"
	"github.com/chrisbbe/GoThesis/GoCCAnalyzer/src/bblock"
)

func main() {
	generateControlFlowGraph()
	//getBasicBlocks()
	//generateGraph()
}

func generateGraph() {
	srcFile := "../directedgraph.txt"

	file, err := os.Open(srcFile)
	if err != nil {
		fmt.Printf("Error opening file %s!\n", srcFile)
		os.Exit(1)
	}
	defer file.Close()

	g := graph.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " -> ")
		left := graph.Node{Value: line[0]}
		right := graph.Node{Value: line[1]}
		g.InsertNode(&left, &right)
	}

	fmt.Println("### DFS Printout ###")
	for _, node := range g.GetDFS() {
		fmt.Printf("%s\n", node.Value)
	}
}

func getBasicBlocks() {
	srcFile := "../codeexamples/_ifelse.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error finding file %s!\n", srcFile)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := bblock.GetBasicBlocksFromSourceCode(fset, file)

	for _, bb := range basicBlocks {
		fmt.Printf("################## BLOCK NR. %d (%d - %d) %s ##################\n", bb.Number, bb.FromLine, bb.ToLine, bb.Value)
	}
}

type cfgNode struct {
	Value      string
	bb         *bblock.BasicBlock
	ifTrueJmp  *bblock.BasicBlock
	ifFalseJmp *bblock.BasicBlock
}

func generateControlFlowGraph() {

	srcFile := "../codeexamples/_ifelse.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error finding file %s!\n", srcFile)
		os.Exit(1)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := bblock.GetBasicBlocksFromSourceCode(fset, file)

	fmt.Printf("Number of basicBlocks: %d\n", len(basicBlocks))
	g := graph.New()



	entryCfg := cfgNode{Value:fmt.Sprintf("%d) %s", 0, "Entry")}
	entry := graph.Node{Value:entryCfg}

	var prevNode *graph.Node = nil

	for i, v := range basicBlocks {
		line := fmt.Sprintf("%d) %s", i, v.Value)
		cfgNode := cfgNode{Value:line, bb:v}

		if v.Value == "If" { //Hack
			cfgNode.ifTrueJmp = basicBlocks[i + 1] //Always are If-True-jump currentBlock + 1
			cfgNode.ifFalseJmp = basicBlocks[i + 3]
		}

		node := graph.Node{Value:cfgNode}

		if prevNode == nil { //First node after entry.
			g.InsertNode(&entry, &node)
		} else {
			g.InsertNode(prevNode, &node)
		}
		prevNode = &node
	}

	exit := graph.Node{Value:fmt.Sprintf("%d) %s", len(basicBlocks), "Exit")}
	g.InsertNode(prevNode, &exit)

	dfs := g.GetDFS()
	fmt.Printf("Number of nodes in DFS %d\n", len(dfs))

	for _, key := range dfs {
		//fmt.Printf("%s\n", key)
		fmt.Println(key)
	}
}