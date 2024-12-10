package main

import (
	"bytes"
	"fmt"
	"os"
)

type position struct {
	x int
	y int
}

func dayTen() {
	input, _ := os.ReadFile("day10input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))

	sum := 0
	sumP2 := 0

	for y, line := range lines {
		for x, char := range line {
			if char == '0' {
				count2 := BFSStep2(lines, position{x: x, y: y})
				count := BFS(lines, position{x: x, y: y})
				sum = sum + count
				sumP2 = sumP2 + count2
			}
		}
	}

	fmt.Println(sum)
	fmt.Println(sumP2)

}

func BFS(grid [][]byte, node position) int {
	marked := make(map[position]bool, 0)
	queue := make([]position, 0)
	queue = append(queue, node)

	foundNine := 0

	marked[node] = true

	for len(queue) > 0 {
		v := queue[0]
		if grid[v.y][v.x] == '9' {
			foundNine++
		}
		queue = queue[1:]
		// up
		var next position
		if v.y > 0 && int(grid[v.y-1][v.x])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x, y: v.y - 1}
			if marked[next] == false {
				marked[next] = true
				queue = append(queue, next)
			}
		}
		// down
		if v.y < len(grid)-1 && int(grid[v.y+1][v.x])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x, y: v.y + 1}
			if marked[next] == false {
				marked[next] = true
				queue = append(queue, next)
			}
		}
		// left
		if v.x > 0 && int(grid[v.y][v.x-1])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x - 1, y: v.y}
			if marked[next] == false {
				marked[next] = true
				queue = append(queue, next)
			}
		}
		// right
		if v.x < len(grid[0])-1 && int(grid[v.y][v.x+1])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x + 1, y: v.y}
			if marked[next] == false {
				marked[next] = true
				queue = append(queue, next)
			}
		}
	}
	return foundNine
}

func BFSStep2(grid [][]byte, node position) int {
	queue := make([]position, 0)
	queue = append(queue, node)

	foundNine := 0

	for len(queue) > 0 {
		v := queue[0]
		if grid[v.y][v.x] == '9' {
			foundNine++
		}
		queue = queue[1:]
		// up
		var next position
		if v.y > 0 && int(grid[v.y-1][v.x])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x, y: v.y - 1}
			queue = append(queue, next)
		}
		// down
		if v.y < len(grid)-1 && int(grid[v.y+1][v.x])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x, y: v.y + 1}
			queue = append(queue, next)
		}
		// left
		if v.x > 0 && int(grid[v.y][v.x-1])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x - 1, y: v.y}
			queue = append(queue, next)
		}
		// right
		if v.x < len(grid[0])-1 && int(grid[v.y][v.x+1])-int(grid[v.y][v.x]) == 1 {
			next = position{x: v.x + 1, y: v.y}
			queue = append(queue, next)
		}
	}
	return foundNine
}
