package apiserver

import (
	"container/list"
	// "fmt"
)

func topologicalSort(nodeID int, nodeLen int, visited []bool, revStack *list.List, Nodes []TemplateItem, nameToIndex map[string]int) {
	visited[nodeID] = true

	neighNodes := Nodes[nodeID].Output

	for _, neighNode := range neighNodes {
		neighNodeID := nameToIndex[neighNode.Name]
		if !visited[neighNodeID] {
			topologicalSort(neighNodeID, nodeLen-1, visited[:], revStack, Nodes, nameToIndex)
		}
	}

	revStack.PushFront(nodeID)

}
