package mbta

import (
	"testing"
)

func Test_AddEdge(t *testing.T) {
	node1 := NewNode()
	node2 := NewNode()

	node1.AddEdge(node2)

	if node1.edges[0] != node2 {
		t.FailNow()
	}

	if node2.edges[0] != node1 {
		t.FailNow()
	}
}

func Test_GetEdges(t *testing.T) {
	node1 := NewNode()
	node2 := &Node{
		edges: []INode {node1},
	}

	edges := *node2.GetEdges()
	if edges[0] != node1 {
		t.FailNow()
	}

	edges[0] = nil

	edges2 := *node2.GetEdges()
	if edges2[0] != nil {
		t.FailNow()
	}
}