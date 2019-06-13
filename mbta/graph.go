package mbta

// INodeAddEdge has the add edge method of INode
type INodeAddEdge interface {
	AddEdge(*INode) bool
}

// INodeFindEdge returns the found node or nil is if there is no match
type INodeFindEdge interface {
	FindEdge() *INode
}

// INode is an interface to represent a graph node
type INode interface {
	INodeAddEdge
	INodeFindEdge
}

// Node represents a node in the graph
type Node struct {
	edges []*INode
}

// AddEdge adds a connected node to this node
func (n *Node) AddEdge(edge *INode) {
	n.edges = append(n.edges, edge)
}

// FindEdge removes a connected edge from this one
func (n *Node) FindEdge() *INode {
	n.edge
}

type Graph struct {
	
}