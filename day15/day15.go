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

type bigboxes struct {
	p1 position
	p2 position
}

type bigWalls struct {
	p1 position
	p2 position
}

func (b boxes) getGPS() int {
	return 100*b.position.y + b.position.x
}

func newBox(x, y int) *boxes {
	return &boxes{position{x: x, y: y}}
}

func (g Grid) ResultP2() int {
	sum := 0
	fmt.Println("Grid dimensions - height:", g.height, ", width:", g.width)
	for i, box := range g.bigBoxes {
		fmt.Println("Processing box", i+1, ":", box)

		// Calculate distances from edges
		leftDistance := box.p1.x
		rightDistance := g.width - box.p2.x - 1
		topDistance := box.p1.y

		fmt.Printf("Distances - Left: %d, Right: %d, Top: %d\n",
			leftDistance, rightDistance, topDistance)

		horizontal := leftDistance

		vertical := topDistance

		fmt.Printf("Chosen - Horizontal: %d, Vertical: %d\n", horizontal, vertical)

		// Calculate GPS coordinate
		gps := 100*vertical + horizontal
		fmt.Println("Calculated GPS:", gps)

		sum += gps
	}
	fmt.Println("Final GPS sum:", sum)
	return sum
}

func newBigBox(x, y int) *bigboxes {
	return &bigboxes{p1: position{x, y}, p2: position{x + 1, y}}
}

type Grid struct {
	width    int
	height   int
	robot    *robot
	boxes    []*boxes
	walls    []position
	bigBoxes []*bigboxes
	bigWalls []*bigWalls
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
	bigWalls := make([]*bigWalls, 0)
	bigBoxes := make([]*bigboxes, 0)
	return Grid{boxes: boxes, robot: robot, walls: walls, bigBoxes: bigBoxes, bigWalls: bigWalls}
}

