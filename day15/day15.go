package day15

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"time"
)

type position struct {
	x int
	y int
}

type robot struct {
	position position
}

func newRobot(x, y int) *robot {
	return &robot{position{x: x, y: y}}
}

type boxes struct {
	position position
}

func (b boxes) getGPS() int {
	return 100*b.position.y + b.position.x
}

func newBox(x, y int) *boxes {
	return &boxes{position{x: x, y: y}}
}

type Grid struct {
	width  int
	height int
	robot  *robot
	boxes  []*boxes
	walls  []position
}

type Direction int

func (d Direction) toString() string {
	if d == 0 {
		return "UP"
	}
	if d == 1 {
		return "RIGHT"
	}
	if d == 2 {
		return "DOWN"
	}
	if d == 3 {
		return "LEFT"
	}
	return ""
}

const (
	UP = iota
	RIGHT
	DOWN
	LEFT
)

func newGrid() Grid {
	boxes := make([]*boxes, 0)
	robot := &robot{}
	walls := make([]position, 0)
	return Grid{boxes: boxes, robot: robot, walls: walls}
}

func Day() {
	input, _ := os.ReadFile("day15/day15input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n\n"))

	gridInput := lines[0]

	grid := newGrid()

	grid = parseGrid(gridInput, grid)
	instruction := parseInstruction(lines[1])
	grid.DoInstruction(instruction)
	fmt.Println(grid.Result())
}

func parseGrid(input []byte, grid Grid) Grid {
	lines := bytes.Split(input, []byte("\n"))
	for y, line := range lines {
		for x, digit := range line {
			if string(digit) == "O" {
				box := newBox(x, y)
				grid.boxes = append(grid.boxes, box)
			}
			if string(digit) == "@" {
				robot := newRobot(x, y)
				grid.robot = robot
			}
			if string(digit) == "#" {
				position := position{x, y}
				grid.walls = append(grid.walls, position)
			}
		}
	}
	grid.width = len(lines)
	grid.height = len(lines[0])
	return grid
}

func parseInstruction(input []byte) []Direction {
	direction := make([]Direction, 0)
	lines := bytes.Split(input, []byte("\n"))
	for _, line := range lines {
		for _, char := range line {
			if string(char) == "v" {
				direction = append(direction, DOWN)
			} else if string(char) == "<" {
				direction = append(direction, LEFT)
			} else if string(char) == ">" {
				direction = append(direction, RIGHT)
			} else if string(char) == "^" {
				direction = append(direction, UP)
			}
		}
	}
	return direction
}

func (g Grid) Result() int {
	sum := 0
	for _, box := range g.boxes {
		sum += box.getGPS()
	}
	return sum
}

func (g Grid) DoInstruction(input []Direction) {
	fmt.Print("\033[H\033[2J")
	for i, instruction := range input {
		fmt.Printf("\033[H") // Move cursor to top without clearing
		fmt.Printf("Move %d/%d: %s\n", i+1, len(input), instruction.toString())
		fmt.Println("----------------------------------------")
		if instruction == UP {
			if g.isLocationFree(g.robot.position.x, g.robot.position.y-1) {
				// There is no box up
				g.robot.position.y = g.robot.position.y - 1
			} else {
				// There is a box up !
				emptyPosition, err := g.nextFreePosition(g.robot.position.x, g.robot.position.y-1, UP)
				if err != nil {
					continue
				}
				for _, box := range g.boxes {
					if box.position.x == g.robot.position.x && box.position.y > emptyPosition.y && box.position.y <= g.robot.position.y {
						box.position.y -= 1
					}
				}
				g.robot.position.y -= 1
			}
		}
		if instruction == DOWN {
			if g.isLocationFree(g.robot.position.x, g.robot.position.y+1) {
				g.robot.position.y = g.robot.position.y + 1
			} else {
				emptyPosition, err := g.nextFreePosition(g.robot.position.x, g.robot.position.y+1, DOWN)
				if err != nil {
					continue
				}
				for _, box := range g.boxes {
					if box.position.x == g.robot.position.x && box.position.y < emptyPosition.y && box.position.y >= g.robot.position.y {
						box.position.y += 1
					}
				}
				g.robot.position.y += 1
			}
		}
		if instruction == RIGHT {
			if g.isLocationFree(g.robot.position.x+1, g.robot.position.y) {
				g.robot.position.x = g.robot.position.x + 1
			} else {
				emptyPosition, err := g.nextFreePosition(g.robot.position.x+1, g.robot.position.y, RIGHT)
				if err != nil {
					continue
				}
				for _, box := range g.boxes {
					if box.position.y == g.robot.position.y && box.position.x < emptyPosition.x && box.position.x >= g.robot.position.x {
						box.position.x += 1
					}
				}
				g.robot.position.x += 1
			}
		}
		if instruction == LEFT {
			if g.isLocationFree(g.robot.position.x-1, g.robot.position.y) {
				g.robot.position.x = g.robot.position.x - 1
			} else {
				emptyPosition, err := g.nextFreePosition(g.robot.position.x-1, g.robot.position.y, LEFT)
				if err != nil {
					continue
				}
				for _, box := range g.boxes {
					if box.position.y == g.robot.position.y && box.position.x > emptyPosition.x && box.position.x <= g.robot.position.x {
						box.position.x -= 1
					}
				}
				g.robot.position.x -= 1
			}
		}
		g.describe()

		// Show some stats
		fmt.Printf("\nRobot position: (%d,%d)\n", g.robot.position.x, g.robot.position.y)
		fmt.Println("----------------------------------------")

		// Small delay to make it visible
		time.Sleep(2 * time.Millisecond)

	}
}

func (g Grid) isLocationFree(x, y int) bool {
	for _, box := range g.boxes {
		if box.position.x == x && box.position.y == y {
			return false
		}
	}
	for _, wall := range g.walls {
		if wall.x == x && wall.y == y {
			return false
		}
	}
	return true
}

func (g Grid) describe() {
	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			// Start by assuming this position is empty
			char := "."

			// Check if it's a wall (assuming walls are at the borders)
			if x == 0 || x == g.width-1 || y == 0 || y == g.height-1 {
				char = "#"
			}

			// Check if there's a box at this position
			for _, box := range g.boxes {
				if box.position.x == x && box.position.y == y {
					char = "O"
					break
				}
			}

			for _, wall := range g.walls {
				if wall.x == x && wall.y == y {
					char = "X"
					break
				}
			}

			// Check if the robot is at this position
			if g.robot.position.x == x && g.robot.position.y == y {
				char = "@"
			}

			fmt.Print(char)
		}
		fmt.Print("\n")
	}
}

func (g Grid) nextFreePosition(x, y int, direction Direction) (position, error) {
	if direction == UP {
		for i := y; i > 0; i-- {
			if g.isLocationFree(x, i) {
				return position{x, i}, nil
			}
			for _, wall := range g.walls {
				if wall.x == x && wall.y == i {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	if direction == DOWN {
		for i := y; i < g.height; i++ {
			if g.isLocationFree(x, i) {
				return position{x, i}, nil
			}
			for _, wall := range g.walls {
				if wall.x == x && wall.y == i {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	if direction == RIGHT {
		for i := x; i < g.width; i++ {
			if g.isLocationFree(i, y) {
				return position{i, y}, nil
			}
			for _, wall := range g.walls {
				if wall.x == i && wall.y == y {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	if direction == LEFT {
		for i := x; i > 0; i-- {
			if g.isLocationFree(i, y) {
				return position{i, y}, nil
			}
			for _, wall := range g.walls {
				if wall.x == i && wall.y == y {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	return position{}, errors.New("The wall is reached")
}
