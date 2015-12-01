package main

import (
	"./graph"
	"bufio"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"go/token"
	"go/parser"
	"./bblock"
)

func main() {
	getBasicBlocks()
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