func Day() {
	input, _ := os.ReadFile("day15/day15input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte("\n\n"))

	gridInput := lines[0]

	grid := newGrid()

	// grid = parseGrid(gridInput, grid)
	instruction := parseInstruction(lines[1])
	// grid.DoInstruction(instruction)
	// fmt.Println(grid.Result())
	grid = parseGridP2(gridInput, grid)
	grid.DoInstructionP2(instruction)
	fmt.Printf("Final GPS sum: %d\n", grid.ResultP2())
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

func parseGridP2(input []byte, grid Grid) Grid {
	lines := bytes.Split(input, []byte("\n"))
	for y, line := range lines {
		for x, digit := range line {
			if string(digit) == "O" {
				box := newBigBox(x*2, y)
				grid.bigBoxes = append(grid.bigBoxes, box)
			}
			if string(digit) == "@" {
				robot := newRobot(x*2, y)
				grid.robot = robot
			}
			if string(digit) == "#" {
				grid.bigWalls = append(grid.bigWalls, &bigWalls{p1: position{x * 2, y}, p2: position{x*2 + 1, y}})
			}
		}
	}
	grid.width = len(lines) * 2
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

func (g Grid) DoInstructionP2(input []Direction) {
	fmt.Print("\033[H\033[2J")
	// reader := bufio.NewReader(os.Stdin)
	for i, instruction := range input {
		fmt.Printf("\033[H") // Move cursor to top without clearing
		fmt.Printf("Move %d/%d: %s\n", i+1, len(input), instruction.toString())
		if i < len(input)-1 {
			fmt.Printf("Next move will be: %s\n", input[i+1].toString())
		} else {
			fmt.Printf("This is the last move\n")
		}
		fmt.Println("----------------------------------------")
		if instruction == UP {
			if g.isLocationFree(g.robot.position.x, g.robot.position.y-1) {
				// There is no box up
				g.robot.position.y = g.robot.position.y - 1
			} else {
				// There is a box up !
				for _, box := range g.bigBoxes {
					if (box.p1.x == g.robot.position.x && box.p1.y == g.robot.position.y-1) || (box.p2.x == g.robot.position.x && box.p2.y == g.robot.position.y-1) {
						canBoxMove, boxes := g.canBoxMove(*box, UP)
						if canBoxMove {
							for _, box := range boxes {
								box.p1.y -= 1
								box.p2.y -= 1
							}
							g.robot.position.y -= 1 // Don't forget to move the robot too!
						}
					}
				}
			}
		}
		if instruction == DOWN {
			if g.isLocationFree(g.robot.position.x, g.robot.position.y+1) {
				g.robot.position.y = g.robot.position.y + 1
			} else {
				// There is a box up !
				for _, box := range g.bigBoxes {
					if (box.p1.x == g.robot.position.x && box.p1.y == g.robot.position.y+1) || (box.p2.x == g.robot.position.x && box.p2.y == g.robot.position.y+1) {
						canBoxMove, boxes := g.canBoxMove(*box, DOWN)
						if canBoxMove {
							for _, box := range boxes {
								box.p1.y += 1
								box.p2.y += 1
							}
							g.robot.position.y += 1
						}
					}
				}
			}
		}
		if instruction == RIGHT {
			if g.isLocationFree(g.robot.position.x+1, g.robot.position.y) {
				g.robot.position.x = g.robot.position.x + 1
			} else {
				// Check for box to push
				for _, box := range g.bigBoxes {
					if (box.p1.x == g.robot.position.x+1 && box.p1.y == g.robot.position.y) ||
						(box.p2.x == g.robot.position.x+1 && box.p2.y == g.robot.position.y) {
						canBoxMove, boxes := g.canBoxMove(*box, RIGHT)
						if canBoxMove {
							for _, box := range boxes {
								box.p1.x += 1
								box.p2.x += 1
							}
							g.robot.position.x += 1
						}
					}
				}
			}
		}

		if instruction == LEFT {
			if g.isLocationFree(g.robot.position.x-1, g.robot.position.y) {
				g.robot.position.x = g.robot.position.x - 1
			} else {
				// Check for box to push
				for _, box := range g.bigBoxes {
					if (box.p1.x == g.robot.position.x-1 && box.p1.y == g.robot.position.y) ||
						(box.p2.x == g.robot.position.x-1 && box.p2.y == g.robot.position.y) {
						canBoxMove, boxes := g.canBoxMove(*box, LEFT)
						if canBoxMove {
							for _, box := range boxes {
								box.p1.x -= 1
								box.p2.x -= 1
							}
							g.robot.position.x -= 1
						}
					}
				}
			}
		}
		// g.describe()

		// Show some stats
		fmt.Printf("\nRobot position: (%d,%d)\n", g.robot.position.x, g.robot.position.y)
		fmt.Println("----------------------------------------")

		// Small delay to make it visible
		// fmt.Println("Press any key to continue...")
		// reader.ReadString('\n') // Wait for Enter key

	}
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
		time.Sleep(1000 * time.Millisecond)

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
	for _, bigWalls := range g.bigWalls {
		if bigWalls.p1.x == x && bigWalls.p1.y == y {
			return false
		}
		if bigWalls.p2.x == x && bigWalls.p2.y == y {
			return false
		}
	}

	for _, bigboxes := range g.bigBoxes {
		if bigboxes.p1.x == x && bigboxes.p1.y == y {
			return false
		}
		if bigboxes.p2.x == x && bigboxes.p2.y == y {
			return false
		}
	}
	return true
}

func (g Grid) IsWall(x, y int) bool {
	for _, bigWalls := range g.bigWalls {
		if bigWalls.p1.x == x && bigWalls.p1.y == y {
			return true
		}
		if bigWalls.p2.x == x && bigWalls.p2.y == y {
			return true
		}
	}
	return false
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

			if g.width > g.height {
				if x == 1 || x == g.width-2 {
					char = "#"
				}
			}

			// Check if there's a box at this position
			for _, box := range g.boxes {
				if box.position.x == x && box.position.y == y {
					char = "O"
					break
				}
			}

			for _, bigBoxes := range g.bigBoxes {
				if bigBoxes.p1.x == x && bigBoxes.p1.y == y {
					char = "["
					break
				}
				if bigBoxes.p2.x == x && bigBoxes.p2.y == y {
					char = "]"
					break
				}
			}

			for _, wall := range g.walls {
				if wall.x == x && wall.y == y {
					char = "X"
					break
				}
			}

			for _, bigWalls := range g.bigWalls {
				if bigWalls.p1.x == x && bigWalls.p1.y == y {
					char = "#"
					break
				}
				if bigWalls.p2.x == x && bigWalls.p2.y == y {
					char = "#"
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

func (g Grid) nextFreePositionP2(x, y int, direction Direction) (position, error) {
	if direction == UP {
		for i := y; i > 0; i-- {
			if g.isLocationFree(i, y) && g.isLocationFree(i+1, y) {
				return position{i, y}, nil
			}
			for _, wall := range g.walls {
				if wall.x == x && wall.y == i {
					return position{}, errors.New("The wall is reached")
				}
			}
			for _, bigWall := range g.bigWalls {
				if (bigWall.p1.x == i && bigWall.p1.y == y) ||
					(bigWall.p2.x == i && bigWall.p2.y == y) {
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
		// For rightward movement, we're looking for two consecutive free spaces
		// that we can push boxes into
		for i := x; i < g.width-1; i++ {
			if g.isLocationFree(i, y) && g.isLocationFree(i+1, y) {
				return position{i + 1, y}, nil
			}
			// Check for walls and return error if we hit one
			for _, bigWall := range g.bigWalls {
				if (bigWall.p1.x == i && bigWall.p1.y == y) ||
					(bigWall.p2.x == i && bigWall.p2.y == y) {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	if direction == LEFT {
		for i := x; i > 1; i-- { // Note: i > 1 because we need two spaces
			if g.isLocationFree(i, y) && g.isLocationFree(i-1, y) {
				return position{i - 1, y}, nil
			}
			// Check for walls and return error if we hit one
			for _, bigWall := range g.bigWalls {
				if (bigWall.p1.x == i && bigWall.p1.y == y) ||
					(bigWall.p2.x == i && bigWall.p2.y == y) {
					return position{}, errors.New("The wall is reached")
				}
			}
		}
	}
	return position{}, errors.New("The wall is reached")
}

func (g Grid) canBoxMove(big bigboxes, direction Direction) (bool, []*bigboxes) {
	boxCanMove := make([]*bigboxes, 0)
	marked := make(map[*bigboxes]bool)
	queue := make([]*bigboxes, 0)

	// Find the initial box pointer
	var startBox *bigboxes
	for _, b := range g.bigBoxes {
		if b.p1 == big.p1 && b.p2 == big.p2 {
			startBox = b
			break
		}
	}
	queue = append(queue, startBox)

	for len(queue) > 0 {
		boxToCheck := queue[0]
		queue = queue[1:]

		if !marked[boxToCheck] {
			if direction == UP {
				// Check if there's a box above that we haven't processed yet
				if g.IsWall(boxToCheck.p1.x, boxToCheck.p1.y-1) || g.IsWall(boxToCheck.p2.x, boxToCheck.p2.y-1) {
					return false, nil
				}
				var foundBox bool

				for _, box := range g.bigBoxes {
					if box.p1.y == boxToCheck.p1.y-1 &&
						(box.p1.x == boxToCheck.p1.x || // aligned
							box.p1.x == boxToCheck.p1.x-1 || // offset left
							box.p1.x == boxToCheck.p1.x+1) { // offset right
						queue = append(queue, box)
						foundBox = true
					}
				}
				// If no box above and space is blocked, can't move
				if !foundBox && (!g.isLocationFree(boxToCheck.p1.x, boxToCheck.p1.y-1) ||
					!g.isLocationFree(boxToCheck.p2.x, boxToCheck.p2.y-1)) {
					fmt.Println("I can't move!")
					return false, nil
				}

				// Find horizontally connected boxes at the same level
				for _, box := range g.bigBoxes {
					if !marked[box] && box.p1.y == boxToCheck.p1.y &&
						(box.p2.x == boxToCheck.p1.x || box.p1.x == boxToCheck.p2.x) {
						queue = append(queue, box)
					}
				}
			}
			if direction == DOWN {
				// Check if there's a box above that we haven't processed yet
				if g.IsWall(boxToCheck.p1.x, boxToCheck.p1.y+1) || g.IsWall(boxToCheck.p2.x, boxToCheck.p2.y+1) {
					return false, nil
				}
				var foundBox bool
				for _, box := range g.bigBoxes {
					if box.p1.y == boxToCheck.p1.y+1 &&
						(box.p1.x == boxToCheck.p1.x || // aligned
							box.p1.x == boxToCheck.p1.x-1 || // offset left
							box.p1.x == boxToCheck.p1.x+1) { // offset right
						queue = append(queue, box)
						foundBox = true
					}
				}
				// If no box above and space is blocked, can't move
				if !foundBox && (!g.isLocationFree(boxToCheck.p1.x, boxToCheck.p1.y+1) ||
					!g.isLocationFree(boxToCheck.p2.x, boxToCheck.p2.y+1)) {
					return false, nil
				}

				// Find horizontally connected boxes at the same level
				for _, box := range g.bigBoxes {
					if !marked[box] && box.p1.y == boxToCheck.p1.y &&
						(box.p2.x == boxToCheck.p1.x || box.p1.x == boxToCheck.p2.x) {
						queue = append(queue, box)
					}
				}
			}
			if direction == RIGHT {
				var foundBox bool
				for _, box := range g.bigBoxes {
					if box.p1.x == boxToCheck.p2.x+1 && box.p1.y == boxToCheck.p1.y {
						queue = append(queue, box)
						foundBox = true
					}
				}
				if !foundBox && !g.isLocationFree(boxToCheck.p2.x+1, boxToCheck.p1.y) {
					return false, nil
				}
			}
			if direction == LEFT {
				var foundBox bool
				for _, box := range g.bigBoxes {
					if box.p2.x == boxToCheck.p1.x-1 && box.p1.y == boxToCheck.p1.y {
						queue = append(queue, box)
						foundBox = true
					}
				}
				if !foundBox && !g.isLocationFree(boxToCheck.p1.x-1, boxToCheck.p1.y) {
					return false, nil
				}
			}

			marked[boxToCheck] = true
			boxCanMove = append(boxCanMove, boxToCheck)
		}
	}

	return len(boxCanMove) > 0, boxCanMove
}
