package main

import (
	"fmt"
	"math"
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
	incoming   [Deep][]TroopSquad
	properties NodeProperties
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

type NodeProperties struct {
	balance [Deep]int
	owner   [Deep]Owner
	cost    []int
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
		state.nodes[i] = &Node{id: i, properties: NodeProperties{cost: make([]int, len(state.graph))}}
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

func (node *Node) CalculateBalance() {
	node.properties.balance[0] = node.garrison
	node.properties.owner[0] = node.owner

	for i := 1; i < Deep; i++ {
		node.properties.owner[i] = node.properties.owner[i-1]
		node.properties.balance[i] = node.properties.balance[i-1]

		if node.properties.owner[i] != Neutral {
			node.properties.balance[i] = node.properties.balance[i-1] + node.production
		}

		incomingBalance := 0
		var incomingWinner Owner
		for _, troop := range node.incoming[i] {
			incomingBalance += troop.size * int(troop.owner)
		}
		if incomingBalance != 0 {
			if incomingBalance > 0 {
				incomingWinner = Player
			} else {
				incomingWinner = Enemy
				incomingBalance *= -1
			}

			if node.properties.owner[i] == incomingWinner {
				node.properties.balance[i] += incomingBalance
			} else {
				node.properties.balance[i] -= incomingBalance
			}

			if node.properties.balance[i] < 0 {
				node.properties.balance[i] *= -1
				node.properties.owner[i] = incomingWinner
			}
		}
	}
}

func (node *Node) CalculateCost(state GameState) {
	if node.owner != Player {
		return
	}
	directions := state.graph[node.id]
	for i := 0; i < len(directions); i++ {
		if i == node.id {
			continue
		}

		currentNode := state.nodes[i]
		troopSize := currentNode.properties.balance[directions[i]]

		var payback float64
		if currentNode.production != 0 {
			payback = float64(troopSize) / float64(currentNode.production)
			node.properties.cost[i] = int(math.Ceil(payback)) + directions[i]
		} else {
			node.properties.cost[i] = 200
		}
	}
}

func (node *Node) SearchCheap() int {
	min := 999
	minCostNode := 0
	for i := 0; i < len(node.properties.cost); i++ {
		if i == node.id {
			continue
		}
		if node.properties.cost[i] < min {
			minCostNode = i
			min = node.properties.cost[i]
		}
	}
	return minCostNode
}

func main() {
	state := GameState{}
	state.InitializeGraph()
	state.InitializeNodes()

	for {
		start := time.Now()
		state.UpdateState()

		for _, node := range state.nodes {
			node.CalculateBalance()
		}
		for _, node := range state.nodes {
			node.CalculateCost(state)
		}

		fmt.Print("WAIT")
		for _, node := range state.nodes {
			if node.owner == Player {
				fmt.Print(";")
				targetNodeId := node.SearchCheap()
				distance := state.graph[node.id][targetNodeId]
				size := state.nodes[targetNodeId].properties.balance[distance] + 1
				fmt.Print("MOVE ", node.id, targetNodeId, size)
			}
		}
		fmt.Println()
		// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"

		duration := time.Since(start)
		fmt.Fprintln(os.Stderr, duration)
	}
}
