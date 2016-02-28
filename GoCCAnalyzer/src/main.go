package main

import (
	"fmt"
	"io/ioutil"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/ccomplexity"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/graph"
	"go/token"
	"go/parser"
	"github.com/chrisbbe/GoAnalysis/GoCCAnalyzer/src/bblock"
	"os"
	"io"
)

func main() {
	//static()
	//testGraph()
	testing()
	//testCC()
	testBB()
}

func testGraph() {
	g := graph.NewGraph()

	a := graph.Node{Value:"A"}
	b := graph.Node{Value:"B"}
	c := graph.Node{Value:"C"}

	g.InsertNode(&a)
	g.InsertNode(&b)
	g.InsertNode(&c)

	fmt.Printf("Length %d\n", len(g.Nodes))

	for i, e := range g.GetDFS() {
		fmt.Printf("%d) %v\n", i, e.Value)
	}
}

func testCC() {
	sourceFile, err := ioutil.ReadFile("./ccomplexity/testcode/_ifelse.go")
	if err != nil {
		fmt.Printf("Error:\n")
	}

	fmt.Printf("CC File Level: %d\n", ccomplexity.GetCyclomaticComplexityFileLevel(sourceFile))

	for _, bbCC := range ccomplexity.GetCyclomaticComplexityFunctionLevel(sourceFile) {
		fmt.Printf("- %s: CC = %d\n", bbCC.FunctionName, bbCC.GetCyclomaticComplexity())
	}
}

func testBB() {
	sourceFile, err := ioutil.ReadFile("./ccomplexity/testcode/_twoifelse.go")
	if err != nil {
		fmt.Printf("Error:\n")
	}
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	for _, bb := range bblock.GetBasicBlocksFromSourceCode(file) {
		fmt.Printf("%s (%d)\n", bb.Type.String(), bb.Number)
		for _, sbb := range bb.Successor {
			fmt.Printf("\t- %s (%d)\n", sbb.Type.String(), sbb.Number)
		}
	}

}

func testing() {
	sourceFile, err := ioutil.ReadFile("./ccomplexity/testcode/_ifelse.go")
	if err != nil {
		fmt.Printf("Error:\n")
	}

	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, "", sourceFile, 0)
	if err != nil {
		panic(err)
	}

	controlFlow := graph.NewGraph()

	for _, bb := range bblock.GetBasicBlocksFromSourceCode(file) {
		controlFlow.InsertNode(&graph.Node{Value:bb})
		for _, sbb := range bb.Successor {
			controlFlow.InsertEdge(&graph.Node{Value:bb}, &graph.Node{Value:sbb})
		}
	}

	dottyFile, err := os.Create("result.gv")
	writeLineToFile("digraph AST {\n", dottyFile)

	for _, node := range controlFlow.Nodes {
		fmt.Printf("%s\n", node.Value.(*bblock.BasicBlock).Type.String())

		for _, sNode := range node.GetOutNodes() {
			fmt.Printf("\t-> %s\n", sNode.Value.(*bblock.BasicBlock).Type.String())
			line := fmt.Sprintf("\t\"(BB #%d) %s\" -> \"(BB #%d) %s\";\n", node.Value.(*bblock.BasicBlock).Number, node.Value.(*bblock.BasicBlock).Type.String(), sNode.Value.(*bblock.BasicBlock).Number, sNode.Value.(*bblock.BasicBlock).Type.String())
			writeLineToFile(line, dottyFile)
		}
		fmt.Println()
	}

	writeLineToFile("}\n", dottyFile)
	dottyFile.Close()

}

func static() {
	start := &graph.Node{Value:"START"}
	exit := &graph.Node{Value:"EXIT"}
	ifNode := &graph.Node{Value:"IF"}
	elseNode := &graph.Node{Value:"ELSE"}

	controlFlow := graph.NewGraph()

	controlFlow.InsertEdge(start, ifNode)
	controlFlow.InsertEdge(start, elseNode)
	controlFlow.InsertEdge(ifNode, exit)
	controlFlow.InsertEdge(elseNode, exit)

	dottyFile, _ := os.Create("result.gv")
	writeLineToFile("digraph AST {\n", dottyFile)

	for _, node := range controlFlow.Nodes {
		fmt.Printf("%s\n", node.Value)

		for _, sNode := range node.GetOutNodes() {
			fmt.Printf("\t-> %s\n", sNode.Value)
			line := fmt.Sprintf("\t\"%s\" -> \"%s\";\n", node.Value, sNode.Value)
			writeLineToFile(line, dottyFile)
		}
		fmt.Println()
	}

	writeLineToFile("}\n", dottyFile)
	dottyFile.Close()
}

func writeLineToFile(line string, f io.Writer) {
	_, err := io.WriteString(f, line)
	if err != nil {
		fmt.Println(err)
	}
}
