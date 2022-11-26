package apiserver

import (
	"container/list"
	// "fmt"
)

func topologicalSort(nodeID int, nodeLen int, visited []bool, revStack *list.List, Nodes []TemplateItem, nameToIndex map[string]int, outputsinks map[string]bool) {
	visited[nodeID] = true

	neighNodes := Nodes[nodeID].Output

	for _, neighNode := range neighNodes {
		if !outputsinks[neighNode.Name] {
			neighNodeID := nameToIndex[neighNode.Name]
			if !visited[neighNodeID] {
				topologicalSort(neighNodeID, nodeLen-1, visited[:], revStack, Nodes, nameToIndex, outputsinks)
			}
		}
	}

	revStack.PushFront(nodeID)

}
