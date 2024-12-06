package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
)

func daySix() {
	// Visited map[]Points
	formatted := formatPoints(12, 67)
	fmt.Println(formatted)

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
			}
			if string(digit) == "#" {
				guard.obstacle = append(guard.obstacle, formatPoints(x, y))
			}
		}
	}

	// Show guard info
	fmt.Println(guard)
	for guard.positionY > 0 && guard.positionY < len(lines) && guard.positionX > 0 && guard.positionX < len(lines[0]) {
		err := guard.move()
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(guard)
	fmt.Println(len(guard.visitedPoints) - 1)
}

func formatPoints(x, y int) string {
	formattedX := fmt.Sprintf("%04d", x)
	formattedY := fmt.Sprintf("%04d", y)
	concat := string(formattedX) + string(formattedY)
	return concat
}

type Guard struct {
	direction     string
	visitedPoints []string
	positionX     int
	positionY     int
	obstacle      []string
}

func (g *Guard) move() error {
	canMove, futureX, futureY := g.canMoveFront()
	if canMove {
		g.positionX = futureX
		g.positionY = futureY
		g.MarkPositionVisited()
		return nil
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
		}
		g.move()
		// return errors.New("Can't move !")
	}
	return nil
}

func (g *Guard) MarkPositionVisited() {
	if !slices.Contains(g.visitedPoints, formatPoints(g.positionX, g.positionY)) {
		g.visitedPoints = append(g.visitedPoints, formatPoints(g.positionX, g.positionY))
	}
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
