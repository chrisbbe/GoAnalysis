package main

import (
	"./graph"
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	g := graph.New()

	file, err := os.Open("GoCCAnalyzer/src/directedgraph.txt")
	if err != nil {
		fmt.Println("Error opening file!")
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), " -> ")
		left := graph.Node{Id: line[0]}
		right := graph.Node{Id: line[1]}
		g.InsertNode(&left, &right)
	}

	fmt.Println("### DFS Printout ###")
	for _, node := range g.GetDFS() {
		fmt.Printf("%s\n", node.Id)
	}

}
