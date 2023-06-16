package main

import (
	"fmt"
)

type Owner int64

const (
	Player  Owner = 1
	Enemy         = -1
	Neutral       = 0
)

type EntityType string

type Graph [][]int

const (
	Factory EntityType = "FACTORY"
	Troop   EntityType = "TROOP"
)

type Node struct {
	id         int
	owner      Owner
	garrison   int
	production int
}

type TroopSquad struct {
	owner       Owner
	source      Node
	destination Node
	size        int
	arrival     int
}

type GameState struct {
	graph  Graph
	nodes  []Node
	troops []TroopSquad
}

func Initialize() Graph {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	graph := make(Graph, factoryCount)
	for i := range graph {
		graph[i] = make([]int, factoryCount)
	}

	for i := 0; i < linkCount; i++ {
		var factory1, factory2, distance int
		fmt.Scan(&factory1, &factory2, &distance)
		graph[factory1][factory2] = distance
		graph[factory2][factory1] = distance
	}
	return graph
}

func (state *GameState) UpdateState() {
	var entityCount int
	fmt.Scan(&entityCount)
	state.nodes = state.nodes[:0]
	state.troops = state.troops[:0]

	for i := 0; i < entityCount; i++ {
		var entityId int
		var entityType EntityType
		var owner Owner
		var arg2, arg3, arg4, arg5 int
		fmt.Scan(&entityId, &entityType, &owner, &arg2, &arg3, &arg4, &arg5)

		switch entityType {
		case Factory:
			node := Node{id: entityId, owner: owner, garrison: arg2, production: arg3}
			state.nodes = append(state.nodes, node)
		case Troop:
			troop := TroopSquad{owner: owner, source: state.nodes[arg2], destination: state.nodes[arg3], size: arg4, arrival: arg5}
			state.troops = append(state.troops, troop)
		}
	}
}

func main() {
	graph := Initialize()

	var nodes []Node
	var troops []TroopSquad
	state := GameState{graph: graph, nodes: nodes, troops: troops}

	for {
		state.UpdateState()
		// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
		fmt.Println("WAIT")
	}
}
