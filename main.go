package main

import (
	"fmt"
	"os"
)

type Owner int64

const (
	Player  Owner = 1
	Enemy         = -1
	Neutral       = 0
)

func initialize() [][]int {
	// factoryCount: the number of factories
	var factoryCount int
	fmt.Scan(&factoryCount)

	// linkCount: the number of links between factories
	var linkCount int
	fmt.Scan(&linkCount)

	graph := make([][]int, factoryCount)

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

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

func main() {
	graph := initialize()
	for i := 0; i < len(graph); i++ {
		fmt.Fprintln(os.Stderr, graph[i])
	}

	for {
		// entityCount: the number of entities (e.g. factories and troops)
		var entityCount int
		fmt.Scan(&entityCount)

		for i := 0; i < entityCount; i++ {
			var entityId int
			var entityType string
			var arg1, arg2, arg3, arg4, arg5 int
			fmt.Scan(&entityId, &entityType, &arg1, &arg2, &arg3, &arg4, &arg5)
		}

		// Any valid action, such as "WAIT" or "MOVE source destination cyborgs"
		fmt.Println("WAIT")
	}
}
