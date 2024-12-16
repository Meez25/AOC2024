package day16

import (
	"bytes"
	"container/heap"
	"fmt"
	"os"
)

type direction int

const (
	NORTH direction = iota
	EAST
	SOUTH
	WEST
)

type Point struct {
	X         int
	Y         int
	Direction direction
	Cost      int
}

type CostMap struct {
	costs map[string]int
}

func NewCostMap() *CostMap {
	return &CostMap{
		costs: make(map[string]int),
	}
}

func (cm *CostMap) Set(p Point) {
	key := fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Direction)
	if existing, exists := cm.costs[key]; !exists || p.Cost < existing {
		cm.costs[key] = p.Cost
	}
}

func (cm *CostMap) Get(p Point) (int, bool) {
	key := fmt.Sprintf("%d,%d,%d", p.X, p.Y, p.Direction)
	cost, exists := cm.costs[key]
	return cost, exists
}

type PriorityQueue []Point

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Cost < pq[j].Cost
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}
func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(Point))
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[0 : n-1]
	return item
}

func turnLeft(d direction) direction {
	switch d {
	case NORTH:
		return WEST
	case WEST:
		return SOUTH
	case SOUTH:
		return EAST
	case EAST:
		return NORTH
	}
	return d
}

func turnRight(d direction) direction {
	switch d {
	case NORTH:
		return EAST
	case EAST:
		return SOUTH
	case SOUTH:
		return WEST
	case WEST:
		return NORTH
	}
	return d
}

func reverseDirection(d direction) direction {
	switch d {
	case NORTH:
		return SOUTH
	case SOUTH:
		return NORTH
	case EAST:
		return WEST
	case WEST:
		return EAST
	}
	return d
}

func getNextPosition(p Point) (int, int) {
	switch p.Direction {
	case NORTH:
		return p.X, p.Y - 1
	case SOUTH:
		return p.X, p.Y + 1
	case EAST:
		return p.X + 1, p.Y
	case WEST:
		return p.X - 1, p.Y
	}
	return p.X, p.Y
}

func isValidPosition(x, y int, grid [][]string) bool {
	return x >= 0 && x < len(grid[0]) && y >= 0 && y < len(grid) && grid[y][x] != "#"
}

func getNeighbors(p Point, grid [][]string) []Point {
	neighbors := make([]Point, 0)

	// Forward
	newX, newY := getNextPosition(p)
	if isValidPosition(newX, newY, grid) {
		neighbors = append(neighbors, Point{
			X:         newX,
			Y:         newY,
			Direction: p.Direction,
			Cost:      p.Cost + 1,
		})
	}

	// Left
	neighbors = append(neighbors, Point{
		X:         p.X,
		Y:         p.Y,
		Direction: turnLeft(p.Direction),
		Cost:      p.Cost + 1000,
	})

	// Right
	neighbors = append(neighbors, Point{
		X:         p.X,
		Y:         p.Y,
		Direction: turnRight(p.Direction),
		Cost:      p.Cost + 1000,
	})

	return neighbors
}

func findStartEnd(grid [][]string) (int, int, int, int) {
	var startX, startY, endX, endY int
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == "S" {
				startX, startY = x, y
			}
			if grid[y][x] == "E" {
				endX, endY = x, y
			}
		}
	}
	return startX, startY, endX, endY
}

func exploreFromPoint(start Point, grid [][]string) *CostMap {
	costMap := NewCostMap()
	visited := make(map[string]bool)
	pq := &PriorityQueue{}
	heap.Init(pq)
	heap.Push(pq, start)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Point)
		key := fmt.Sprintf("%d,%d,%d", current.X, current.Y, current.Direction)

		if visited[key] {
			continue
		}

		visited[key] = true
		costMap.Set(current)

		for _, neighbor := range getNeighbors(current, grid) {
			neighborKey := fmt.Sprintf("%d,%d,%d", neighbor.X, neighbor.Y, neighbor.Direction)
			if !visited[neighborKey] {
				heap.Push(pq, neighbor)
			}
		}
	}

	return costMap
}

func findOptimalPaths(grid [][]string) (int, int) {
	startX, startY, endX, endY := findStartEnd(grid)

	// Forward exploration (start to end)
	startPoint := Point{X: startX, Y: startY, Direction: EAST, Cost: 0}
	forwardCosts := exploreFromPoint(startPoint, grid)

	// Backward exploration (end to start)
	// Try all possible end directions
	minCost := -1
	var backwardCosts *CostMap
	for _, dir := range []direction{NORTH, EAST, SOUTH, WEST} {
		endPoint := Point{X: endX, Y: endY, Direction: dir, Cost: 0}
		tempCosts := exploreFromPoint(endPoint, grid)

		// Check if this direction gives us a better total cost
		for y := range grid {
			for x := range grid[y] {
				if grid[y][x] == "#" {
					continue
				}
				for _, fwdDir := range []direction{NORTH, EAST, SOUTH, WEST} {
					fwdPoint := Point{X: x, Y: y, Direction: fwdDir}
					bwdPoint := Point{X: x, Y: y, Direction: reverseDirection(fwdDir)}

					if fwdCost, fwdExists := forwardCosts.Get(fwdPoint); fwdExists {
						if bwdCost, bwdExists := tempCosts.Get(bwdPoint); bwdExists {
							totalCost := fwdCost + bwdCost
							if minCost == -1 || totalCost < minCost {
								minCost = totalCost
								backwardCosts = tempCosts
							}
						}
					}
				}
			}
		}
	}

	// Find all positions that are part of optimal paths
	optimalPositions := make(map[string]bool)
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == "#" {
				continue
			}

			for _, dir := range []direction{NORTH, EAST, SOUTH, WEST} {
				fwdPoint := Point{X: x, Y: y, Direction: dir}
				bwdPoint := Point{X: x, Y: y, Direction: reverseDirection(dir)}

				if fwdCost, fwdExists := forwardCosts.Get(fwdPoint); fwdExists {
					if bwdCost, bwdExists := backwardCosts.Get(bwdPoint); bwdExists {
						if fwdCost+bwdCost == minCost {
							posKey := fmt.Sprintf("%d,%d", x, y)
							optimalPositions[posKey] = true
						}
					}
				}
			}
		}
	}

	return minCost, len(optimalPositions)
}

func Day() {
	input, _ := os.ReadFile("day16/day16input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))

	grid := make([][]string, len(lines))
	for i := range grid {
		grid[i] = make([]string, len(lines[0]))
	}

	for y, line := range lines {
		for x, char := range line {
			grid[y][x] = string(char)
		}
	}

	minCost, optimalTiles := findOptimalPaths(grid)
	fmt.Printf("Part 1 - Lowest possible score: %d\n", minCost)
	fmt.Printf("Part 2 - Number of tiles on optimal paths: %d\n", optimalTiles)
}
