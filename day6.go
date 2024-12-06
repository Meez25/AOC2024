package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"slices"
)

func daySix() {
	// Visited map[]Points
	input, _ := os.ReadFile("day6input.txt")
	lines := bytes.Split(input, []byte("\n"))

	//remove last one
	lines = lines[:len(lines)-1]
	var guard Guard

	// Place guard
	for y, line := range lines {
		for x, digit := range line {
			if string(digit) != "." && string(digit) != "#" {
				switch string(digit) {
				case ">":
					guard.direction = "right"
				case "^":
					guard.direction = "up"
				case "<":
					guard.direction = "left"
				case "v":
					guard.direction = "down"
				}
				guard.positionX = x
				guard.positionY = y
				guard.initialX = x
				guard.initialY = y
				guard.initialDirection = guard.direction
			}
			if string(digit) == "#" {
				guard.obstacle = append(guard.obstacle, formatPoints(x, y))
				guard.initialObstacle = guard.obstacle
			}
		}
	}

	// Show guard info

	// Move guard
	for guard.positionY > 0 && guard.positionY < len(lines) && guard.positionX > 0 && guard.positionX < len(lines[0]) {
		err := guard.move()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(len(guard.visitedPoints) - 1)
	visitedPoints := guard.visitedPoints

	// PART 2
	// Now that I have all the visitedPoints, I can try to put an obstacle after each points, and check if a loop happens

	timesInLoop := 0
Outer:
	for i := range visitedPoints {
		fmt.Println(i)
		guard.resetSimulation()
		if visitedPoints[i] == formatPoints(guard.initialX, guard.positionY) {
			continue Outer
		}
		if i < len(visitedPoints)-1 {
			guard.obstacle = append(guard.obstacle, visitedPoints[i])
			for guard.positionY > 0 && guard.positionY < len(lines) && guard.positionX > 0 && guard.positionX < len(lines[0]) {
				err := guard.move()
				if err != nil {
					fmt.Println(err)
					timesInLoop++
					break
				}
			}
			continue Outer
		}
	}

	fmt.Println(timesInLoop)

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
}

func (g *Guard) resetSimulation() {
	g.direction = g.initialDirection
	g.positionX = g.initialX
	g.positionY = g.initialY
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
		// fmt.Println("I already saw", g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-1], "at position", firstTime)
		// fmt.Println(g.iVeBeenHereHmmm[firstTime-1], firstTime, g.iVeBeenHereHmmm[firstTime-1], g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-2], firstTime != len(g.iVeBeenHereHmmm)-1)
		if firstTime != -1 && firstTime != 0 && g.iVeBeenHereHmmm[firstTime-1] == g.iVeBeenHereHmmm[len(g.iVeBeenHereHmmm)-2] && firstTime != len(g.iVeBeenHereHmmm)-1 {
			return true
		}
	}
	return false
}

func (g *Guard) canMoveFront() (bool, int, int) {
	var futurePosition string
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

	futurePosition = formatPoints(futureX, futureY)

	if slices.Contains(g.obstacle, futurePosition) {
		return false, 0, 0
	}
	return true, futureX, futureY
}
