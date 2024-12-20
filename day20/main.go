package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"time"
)

type Grid [141][141]byte

type Position struct {
	x, y int
}

type Shortcut struct {
	from     Position
	to       Position
	distance int
}

type ReachablePosition struct {
	pos   Position
	steps int
}

type Queue struct {
	nodes    [500]Node
	head     int
	tail     int
	size     int
	capacity int
}

func NewQueue() *Queue {
	return &Queue{
		nodes:    [500]Node{},
		head:     0,
		tail:     0,
		size:     0,
		capacity: 500,
	}
}

func (q *Queue) Push(n Node) bool {
	if q.size >= q.capacity {
		return false // Queue is full
	}

	q.nodes[q.tail] = n
	q.tail = (q.tail + 1) % q.capacity
	q.size++
	return true
}

func (q *Queue) Pop() (Node, bool) {
	if q.size == 0 {
		return Node{}, false // Queue is empty
	}

	node := q.nodes[q.head]
	q.head = (q.head + 1) % q.capacity
	q.size--
	return node, true
}

func (q *Queue) IsEmpty() bool {
	return q.size == 0
}

type Node struct {
	step     int
	position Position
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

	p1 := solve(grid, 2)
	partOne := time.Since(start)
	fmt.Println("Part 1:", p1, "in", partOne)

	startp2 := time.Now()
	p2 := solve(grid, 20)
	partTwo := time.Since(startp2)
	fmt.Println("Part 2:", p2, "in", partTwo)

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

func solve(grid Grid, cheat int) int {
	baseLinePath := findNoCheating(grid)
	shortcuts := FindShortcuts(grid, baseLinePath, cheat)

	count := 0
	timeSaved := 100
	for _, shortcut := range shortcuts {
		if shortcut.distance >= timeSaved {
			count++
		}
	}

	return count
}

func findNoCheating(grid Grid) []Position {
	startPos, endPos := grid.GetStartAndEndPosition()
	x, y := startPos.x, startPos.y
	directions := [4][2]int{
		{0, 1},  // bottom
		{0, -1}, // up
		{1, 0},  // right
		{-1, 0}, // left
	}
	visited := make(map[Position]bool, 1000)
	path := []Position{startPos}
	visited[startPos] = true

	for x != endPos.x || y != endPos.y {
		moved := false
		for _, direction := range directions {
			newX, newY := x+direction[0], y+direction[1]
			newPos := Position{x: newX, y: newY}

			if newX >= 0 && newX < len(grid[0]) &&
				newY >= 0 && newY < len(grid) &&
				grid[newY][newX] != '#' &&
				!visited[newPos] {
				x, y = newX, newY
				path = append(path, newPos)
				visited[newPos] = true
				moved = true
				break
			}
		}
		if !moved {
			break
		}
	}
	return path
}

func FindShortcuts(grid Grid, positions []Position, cheatNo int) []Shortcut {
	var shortcuts []Shortcut
	for i, pos := range positions {
		reachable := miniBFS(grid, pos, cheatNo)
		for _, reach := range reachable {
			// Look for this position starting after where we would end up after the shortcut
			startLookingAt := i + reach.steps
			if startLookingAt < len(positions) {
				for j := startLookingAt; j < len(positions); j++ {
					if positions[j] == reach.pos {
						// We saved: distance in original path - steps taken in shortcut
						savedSteps := j - i - reach.steps
						if savedSteps > 0 {
							shortcuts = append(shortcuts, Shortcut{
								from:     pos,
								to:       reach.pos,
								distance: savedSteps,
							})
						}
						break
					}
				}
			}
		}
	}
	return shortcuts
}

func miniBFS(grid Grid, position Position, cheatNo int) []ReachablePosition {
	visited := make(map[Position]bool, 5000)
	queue := NewQueue()
	reachablePositions := make([]ReachablePosition, 0, 100)
	directions := [4][2]int{
		{0, 1},  // bottom
		{0, -1}, // up
		{1, 0},  // right
		{-1, 0}, // left
	}

	visited[position] = true
	queue.Push(Node{position: position, step: 0})

	for !queue.IsEmpty() {
		v, _ := queue.Pop()

		if v.step >= cheatNo {
			continue
		}

		for _, direction := range directions {
			newX := v.position.x + direction[0]
			newY := v.position.y + direction[1]
			newPos := Position{x: newX, y: newY}

			if newX >= 0 && newX < len(grid[0]) &&
				newY >= 0 && newY < len(grid) &&
				!visited[newPos] {

				visited[newPos] = true
				queue.Push(Node{
					position: newPos,
					step:     v.step + 1,
				})

				if grid[newY][newX] != '#' {
					reachablePositions = append(reachablePositions, ReachablePosition{
						pos:   newPos,
						steps: v.step + 1,
					})
				}
			}
		}
	}
	return reachablePositions
}
