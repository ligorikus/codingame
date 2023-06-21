package main

import (
	"fmt"
	"os"
	"time"
)

type Graph1 [][]int

const MaxDistanceNum = 64

func Dijkstra(graph Graph1, startNode int) ([]int, []int) {
	u := make([]bool, len(graph))
	d := make([]int, len(graph))
	p := make([]int, len(graph))
	for i := 0; i < len(d); i++ {
		d[i] = MaxDistanceNum
	}
	d[startNode] = 0

	for i := 0; i < len(graph); i++ {
		v := -1
		for j := 0; j < len(graph); j++ {
			if !u[j] && (v == -1 || d[j] < d[v]) {
				v = j
			}
		}

		if d[v] == MaxDistanceNum {
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

func InitMoveMatrix(graph Graph1) [][]int {
	length := len(graph)

	moveMatrix := make([][]int, length)
	for i := 0; i < length; i++ {
		moveMatrix[i] = make([]int, length)
		for j := 0; j < length; j++ {
			moveMatrix[i][j] = -1
		}
	}

	for i := 0; i < length; i++ {
		_, path := Dijkstra(graph, i)
		for j := 0; j < length; j++ {
			if i != j {
				restored := RestoreDijkstra(length, i, j, path)
				moveMatrix = UpdateMoveMatrix(moveMatrix, restored)
			}
		}
	}
	return moveMatrix
}

func main() {
	start := time.Now()
	graph := Graph1{
		{0, 4, 4, 8, 8, 7, 7, 3, 3, 1, 1, 2, 2, 3, 3},
		{4, 0, 10, 3, 13, 2, 13, 2, 8, 4, 5, 1, 7, 1, 9},
		{4, 10, 0, 13, 3, 13, 2, 8, 2, 5, 4, 7, 1, 9, 1},
		{8, 3, 13, 0, 17, 3, 16, 3, 12, 7, 9, 4, 11, 5, 12},
		{8, 13, 3, 17, 0, 16, 3, 12, 3, 9, 7, 11, 4, 12, 5},
		{7, 2, 13, 3, 16, 0, 16, 5, 11, 8, 8, 4, 10, 3, 12},
		{7, 13, 2, 16, 3, 16, 0, 11, 5, 8, 8, 10, 4, 12, 3},
		{3, 2, 8, 3, 12, 5, 11, 0, 8, 2, 5, 1, 6, 4, 7},
		{3, 8, 2, 12, 3, 11, 5, 8, 0, 5, 2, 6, 1, 7, 4},
		{1, 4, 5, 7, 9, 8, 8, 2, 5, 0, 3, 2, 3, 5, 3},
		{1, 5, 4, 9, 7, 8, 8, 5, 2, 3, 0, 3, 2, 3, 5},
		{2, 1, 7, 4, 11, 4, 10, 1, 6, 2, 3, 0, 5, 2, 6},
		{2, 7, 1, 11, 4, 10, 4, 6, 1, 3, 2, 5, 0, 6, 2},
		{3, 1, 9, 5, 12, 3, 12, 4, 7, 5, 3, 2, 6, 0, 8},
		{3, 9, 1, 12, 5, 12, 3, 7, 4, 3, 5, 6, 2, 8, 0},
	}

	moveMatrix := InitMoveMatrix(graph)
	duration := time.Since(start)
	fmt.Fprintln(os.Stderr, duration)
	for _, item := range moveMatrix {
		fmt.Println(item)
	}
}
