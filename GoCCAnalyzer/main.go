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
	generateGraph()
}

func generateGraph() {
	g := graph.New()

	file, err := os.Open("./directedgraph.txt")
	if err != nil {
		fmt.Println("Error opening file!")
		os.Exit(1)
	}
	defer file.Close()

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
	srcFile := "./test_src.go"

	sourceFile, err := ioutil.ReadFile(srcFile)
	if err != nil {
		fmt.Printf("Error finding file\n")
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	basicBlocks := bblock.GetBasicBlocksFromSourceCode(fset, file)

	for _, bb := range basicBlocks {
		fmt.Printf("################## BLOCK NR. %d (%d - %d) ##################\n", bb.Number, bb.FromLine, bb.ToLine)
	}
}
