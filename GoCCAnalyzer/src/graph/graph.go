package graph

type Graph struct {
	Root  *Node
	Nodes map[string]*Node
}

type Node struct {
	Id       string
	visited  bool
	inEdges  []*Node
	outEdges []*Node
}

func New() *Graph {
	return &Graph{Nodes: map[string]*Node{}}
}

func (n *Node) insertOutNode(node *Node) {
	n.outEdges = append(n.outEdges, node)
}

func (n *Node) insertInNode(node *Node) {
	n.inEdges = append(n.inEdges, node)
}

func (graph *Graph) InsertNode(leftNode *Node, rightNode *Node) {

	if len(graph.Nodes) == 0 {
		graph.Root = leftNode
		graph.Root.insertOutNode(rightNode)
		rightNode.insertInNode(graph.Root)
		graph.Nodes[graph.Root.Id] = graph.Root
		graph.Nodes[rightNode.Id] = rightNode
	} else {
		//Get right left and right node if they already exist.
		if graph.Nodes[leftNode.Id] != nil {
			leftNode = graph.Nodes[leftNode.Id]
		}
		if graph.Nodes[rightNode.Id] != nil {
			rightNode = graph.Nodes[rightNode.Id]
		}

		if graph.Nodes[leftNode.Id] == nil {
			leftNode.insertOutNode(rightNode)
			rightNode.insertInNode(leftNode)
			graph.Nodes[leftNode.Id] = leftNode
			graph.Nodes[rightNode.Id] = rightNode
		} else {
			leftNode.insertOutNode(rightNode)
			rightNode.insertInNode(leftNode)
			graph.Nodes[rightNode.Id] = rightNode
		}
	}
}

func (node *Node) getDFS() (nodes []*Node) {
	if !node.visited {
		nodes = append(nodes, node)
	}
	for _, n := range node.outEdges {
		nodes = append(nodes, n.getDFS()...)
	}
	node.visited = true
	return nodes
}

func (graph *Graph) GetDFS() (nodes []*Node) {
	nodes = graph.Root.getDFS()
	//Clean up nodes by setting visited = false
	for _, node := range graph.Nodes {
		node.visited = false
	}
	return nodes
}
