package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
	"time"
)

func main() {
	start := time.Now()
	input, _ := os.ReadFile("day6input.txt")
	lines := bytes.Split(input, []byte("\n"))

	//remove last one
	lines = lines[:len(lines)-1]
	var guard Guard

	guard.obstacle = make([]string, 0, 100)
	guard.iVeBeenHereHmmm = make([]string, 0, 1000)
	guard.visitedPoints = make([]string, 0, 1000)

	// Place guard
	for y, line := range lines {
		for x, digit := range line {
			if digit != '.' && digit != '#' {
				switch digit {
				case '>':
					guard.direction = "right"
				case '^':
					guard.direction = "up"
				case '<':
					guard.direction = "left"
				case 'v':
					guard.direction = "down"
				}
				guard.positionX = x
				guard.positionY = y
				guard.initialX = x
				guard.initialY = y
				guard.initialDirection = guard.direction
				guard.states = make(map[int]State, 100)
			}
			if digit == '#' {
				guard.obstacle = append(guard.obstacle, formatPoints(x, y))
				guard.initialObstacle = guard.obstacle
			}
		}
	}

	guard.obstacleMap = make(map[string]bool, len(guard.obstacle))
	for _, obs := range guard.obstacle {
		guard.obstacleMap[obs] = true
	}

	// Show guard info
	guard.states[0] = State{direction: guard.direction, positionX: guard.positionX, positionY: guard.positionY}

	// Move guard
	for guard.positionY > 0 && guard.positionY < len(lines) && guard.positionX > 0 && guard.positionX < len(lines[0]) {
		// Save each position
		guard.states[len(guard.visitedPoints)] = State{direction: guard.direction, positionX: guard.positionX, positionY: guard.positionY}
		err := guard.move()
		if err != nil {
			fmt.Println(err)
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", len(guard.visitedPoints)-1, "in", elapsed)
	start = time.Now()

	visitedPoints := guard.visitedPoints

	// PART 2
	// Now that I have all the visitedPoints, I can try to put an obstacle after each points, and check if a loop happens

	timesInLoop := 0
	initialPoint := formatPoints(guard.initialX, guard.initialY)

	chunkSize := 100
	for i := 0; i < len(visitedPoints); i += chunkSize {
		end := i + chunkSize
		if end > len(visitedPoints) {
			end = len(visitedPoints)
		}
		fmt.Printf("Processing points %d-%d of %d\n", i, end-1, len(visitedPoints))

		for j := i; j < end; j++ {
			if visitedPoints[j] == initialPoint {
				continue
			}

			if j > 0 {
				guard.loadState(j)
			}

			guard.obstacleMap[visitedPoints[j]] = true
			loopFound := false

			for guard.positionY > 0 && guard.positionY < len(lines) && guard.positionX > 0 && guard.positionX < len(lines[0]) {
				err := guard.move()
				if err != nil {
					timesInLoop++
					loopFound = true
					break
				}
			}

			guard.obstacleMap[visitedPoints[j]] = false
			if !loopFound {
				guard.resetSimulation()
			}
		}
	}
	elapsed = time.Since(start)
	fmt.Println("Part 2 :", timesInLoop, "in", elapsed)

}

func formatPoints(x, y int) string {
	formattedX := fmt.Sprintf("%04d", x)
	formattedY := fmt.Sprintf("%04d", y)
	concat := string(formattedX) + string(formattedY)
	return concat
}

type Guard struct {
	direction        string
	visitedPoints    []string
	positionX        int
	positionY        int
	obstacle         []string
	iVeBeenHereHmmm  []string
	initialX         int
	initialY         int
	initialDirection string
	initialObstacle  []string
	states           map[int]State
	obstacleMap      map[string]bool
}

type State struct {
	direction string
	positionX int
	positionY int
}

func (g *Guard) resetSimulation() {
	g.direction = g.initialDirection
	g.positionX = g.initialX
	g.positionY = g.initialY
	g.obstacle = g.initialObstacle
	g.iVeBeenHereHmmm = nil
	g.visitedPoints = nil
}

func (g *Guard) loadState(i int) {
	state := g.states[i-1]
	g.direction = state.direction
	g.positionX = state.positionX
	g.positionY = state.positionY
	g.obstacle = g.initialObstacle
	g.iVeBeenHereHmmm = nil
	g.visitedPoints = nil
}

func (g *Guard) move() error {
	canMove, futureX, futureY := g.canMoveFront()
	if canMove {
		g.positionX = futureX
		g.positionY = futureY
		g.MarkPositionVisited()
	} else {
		switch g.direction {
		case "up":
			g.direction = "right"
		case "right":
			g.direction = "down"
		case "down":
			g.direction = "left"
		case "left":
			g.direction = "up"
		default:
			fmt.Println("J'ai pas de direction !")
		}
		g.move()
	}
	loop := g.amILooping()
	if loop {
		return errors.New("I'm looping !")
	}
	return nil
}

func (g *Guard) MarkPositionVisited() {
	if !slices.Contains(g.visitedPoints, formatPoints(g.positionX, g.positionY)) {
		g.visitedPoints = append(g.visitedPoints, formatPoints(g.positionX, g.positionY))
	}
	g.iVeBeenHereHmmm = append(g.iVeBeenHereHmmm, formatPoints(g.positionX, g.positionY))
}

func (g *Guard) amILooping() bool {
	// If the last 2 movement have already happened, I'm in a loop
	if len(g.iVeBeenHereHmmm) > 1 {
		firstTime := slices.Index(g.iVeBeenHereHmmm, g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-1])
		if firstTime == 0 {
			firstTime = slices.Index(g.iVeBeenHereHmmm[1:], g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-1])
		}
		if firstTime != -1 && firstTime != 0 && g.iVeBeenHereHmmm[firstTime-1] == g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-2] && firstTime != len(g.iVeBeenHereHmmm)-1 {
			return true
		}
	}
	return false
}

func (g *Guard) canMoveFront() (bool, int, int) {
	var futureX = g.positionX
	var futureY = g.positionY
	switch g.direction {
	case "up":
		futureY--
	case "down":
		futureY++
	case "right":
		futureX++
	case "left":
		futureX--
	}

	futurePosition := formatPoints(futureX, futureY)
	return !g.obstacleMap[futurePosition], futureX, futureY
}
