package main

import "fmt"

type Graph1 [][]int

const MaxDistanceNum = 64

func MinKeyValue(arr []int, exclude []bool) (int, int) {
	min := MaxDistanceNum
	minKey := -1
	for i, item := range arr {
		if exclude[i] {
			continue
		}
		if item < min {
			minKey = i
			min = item
		}
	}
	return minKey, min
}

func Dijkstra(graph Graph1, movementTable [][]int, startNode int) [][]int {
	var dijkstraTable, dijkstraTableUtil [][]int
	dijkstraTable = make([][]int, len(graph))
	dijkstraTableUtil = make([][]int, len(graph))
	visited := make([]bool, len(graph))

	for i := 0; i < len(graph); i++ {
		dijkstraTable[i] = make([]int, len(graph))
		dijkstraTableUtil[i] = make([]int, len(graph))
		for j := 0; j < len(movementTable); j++ {
			dijkstraTable[i][j] = MaxDistanceNum
			dijkstraTableUtil[i][j] = -1
		}
	}
	dijkstraTable[0][startNode] = 0
	dijkstraTableUtil[0][startNode] = startNode

	for i := 1; i < len(dijkstraTable); i++ {
		key, cost := MinKeyValue(dijkstraTable[i-1], visited)
		if key == -1 {
			break
		}

		visited[key] = true
		for j := 0; j < len(dijkstraTable[i]); j++ {
			if visited[j] {
				dijkstraTable[i][j] = dijkstraTable[i-1][j]
				dijkstraTableUtil[i][j] = dijkstraTableUtil[i-1][j]
				continue
			}

			oldCost := dijkstraTable[i-1][j]
			newCost := cost + graph[key][j]
			if newCost < oldCost {
				dijkstraTable[i][j] = newCost
				dijkstraTableUtil[i][j] = key
			} else {
				dijkstraTable[i][j] = oldCost
				dijkstraTableUtil[i][j] = dijkstraTableUtil[i-1][j]
			}
		}
	}

	for i, item := range dijkstraTableUtil[len(graph)-1] {
		if i == item {
			continue
		}
		movementTable[i][item] = item
		movementTable[item][i] = i

		var inversedI, inversedItem int
		if i != 0 {
			if i%2 == 0 {
				inversedI = i - 1
			} else {
				inversedI = i + 1
			}
		}
		if item != 0 {
			if item%2 == 0 {
				inversedItem = i - 1
			} else {
				inversedItem = i + 1
			}
		}

		if inversedI == inversedItem {
			continue
		}
		movementTable[inversedI][inversedItem] = item
		movementTable[inversedItem][inversedI] = i
	}

	return movementTable
}

func main() {
	graph := Graph1{
		{0, 6, 6, 3, 3, 5, 5, 1, 1},
		{6, 0, 14, 2, 10, 1, 13, 4, 9},
		{6, 14, 0, 10, 2, 13, 1, 9, 4},
		{3, 2, 10, 0, 7, 2, 9, 1, 5},
		{3, 10, 2, 7, 0, 9, 2, 5, 1},
		{5, 1, 13, 2, 9, 0, 12, 2, 8},
		{5, 13, 1, 9, 2, 12, 0, 8, 2},
		{1, 4, 9, 1, 5, 2, 8, 0, 4},
		{1, 9, 4, 5, 1, 8, 2, 4, 0},
	}

	var movementTable [][]int

	movementTable = make([][]int, len(graph))
	for i := 0; i < len(movementTable); i++ {
		movementTable[i] = make([]int, len(graph))
		for j := 0; j < len(movementTable); j++ {
			movementTable[i][j] = -1
		}
	}

	movementTable = Dijkstra(graph, movementTable, 1)

	for _, item := range movementTable {
		fmt.Println(item)
	}
}
