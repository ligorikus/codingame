package main

import (
	"fmt"
)

type Graph [][]int

type Path struct {
	path     []int
	cost     int
	startKey int
	endKey   int
	length   int
}

func MinKey(min int, arr []int) int {
	minKey := -1
	for i := 0; i < len(arr); i++ {
		if arr[i] < min {
			minKey = i
			min = arr[i]
		}
	}
	return minKey
}

func RelationWithMinCost(graph Graph, path Path) Path {
	result := make([]int, len(graph))
	resultPath := make([]Path, len(graph))
	for i := 0; i < len(result); i++ {
		result[i] = 1000
	}
	for node := 0; node < len(graph); node++ {
		continueFlag := false
		for j := 0; j < path.length; j++ {
			if node == path.path[j] {
				continueFlag = true
				break
			}
		}

		if continueFlag {
			continue
		}
		calcPath := Path{}
		calcPath.path = make([]int, len(graph))
		calcPath.length = path.length
		calcPath.startKey = path.startKey
		calcPath.endKey = path.endKey
		copy(calcPath.path, path.path)

		rightSidePath := make([]int, len(graph))
		copy(rightSidePath, calcPath.path[calcPath.endKey:])
		calcPath.path[calcPath.endKey] = node
		for i := 0; i < len(rightSidePath); i++ {
			if rightSidePath[i] == -1 {
				break
			}
			calcPath.path[calcPath.endKey+i+1] = rightSidePath[i]
		}

		calcPath.length++

		calcPath.cost = GetCostPath(graph, calcPath)
		if calcPath.cost < path.cost {
			result[node] = calcPath.cost
			resultPath[node] = calcPath
		}
	}
	minKey := MinKey(path.cost, result)

	if minKey != -1 {
		leftPath := Path{}
		leftPath.path = make([]int, len(graph))
		leftPath.length = resultPath[minKey].length
		leftPath.cost = resultPath[minKey].cost
		leftPath.startKey = resultPath[minKey].startKey
		leftPath.endKey = resultPath[minKey].startKey + 1
		copy(leftPath.path, resultPath[minKey].path)
		path = RelationWithMinCost(graph, leftPath)

		rightPath := Path{}
		rightPath.path = make([]int, len(graph))
		rightPath.length = resultPath[minKey].length
		rightPath.startKey = resultPath[minKey].cost
		rightPath.startKey = resultPath[minKey].endKey
		rightPath.endKey = resultPath[minKey].endKey + 1
		copy(rightPath.path, resultPath[minKey].path)
		path = RelationWithMinCost(graph, rightPath)
	}
	return path
}

func GetCostPath(graph Graph, path Path) int {
	cost := 0
	for i := 1; i < path.length; i++ {
		cost += graph[path.path[i-1]][path.path[i]]
	}
	return cost
}

func main() {
	graph := Graph{
		{0, 3, 3, 3, 3, 7, 7, 2, 2, 6, 6, 1, 1, 8, 8},
		{3, 0, 7, 2, 7, 6, 9, 5, 2, 4, 9, 1, 5, 8, 9},
		{3, 7, 0, 7, 2, 9, 6, 2, 5, 9, 4, 5, 1, 9, 8},
		{3, 2, 7, 0, 8, 3, 11, 4, 5, 1, 10, 1, 5, 5, 11},
		{3, 7, 2, 8, 0, 11, 3, 5, 4, 10, 1, 5, 1, 11, 5},
		{7, 6, 9, 3, 11, 0, 15, 5, 9, 2, 14, 5, 9, 1, 15},
		{7, 9, 6, 11, 3, 15, 0, 9, 5, 14, 2, 9, 5, 15, 1},
		{2, 5, 2, 4, 5, 5, 9, 0, 5, 6, 7, 2, 3, 5, 10},
		{2, 2, 5, 5, 4, 9, 5, 5, 0, 7, 6, 3, 2, 10, 5},
		{6, 4, 9, 1, 10, 2, 14, 6, 7, 0, 13, 3, 8, 4, 14},
		{6, 9, 4, 10, 1, 14, 2, 7, 6, 13, 0, 8, 3, 14, 4},
		{1, 1, 5, 1, 5, 5, 9, 2, 3, 3, 8, 0, 3, 6, 9},
		{1, 5, 1, 5, 1, 9, 5, 3, 2, 8, 3, 3, 0, 9, 6},
		{8, 8, 9, 5, 11, 1, 15, 5, 10, 4, 14, 6, 9, 0, 17},
		{8, 9, 8, 11, 5, 15, 1, 10, 5, 14, 4, 9, 6, 17, 0},
	}

	pathArr := make([]int, len(graph))
	for i := 2; i < len(pathArr); i++ {
		pathArr[i] = -1
	}
	pathArr[0] = 9
	pathArr[1] = 0

	path := Path{
		path:     pathArr,
		cost:     graph[pathArr[0]][pathArr[1]],
		startKey: 0,
		endKey:   1,
		length:   2,
	}

	result := RelationWithMinCost(graph, path)
	fmt.Println(result.path[:result.length])
}
