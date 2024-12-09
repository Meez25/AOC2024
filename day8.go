package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

type Point struct {
	x int
	y int
}

func dayEight() {
	// start := time.Now()
	input, _ := os.ReadFile("day8input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))

	nodes := make(map[string][]Point, 0)
	antiNodes := make([]Point, 0)
	fmt.Println(antiNodes)

	maxXGrid := len(lines[0])
	maxYGrid := len(lines)

	for y, line := range lines {
		for x, char := range line {
			if string(char) != "." {
				nodes[string(char)] = append(nodes[string(char)], Point{x: x, y: y})
			}
		}
	}

	// Compute antinodes
	for k, v := range nodes {
		for i := 0; i < len(v); i++ {
			fmt.Println("main", k, v[i])
			for j := 0; j < len(v); j++ {
				if i == j {
					continue
				}
				fmt.Println("compared with", v[j])
				antiNodeX := v[i].x + (v[i].x - v[j].x)
				antiNodeY := v[i].y + (v[i].y - v[j].y)
				fmt.Println(antiNodeX, antiNodeY)
				if antiNodeX < maxXGrid && antiNodeX >= 0 && antiNodeY >= 0 && antiNodeY < maxYGrid {
					newAntiNode := Point{x: antiNodeX, y: antiNodeY}
					if !slices.Contains(antiNodes, newAntiNode) {
						antiNodes = append(antiNodes, newAntiNode)
					}
				}
			}
		}
	}

	fmt.Println(len(antiNodes))

}
