package mbta

// INode interface defines a node in a graph
type INode interface {
	INodeGetEdges
	INodeAddEdge
}

// Node represents a node in a graph
type Node struct {
	edges []INode
}

// NewNode returns a initialized INode struct
func NewNode() *Node {
	edges := make([]INode, 0)
	return &Node{
		edges: edges,
	}
}

// INodeGetEdges interface defines a node that can return its edges
type INodeGetEdges interface {
	GetEdges() *[]INode
}

// GetEdges returns the edges of this node
func (n *Node) GetEdges() *[]INode {
	return &n.edges
}

// INodeAddEdge interface defines a node that can add edges
type INodeAddEdge interface {
	AddEdge(INode)
}

// AddEdge makes a bidirectional edge between this node and the parameter node
func (n *Node) AddEdge(node INode) {
	n.edges = append(n.edges, node)
	edges := node.GetEdges()
	*edges = append(*edges, n)
}
