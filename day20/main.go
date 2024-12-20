package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"time"
)

type Grid [141][141]byte

type Position struct {
	x, y int
}

type Queue struct {
	nodes    [100000]Node // Fixed-size array for better memory characteristics
	head     int          // Index of the first element
	tail     int          // Index where the next element will be inserted
	size     int          // Current number of elements
	capacity int          // Maximum capacity
}

// NewQueue creates a new queue with the specified capacity
func NewQueue() *Queue {
	return &Queue{
		nodes:    [100000]Node{},
		head:     0,
		tail:     0,
		size:     0,
		capacity: 100000,
	}
}

// Push adds a node to the queue
func (q *Queue) Push(n Node) bool {
	if q.size >= q.capacity {
		return false // Queue is full
	}

	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return true
}

// Pop removes and returns the first node in the queue
func (q *Queue) Pop() (Node, bool) {
	if q.size == 0 {
		return Node{}, false // Queue is empty
	}

	node := q.nodes[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return node, true
}

// Peek returns the first node without removing it
func (q *Queue) Peek() (Node, bool) {
	if q.size == 0 {
		return Node{}, false
	}
	return q.nodes[q.head], true
}

// IsEmpty returns true if the queue has no elements
func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

// Size returns the current number of elements in the queue
func (q *Queue) Size() int {
	return q.size
}

// Clear resets the queue to empty state
func (q *Queue) Clear() {
	q.head = 0
	q.tail = 0
	q.size = 0
}

type Node struct {
	path        []Position
	position    Position
	cheatingPos []Position
}

func (g Grid) Display() {
	for y := range g {
		for x := range g[y] {
			fmt.Print(string(g[y][x]))
		}
		fmt.Println()
	}
}

func (g Grid) VisualizeSolution(input []Position) {
	fmt.Println("Visualizer", len(input))
	fmt.Println(input)
	for y := range g {
		for x := range g[y] {
			isPath := false
			for _, pos := range input {
				if pos.x == x && pos.y == y {
					isPath = true
					break
				}
			}
			if isPath {
				fmt.Print("O")
			} else {
				fmt.Print(string(g[y][x]))
			}
		}
		fmt.Println()
	}
}

func (g Grid) GetStartAndEndPosition() (Position, Position) {
	var start, end Position
	for y := range g {
		for x := range g[y] {
			if g[y][x] == 'S' {
				start = Position{x: x, y: y}
			}
			if g[y][x] == 'E' {
				end = Position{x: x, y: y}
			}
		}
	}
	return start, end
}

//go:embed "day20input.txt"
var inputFile []byte

func main() {
	start := time.Now()
	grid := Grid{}

	// Parse input grid
	grid = fillGrid(inputFile, grid)

	p1 := p1(grid)

	// for _, node := range p1 {
	// 	grid.VisualizeSolution(node.path)
	// }

	partOne := time.Since(start)
	fmt.Println("Part 1:", p1, "in", partOne)
}

func fillGrid(inputFile []byte, grid Grid) Grid {
	lines := bytes.Split(inputFile, []byte("\n"))
	for y, line := range lines {
		for x, byte := range line {
			grid[y][x] = byte
		}
	}
	return grid
}

func createNodeKey(node Node) string {
	var nextNodekey string
	if len(node.cheatingPos) == 2 {
		nextNodekey = fmt.Sprintf("%d,%d,%d,%d,%d,%d", node.position.x, node.position.y, node.cheatingPos[0].x, node.cheatingPos[0].y, node.cheatingPos[1].x, node.cheatingPos[1].y)
	} else if len(node.cheatingPos) == 1 {
		nextNodekey = fmt.Sprintf("%d,%d,%d,%d", node.position.x, node.position.y, node.cheatingPos[0].x, node.cheatingPos[0].y)
	} else {
		nextNodekey = fmt.Sprintf("%d,%d", node.position.x, node.position.y)
	}
	return nextNodekey
}

func p1(grid Grid) int {
	startPos, endPos := grid.GetStartAndEndPosition()
	nodeThatEnded := make([]Node, 0)
	baseLinePath := findNoCheating(grid)
	baseLinePico := len(baseLinePath[0].path)
	directions := [][]int{
		{0, 1},  // bottom
		{0, -1}, // up
		{1, 0},  // right
		{-1, 0}, // left
	}
	// count := 0
	startNode := Node{path: make([]Position, 0), position: startPos, cheatingPos: make([]Position, 0)}
	queue := NewQueue()
	visited := make(map[string]bool, 2000)

	key := fmt.Sprintf("%d,%d", startNode.position.x, startNode.position.y)
	queue.Push(startNode)
	visited[key] = true

	for queue.IsEmpty() {
		v, _ := queue.Pop()
		if v.position.x == endPos.x && v.position.y == endPos.y {
			if baseLinePico > len(v.path) {
				nodeThatEnded = append(nodeThatEnded, v)
			}
			continue
		}
		if grid[v.position.y][v.position.x] == '#' && len(v.cheatingPos) == 2 {
			continue
		}
		for _, direction := range directions {
			nextPosition := Position{x: v.position.x + direction[0], y: v.position.y + direction[1]}
			if slices.Contains(v.cheatingPos, nextPosition) {
				continue
			}

			if nextPosition.x < 0 || nextPosition.x > len(grid[0])-1 || nextPosition.y < 0 || nextPosition.y > len(grid)-1 {
				continue
			}

			nextNode := Node{position: nextPosition}
			nextNode.path = make([]Position, len(v.path))
			nextNode.cheatingPos = make([]Position, len(v.cheatingPos))
			copy(nextNode.cheatingPos, v.cheatingPos)
			copy(nextNode.path, v.path)
			nextNode.path = append(nextNode.path, nextPosition)

			nextNodekey := createNodeKey(nextNode)
			if visited[nextNodekey] == false && grid[nextPosition.y][nextPosition.x] != '#' && len(nextNode.path) < baseLinePico {
				queue.Push(nextNode)
				visited[nextNodekey] = true
			}

			if len(v.cheatingPos) < 2 && grid[nextPosition.y][nextPosition.x] == '#' {
				// Cheating
				nextCheatingNode := Node{position: nextPosition}
				nextCheatingNode.path = make([]Position, len(v.path))
				nextCheatingNode.cheatingPos = make([]Position, len(v.cheatingPos))
				copy(nextCheatingNode.path, v.path)
				copy(nextCheatingNode.cheatingPos, v.cheatingPos)
				nextCheatingNode.path = append(nextCheatingNode.path, nextPosition)
				nextCheatingNode.cheatingPos = append(nextCheatingNode.cheatingPos, nextPosition)

				nextCheatingNodekey := createNodeKey(nextCheatingNode)
				if visited[nextCheatingNodekey] == false && len(nextCheatingNode.path) < baseLinePico {
					queue.Push(nextCheatingNode)
					visited[nextCheatingNodekey] = true
				}
			}
		}
	}

	count := 0
	timeSaved := 100
	for _, node := range nodeThatEnded {
		if len(node.path) < baseLinePico-timeSaved {
			count++
		}
	}

	return count
}

func findNoCheating(grid Grid) []Node {
	startPos, endPos := grid.GetStartAndEndPosition()
	nodeThatEnded := make([]Node, 0)
	directions := [][]int{
		{0, 1},  // bottom
		{0, -1}, // up
		{1, 0},  // right
		{-1, 0}, // left
	}
	// count := 0
	startNode := Node{path: make([]Position, 0), position: startPos, cheatingPos: make([]Position, 0)}
	queue := NewQueue()
	visited := make(map[string]bool, 1000)

	key := fmt.Sprintf("%d,%d", startNode.position.x, startNode.position.y)
	queue.Push(startNode)
	visited[key] = true

	for !queue.IsEmpty() {
		v, _ := queue.Pop()
		if v.position.x == endPos.x && v.position.y == endPos.y {
			nodeThatEnded = append(nodeThatEnded, v)
			continue
		}
		if grid[v.position.y][v.position.x] == '#' && len(v.cheatingPos) == 2 {
			continue
		}
		for _, direction := range directions {
			nextPosition := Position{x: v.position.x + direction[0], y: v.position.y + direction[1]}
			if slices.Contains(v.cheatingPos, nextPosition) {
				continue
			}

			if nextPosition.x < 0 || nextPosition.x > len(grid[0])-1 || nextPosition.y < 0 || nextPosition.y > len(grid)-1 {
				continue
			}

			nextNode := Node{position: nextPosition}
			nextNode.path = make([]Position, len(v.path))
			nextNode.cheatingPos = make([]Position, len(v.cheatingPos))
			copy(nextNode.cheatingPos, v.cheatingPos)
			copy(nextNode.path, v.path)
			nextNode.path = append(nextNode.path, nextPosition)

			nextNodekey := createNodeKey(nextNode)
			if visited[nextNodekey] == false && grid[nextPosition.y][nextPosition.x] != '#' {
				queue.Push(nextNode)
				visited[nextNodekey] = true
			}
		}
	}

	return nodeThatEnded
}
