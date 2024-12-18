package main

import (
	_ "embed"
	"fmt"
	"strings"
	"time"
)

//go:embed day18input.txt
var inputFile string

type Node struct {
	x    int
	y    int
	path []Position
}

type Position struct {
	x int
	y int
}

func newNode(x, y int) *Node {
	path := make([]Position, 0)
	return &Node{x: x, y: y, path: path}
}

func (n *Node) AddPosition(position Position) {
	n.path = append(n.path, position)
}

func main() {
	start := time.Now()
	grid := make([][]string, 71)
	for i := range grid {
		grid[i] = make([]string, 71)
		for j := range grid[i] {
			grid[i][j] = "."
		}
	}

	for _, coord := range strings.Split(strings.TrimSpace(inputFile), "\n")[:1024] {
		parts := strings.Split(coord, ",")
		x, y := 0, 0
		fmt.Sscanf(parts[0]+","+parts[1], "%d,%d", &x, &y)
		grid[y][x] = "#"
	}

	// Printing grid
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}

	// BFS to find exit
	successNodes := BFS(grid)
	path := successNodes[0].path

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", len(successNodes[0].path), "in", elapsed)

	// Part 2, doing BFS for each new corrupted memory
	firstPosition := Part2(inputFile, path, grid)

	elapsed = time.Since(start)
	fmt.Println("Part 2 :", firstPosition, elapsed)

	// Printing end grid
	for y := range grid {
		for x := range grid[y] {
			for _, position := range path {
				if y == position.y && x == position.x {
					grid[position.y][position.x] = "O"
				}
			}
		}
	}
	for _, row := range grid {
		fmt.Println(strings.Join(row, ""))
	}
}

func Part2(inputFile string, path []Position, grid [][]string) Position {
	pathAsPosition := make([]Position, len(path))

	var firstPositionToMatch Position
	var allInputFalling []Position

	for i, pos := range path {
		position := Position{x: pos.x, y: pos.y}
		pathAsPosition[i] = position
	}

	for _, coord := range strings.Split(strings.TrimSpace(inputFile), "\n") {
		parts := strings.Split(coord, ",")
		x, y := 0, 0
		fmt.Sscanf(parts[0]+","+parts[1], "%d,%d", &x, &y)
		position := Position{x: x, y: y}
		allInputFalling = append(allInputFalling, position)
	}

	i := 1024
	for i := i; i < len(allInputFalling); i++ {
		grid[allInputFalling[i].y][allInputFalling[i].x] = "#"
		resultAfterFalling := BFS(grid)
		if len(resultAfterFalling) == 0 {
			firstPositionToMatch = allInputFalling[i]
			break
		}
	}
	return firstPositionToMatch
}

func BFS(grid [][]string) []Node {
	startNode := newNode(0, 0)
	visited := make(map[Position]bool)
	queue := make([]*Node, 0)
	successNode := make([]Node, 0)

	queue = append(queue, startNode)
	visited[Position{x: 0, y: 0}] = true

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v.x == len(grid[0])-1 && v.y == len(grid)-1 {
			successNode = append(successNode, *v)
			continue
		}
		// Explore 4 sides
		// Up
		if v.y > 0 {
			up := grid[v.y-1][v.x]
			if up != "#" {
				upNode := newNode(v.x, v.y-1)
				upNode.path = make([]Position, len(v.path))
				copy(upNode.path, v.path)
				nextPosition := Position{x: upNode.x, y: upNode.y}
				upNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, upNode)
					visited[nextPosition] = true
				}
			}
		}
		// Down
		if v.y < len(grid)-1 {
			down := grid[v.y+1][v.x]
			if down != "#" {
				downNode := newNode(v.x, v.y+1)
				downNode.path = make([]Position, len(v.path))
				copy(downNode.path, v.path)
				nextPosition := Position{x: downNode.x, y: downNode.y}
				downNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, downNode)
					visited[nextPosition] = true
				}
			}
		}
		// Right
		if v.x < len(grid[0])-1 {
			right := grid[v.y][v.x+1]
			if right != "#" {
				rightNode := newNode(v.x+1, v.y)
				rightNode.path = make([]Position, len(v.path))
				copy(rightNode.path, v.path)
				nextPosition := Position{x: rightNode.x, y: rightNode.y}
				rightNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, rightNode)
					visited[nextPosition] = true
				}
			}
		}
		// Left
		if v.x > 0 {
			left := grid[v.y][v.x-1]
			if left != "#" {
				leftNode := newNode(v.x-1, v.y)
				leftNode.path = make([]Position, len(v.path))
				copy(leftNode.path, v.path)
				nextPosition := Position{x: leftNode.x, y: leftNode.y}
				leftNode.AddPosition(nextPosition)
				if !visited[nextPosition] {
					queue = append(queue, leftNode)
					visited[nextPosition] = true
				}
			}
		}
	}
	return successNode
}
