package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type Grid struct {
	height int
	width  int
	robots []*Robot
}

func (g *Grid) Move() {
	for _, robot := range g.robots {
		robot.x += robot.speedX
		if robot.x >= g.width {
			robot.x -= g.width
		}
		if robot.x < 0 {
			robot.x += g.width
		}
		robot.y += robot.speedY
		if robot.y >= g.height {
			robot.y -= g.height
		}
		if robot.y < 0 {
			robot.y += g.height
		}
	}
}

func (g *Grid) Count() (int, int, int, int) {
	countA := 0
	countB := 0
	countC := 0
	countD := 0
	for i := 0; i < g.height/2; i++ {
		for j := 0; j < g.width/2; j++ {
			for _, robot := range g.robots {
				if robot.x == j && robot.y == i {
					countA++
				}
			}
		}
	}

	for i := g.height/2 + 1; i < g.height; i++ {
		for j := 0; j < g.width/2; j++ {
			for _, robot := range g.robots {
				if robot.x == j && robot.y == i {
					countB++
				}
			}
		}
	}

	for i := g.height/2 + 1; i < g.height; i++ {
		for j := g.width/2 + 1; j < g.width; j++ {
			for _, robot := range g.robots {
				if robot.x == j && robot.y == i {
					countD++
				}
			}
		}
	}

	for i := 0; i < g.height/2; i++ {
		for j := g.width/2 + 1; j < g.width; j++ {
			for _, robot := range g.robots {
				if robot.x == j && robot.y == i {
					countC++
				}
			}
		}
	}
	return countA, countC, countB, countD
}

func (g *Grid) IsChristmasTree() bool {
	emptyLines := 0
	emptyCol := 0
	for i := 0; i < g.height; i++ {
		robotInLine := 0
		for _, robot := range g.robots {
			if robot.y == i {
				robotInLine++
			}
		}
		if robotInLine == 0 {
			emptyLines++
		}
	}

	for i := 0; i < g.width; i++ {
		robotInLine := 0
		for _, robot := range g.robots {
			if robot.x == i {
				robotInLine++
			}
		}
		if robotInLine == 0 {
			emptyCol++
		}
	}
	if emptyLines > 10 && emptyCol > 10 {
		return true
	}
	return false
}

func (g *Grid) Describe() {
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			count := 0
			for _, robot := range g.robots {
				if robot.x == j && robot.y == i {
					count++
				}
			}
			if count > 0 {
				fmt.Print(count)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func newGrid() *Grid {
	var robotList []*Robot
	return &Grid{height: 103, width: 101, robots: robotList}
}

type Robot struct {
	x      int
	y      int
	speedX int
	speedY int
}

func newRobot(x, y, speedX, speedY int) *Robot {
	return &Robot{x: x, y: y, speedX: speedX, speedY: speedY}
}

func dayFourteen() {
	input, _ := os.ReadFile("day14input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n"))
	exp := regexp.MustCompile(`(-?\d+)`)

	var robots []*Robot

	grid := newGrid()
	for _, line := range lines {
		result := exp.FindAllStringSubmatch(string(line), -1)
		x, _ := strconv.Atoi(result[0][0])
		y, _ := strconv.Atoi(result[1][0])
		speedX, _ := strconv.Atoi(result[2][0])
		speedY, _ := strconv.Atoi(result[3][0])
		robots = append(robots, newRobot(x, y, speedX, speedY))
	}
	grid.robots = robots

	i := 0
	for i < 12000 {
		grid.Move()
		if grid.IsChristmasTree() {
			fmt.Println("Christmas tree pattern found after", i+1, "iterations:")
			grid.Describe()
			break
		}
		i++
	}
}
