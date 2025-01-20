package main

import (
	"container/list"
	"fmt"
	"math"
	"os"
)

const (
	MaxEdgeWeight = 9    // Maximum edge weight for the graph
	NullParent    = -1   // Indicator for no parent
)

type Cell struct {
	Row, Col int64
}

// Check if the start and end cells are the same
func isSameCell(startRow, startCol, endRow, endCol int64) bool {
	return startRow == endRow && startCol == endCol
}

// Print the path from start to end
func printPath(path []Cell) {
	for i := len(path) - 1; i >= 0; i-- {
		fmt.Println(path[i].Row, path[i].Col)
	}
	fmt.Println(".")
}

// Build a graph representation from the grid weights
func buildGraph(rows, cols int64, weights []int8, graph [][]int64) {
	for i := int64(0); i < rows; i++ {
		for j := int64(0); j < cols; j++ {
			currentIndex := cols*i + j
			if weights[currentIndex] == 0 {
				continue // Skip impassable cells
			}
			// Add edges for valid neighbors
			if i > 0 && weights[currentIndex-cols] != 0 {
				graph[currentIndex] = append(graph[currentIndex], currentIndex-cols)
			}
			if i < rows-1 && weights[currentIndex+cols] != 0 {
				graph[currentIndex] = append(graph[currentIndex], currentIndex+cols)
			}
			if j > 0 && weights[currentIndex-1] != 0 {
				graph[currentIndex] = append(graph[currentIndex], currentIndex-1)
			}
			if j < cols-1 && weights[currentIndex+1] != 0 {
				graph[currentIndex] = append(graph[currentIndex], currentIndex+1)
			}
		}
	}
}

// Reconstruct the path from end to start
func buildPath(cols, startRow, startCol, endRow, endCol int64, parents []int64, path *[]Cell) bool {
	for endRow != startRow || endCol != startCol {
		*path = append(*path, Cell{endRow, endCol})
		parentIndex := parents[endRow*cols+endCol]
		if parentIndex == NullParent {
			return false // Path does not exist
		}
		endRow = parentIndex / cols
		endCol = parentIndex % cols
	}
	*path = append(*path, Cell{startRow, startCol})
	return true
}

// Modified BFS to handle weighted graphs using multi-level queues
func bfs(graph [][]int64, distances []int64, levels []list.List, visited []bool, weights []int8, parents []int64) {
	level := 0
	activeNodes := 1

	for activeNodes > 0 {
		// Skip empty queues
		for levels[level%(MaxEdgeWeight+1)].Len() == 0 {
			level++
		}

		element := levels[level%(MaxEdgeWeight+1)].Front()
		currentNode := element.Value.(int64)
		levels[level%(MaxEdgeWeight+1)].Remove(element)

		activeNodes--

		if visited[currentNode] {
			continue
		}

		visited[currentNode] = true
		for _, neighbor := range graph[currentNode] {
			neighborWeight := weights[neighbor]
			if distances[currentNode]+int64(neighborWeight) < distances[neighbor] {
				distances[neighbor] = distances[currentNode] + int64(neighborWeight)
				parents[neighbor] = currentNode
				levels[distances[neighbor]%(MaxEdgeWeight+1)].PushBack(neighbor)
				activeNodes++
			}
		}
	}
}

func main() {
	var rows, cols int64
	fmt.Scan(&rows, &cols)

	weights := make([]int8, rows*cols)
	parents := make([]int64, rows*cols)
	distances := make([]int64, rows*cols)
	visited := make([]bool, rows*cols)
	graph := make([][]int64, rows*cols)

	// Initialize distances and parents
	for i := range distances {
		distances[i] = math.MaxInt64
		parents[i] = NullParent
	}

	for i := int64(0); i < rows*cols; i++ {
		fmt.Scan(&weights[i])
	}

	var startRow, startCol, endRow, endCol int64
	fmt.Scan(&startRow, &startCol, &endRow, &endCol)

	// Validate input coordinates
	if startRow < 0 || startRow >= rows || endRow < 0 || endCol >= cols {
		fmt.Fprintln(os.Stderr, "Invalid input")
		os.Exit(1)
	}

	// If the start and end cells are the same, print the path
	if isSameCell(startRow, startCol, endRow, endCol) {
		printPath([]Cell{{startRow, startCol}})
		os.Exit(0)
	}

	// Build the graph
	buildGraph(rows, cols, weights, graph)

	// Set up the starting point
	startIndex := startRow*cols + startCol
	parents[startIndex] = startIndex
	distances[startIndex] = 0

	levels := make([]list.List, MaxEdgeWeight+1)
	levels[0].PushBack(startIndex)

	// Perform BFS
	bfs(graph, distances, levels, visited, weights, parents)

	// Construct and print the path
	var path []Cell
	if !buildPath(cols, startRow, startCol, endRow, endCol, parents, &path) {
		fmt.Fprintln(os.Stderr, "End cell is unreachable")
		os.Exit(1)
	}

	printPath(path)
}
