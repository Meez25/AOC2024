package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

type Region struct {
	name      string
	area      []Point
	perimeter int
}

func dayTwelve() {
	input, _ := os.ReadFile("day12input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))

	regions := make([]Region, 0)

	for y, line := range lines {
		for x, char := range line {
			region, perimeter := bfs(lines, Point{x: x, y: y})
			r := Region{name: string(char), area: region, perimeter: perimeter}
			contains := false
			for i := range regions {
				if regions[i].name == r.name {
					if slices.Contains(regions[i].area, r.area[0]) {
						contains = true
					}
				}

			}
			if !contains {
				regions = append(regions, r)
			}
		}
	}

	price := 0

	for _, region := range regions {
		price = price + (region.perimeter * len(region.area))
	}

	fmt.Println(price)
}

func bfs(grid [][]byte, point Point) ([]Point, int) {
	queue := make([]Point, 0)
	marked := make(map[Point]bool)
	queue = append(queue, point)

	perimeter := 0

	marked[point] = true

	region := make([]Point, 0)
	// AAAA
	// BBCD
	// BBCC
	// EEEC

	for len(queue) > 0 {
		var next Point
		p := queue[0]
		if grid[p.y][p.x] == grid[point.y][point.x] {
			region = append(region, p)
		}
		queue = queue[1:]
		// up
		if p.x > 0 {
			next = Point{x: p.x - 1, y: p.y}
			if grid[next.y][next.x] == grid[p.y][p.x] {
				if marked[next] == false {
					queue = append(queue, next)
					marked[next] = true
				}
			} else {
				perimeter++
			}
		} else {
			perimeter++
		}
		// down
		if p.x < len(grid[0])-1 {
			next = Point{x: p.x + 1, y: p.y}
			if grid[next.y][next.x] == grid[p.y][p.x] {
				if marked[next] == false {
					queue = append(queue, next)
					marked[next] = true
				}
			} else {
				perimeter++
			}
		} else {
			perimeter++
		}
		// left
		if p.y > 0 {
			next = Point{x: p.x, y: p.y - 1}
			if grid[next.y][next.x] == grid[p.y][p.x] {
				if marked[next] == false {
					queue = append(queue, next)
					marked[next] = true
				}
			} else {
				perimeter++
			}
		} else {
			perimeter++
		}
		// right
		if p.y < len(grid)-1 {
			next = Point{x: p.x, y: p.y + 1}
			if grid[next.y][next.x] == grid[p.y][p.x] {
				if marked[next] == false {
					queue = append(queue, next)
					marked[next] = true
				}
			} else {
				perimeter++
			}
		} else {
			perimeter++
		}
	}

	return region, perimeter
}
