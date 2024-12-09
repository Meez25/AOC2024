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
				antiNodes = append(antiNodes, Point{x: x, y: y})
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

				differenceX := v[i].x - v[j].x
				differenceY := v[i].y - v[j].y

				computedX := antiNodeX
				computedY := antiNodeY

				for computedX < maxXGrid && computedX >= 0 && computedY < maxYGrid && computedY >= 0 {
					if antiNodeX < maxXGrid && antiNodeX >= 0 && antiNodeY >= 0 && antiNodeY < maxYGrid {
						newAntiNode := Point{x: computedX, y: computedY}
						if !slices.Contains(antiNodes, newAntiNode) {
							antiNodes = append(antiNodes, newAntiNode)
						}
						computedX += differenceX
						computedY += differenceY
					}
				}
			}
		}
	}

	fmt.Println(antiNodes)
	fmt.Println(len(antiNodes))

}
