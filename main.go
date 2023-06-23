package main

import (
	"fmt"
	"os"
	"time"
)

const MaxDeep = 21

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

type Matrix [][]int
type Graph Matrix

const (
	Factory EntityType = "FACTORY"
	Troop   EntityType = "TROOP"
)

type Node struct {
	id         int
	owner      Owner
	garrison   int
	production int
	incoming   [MaxDeep][]TroopSquad
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
	graph      Graph
	nodes      []*Node
	costMatrix Matrix
	moveMatrix Matrix
}

type NodeProperties struct {
	balance [MaxDeep]int
	owner   [MaxDeep]Owner
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

	costMatrix, moveMatrix := InitMoveMatrix(state.graph)
	state.costMatrix = costMatrix
	state.moveMatrix = moveMatrix
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
	node.incoming = [MaxDeep][]TroopSquad{}
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

func Dijkstra(graph Graph, startNode int) ([]int, []int) {
	u := make([]bool, len(graph))
	d := make([]int, len(graph))
	p := make([]int, len(graph))
	for i := 0; i < len(d); i++ {
		d[i] = MaxDeep
	}
	d[startNode] = 0

	for i := 0; i < len(graph); i++ {
		v := -1
		for j := 0; j < len(graph); j++ {
			if !u[j] && (v == -1 || d[j] < d[v]) {
				v = j
			}
		}

		if d[v] == MaxDeep {
			break
		}

		u[v] = true

		for j := 0; j < len(graph[v]); j++ {
			to := j
			len := graph[v][j]
			if d[v]+len < d[to] {
				d[to] = d[v] + len
				p[to] = v
			}
		}
	}

	return d, p
}

func RestoreDijkstra(length int, startNode int, finishNode int, p []int) []int {
	path := make([]int, length)
	for v := finishNode; v != startNode; v = p[v] {
		path = append(path, v)
	}
	path = append(path, startNode)

	rev_slc := []int{}
	for i := range path {
		rev_slc = append(rev_slc, path[len(path)-1-i])
	}

	j := 0
	for i := range rev_slc {
		j++
		if rev_slc[i] == finishNode {
			break
		}
	}
	return rev_slc[:j]
}

func UpdateMoveMatrix(moveMatrix [][]int, restored []int) [][]int {
	target := restored[len(restored)-1]
	for i := 0; i < len(restored)-1; i++ {
		moveMatrix[restored[i]][target] = restored[i+1]
	}
	return moveMatrix
}

func InitMoveMatrix(graph Graph) ([][]int, [][]int) {
	length := len(graph)

	moveMatrix := make([][]int, length)
	costMatrix := make([][]int, length)
	for i := 0; i < length; i++ {
		moveMatrix[i] = make([]int, length)
		costMatrix[i] = make([]int, length)
		for j := 0; j < length; j++ {
			moveMatrix[i][j] = -1
		}
	}

	for i := 0; i < length; i++ {
		cost, path := Dijkstra(graph, i)
		costMatrix[i] = cost
		for j := 0; j < length; j++ {
			if i != j {
				restored := RestoreDijkstra(length, i, j, path)
				moveMatrix = UpdateMoveMatrix(moveMatrix, restored)
			}
		}
	}
	return costMatrix, moveMatrix
}

func main() {
	start := time.Now()
	state := GameState{}
	state.InitializeGraph()
	state.InitializeNodes()

	for _, item := range state.graph {
		fmt.Fprintln(os.Stderr, item)
	}

	duration := time.Since(start)
	fmt.Fprintln(os.Stderr, duration)

	for {
		state.UpdateState()

		fmt.Println("WAIT")
		// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
	}
}
