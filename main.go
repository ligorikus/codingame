package main

import (
	"fmt"
	"os"
	"time"
)

const MaxDistance = 20
const Deep = MaxDistance + 1

type Owner int64

type GameInput struct {
	owner Owner
	arg2  int
	arg3  int
	arg4  int
	arg5  int
}

const (
	Player  Owner = 1
	Enemy         = -1
	Neutral       = 0
)

type EntityType string

type Graph1 [][]int

const (
	Factory EntityType = "FACTORY"
	Troop   EntityType = "TROOP"
)

type Node struct {
	id         int
	owner      Owner
	garrison   int
	production int
	incoming   [Deep][]TroopSquad
}

type TroopSquad struct {
	owner       Owner
	source      int
	destination int
	size        int
	arrival     int
}

type GameState struct {
	graph Graph
	nodes []*Node
}

func (state *GameState) InitializeGraph() {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	state.graph = make(Graph, factoryCount)
	for i := range state.graph {
		state.graph[i] = make([]int, factoryCount)
	}

	for i := 0; i < linkCount; i++ {
		var factory1, factory2, distance int
		fmt.Scan(&factory1, &factory2, &distance)
		state.graph[factory1][factory2] = distance
		state.graph[factory2][factory1] = distance
	}
}

func (state *GameState) InitializeNodes() {
	state.nodes = make([]*Node, len(state.graph))
	for i := range state.graph {
		state.nodes[i] = &Node{id: i}
	}
}

func (state *GameState) UpdateState() {
	var entityCount int
	fmt.Scan(&entityCount)

	for i := 0; i < entityCount; i++ {
		var entityId int
		var entityType EntityType
		var owner Owner
		var arg2, arg3, arg4, arg5 int
		fmt.Scan(&entityId, &entityType, &owner, &arg2, &arg3, &arg4, &arg5)
		input := GameInput{owner: owner, arg2: arg2, arg3: arg3, arg4: arg4, arg5: arg5}

		switch entityType {
		case Factory:
			state.UpdateNode(entityId, input)
		case Troop:
			state.InputTroop(input)
		}
	}
}

func (state *GameState) UpdateNode(nodeId int, input GameInput) {
	node := state.nodes[nodeId]
	node.owner = input.owner
	node.garrison = input.arg2
	node.production = input.arg3
	node.incoming = [Deep][]TroopSquad{}
}

func (state *GameState) InputTroop(input GameInput) {
	troop := TroopSquad{}
	troop.owner = input.owner
	troop.source = input.arg2
	troop.destination = input.arg3
	troop.size = input.arg4
	troop.arrival = input.arg5

	nodeIncoming := state.nodes[troop.destination].incoming[troop.arrival]
	state.nodes[troop.destination].incoming[troop.arrival] = append(nodeIncoming, troop)
}

func (state *GameState) PrintGraphMatrix() {
	for _, graph := range state.graph {
		fmt.Fprintln(os.Stderr, graph)
	}
}

func main2() {
	state := GameState{}
	state.InitializeGraph()

	start := time.Now()
	state.InitializeNodes()

	state.PrintGraphMatrix()

	duration := time.Since(start)
	fmt.Fprintln(os.Stderr, duration)
	for {
		start := time.Now()
		state.UpdateState()

		fmt.Print("WAIT")

		fmt.Println()
		// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"

		duration := time.Since(start)
		fmt.Fprintln(os.Stderr, duration)
	}
}
