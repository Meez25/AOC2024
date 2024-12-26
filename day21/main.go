package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
)

//go:embed "day21input.txt"
var inputFile []byte

var numericKeypad = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{"", "0", "A"},
}

var directionalKeypad = [][]string{
	{"", "^", "A"},
	{"<", "v", ">"},
}

type Position struct {
	x    int
	y    int
	path string
}

type VisitedPosition struct {
	x int
	y int
}

var Directions = [][]int{
	{0, 1},  // DOWN
	{1, 0},  // RIGHT
	{-1, 0}, // LEFT
	{0, -1}, // UP
}

func main() {
	count := 0
	for _, line := range bytes.Split(bytes.TrimSpace(inputFile), []byte("\n")) {
		digit, _ := strconv.Atoi(string(line[:3]))
		paths := generatePossibilitiesNumericalPad(numericKeypad, line)
		possibilities := generateString(paths)
		var secondLevel []string
		var thirdLevel []string
		for i := range possibilities {
			secondLevel = append(secondLevel, generateString(generatePossibilitiesForDirectionalPad(directionalKeypad, possibilities[i]))...)
		}
		for i := range secondLevel {
			thirdLevel = append(thirdLevel, generateString(generatePossibilitiesForDirectionalPad(directionalKeypad, secondLevel[i]))...)
		}
		minValue := -1
		for i := range thirdLevel {
			if len(thirdLevel[i]) < minValue || minValue == -1 {
				minValue = len(thirdLevel[i])
			}
		}
		fmt.Println(minValue)
		count += digit * minValue
	}
	fmt.Println(count)
}

func generatePossibilitiesForDirectionalPad(keypad [][]string, input string) [][]string {
	var output [][]string
	current := Position{x: 2, y: 0}
	for _, char := range input {
		var targetX, targetY int
		switch char {
		case '^':
			targetX, targetY = 1, 0
		case '<':
			targetX, targetY = 0, 1
		case '>':
			targetX, targetY = 2, 1
		case 'v':
			targetX, targetY = 1, 1
		case 'A':
			targetX, targetY = 2, 0
		default:
			continue
		}

		paths := findPathPossibility(keypad, current.x, current.y, targetX, targetY)
		var pathToAdd []string
		for _, path := range paths {
			pathToAdd = append(pathToAdd, path.path)
		}
		output = append(output, pathToAdd)

		current.x = targetX
		current.y = targetY
	}

	return output

}

func generateString(input [][]string) []string {
	// If there are no groups, return empty slice
	if len(input) == 0 {
		return []string{}
	}

	// Start with the paths from first group
	result := make([]string, len(input[0]))
	copy(result, input[0])

	// For each subsequent group
	for i := 1; i < len(input); i++ {
		var newResult []string

		// For each existing combination
		for _, existing := range result {
			// For each path in current group
			for _, newPath := range input[i] {
				// Create new combination
				newResult = append(newResult, existing+newPath)
			}
		}
		result = newResult
	}

	return result
}

func generatePossibilitiesNumericalPad(keypad [][]string, input []byte) [][]string {
	var output [][]string
	current := Position{x: 2, y: 3}
	for _, char := range input {
		var targetX, targetY int
		switch char {
		case '7':
			targetX, targetY = 0, 0
		case '8':
			targetX, targetY = 1, 0
		case '9':
			targetX, targetY = 2, 0
		case '4':
			targetX, targetY = 0, 1
		case '5':
			targetX, targetY = 1, 1
		case '6':
			targetX, targetY = 2, 1
		case '1':
			targetX, targetY = 0, 2
		case '2':
			targetX, targetY = 1, 2
		case '3':
			targetX, targetY = 2, 2
		case '0':
			targetX, targetY = 1, 3
		case 'A':
			targetX, targetY = 2, 3
		default:
			continue
		}

		paths := findPathPossibility(keypad, current.x, current.y, targetX, targetY)
		var pathToAdd []string
		for _, path := range paths {
			pathToAdd = append(pathToAdd, path.path)
		}
		output = append(output, pathToAdd)

		current.x = targetX
		current.y = targetY
	}

	return output
}

func findPathPossibility(keypad [][]string, startX, startY, endX, endY int) []Position {
	queue := []Position{}
	start := Position{startX, startY, ""}
	queue = append(queue, start)
	visited := map[VisitedPosition]int{}
	visited[VisitedPosition{start.x, start.y}] = 0

	shortestPaths := []Position{}
	shortestLength := -1

	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		if v.x == endX && v.y == endY {
			if shortestLength == -1 {
				shortestLength = len(v.path)
				v.path += "A"
				shortestPaths = append(shortestPaths, v)
			} else if len(v.path) == shortestLength {
				v.path += "A"
				shortestPaths = append(shortestPaths, v)
			}
			continue
		}
		if shortestLength != -1 && len(v.path) > shortestLength {
			continue
		}
		for _, direction := range Directions {
			var directionArrow string
			for i, dir := range Directions {
				if dir[0] == direction[0] && dir[1] == direction[1] {
					switch i {
					case 0: // DOWN
						directionArrow += "v"
					case 1: // RIGHT
						directionArrow += ">"
					case 2: // LEFT
						directionArrow += "<"
					case 3: // UP
						directionArrow += "^"
					}
					break
				}
			}
			nextPosition := Position{x: v.x + direction[0], y: v.y + direction[1], path: v.path + directionArrow}
			if nextPosition.y < 0 || nextPosition.y >= len(keypad) || nextPosition.x < 0 || nextPosition.x >= len(keypad[0]) || keypad[nextPosition.y][nextPosition.x] == "" {
				continue
			}
			visitedPos := VisitedPosition{nextPosition.x, nextPosition.y}
			// Allow revisiting if we found a path of the same length
			pathLength := len(nextPosition.path)
			if prevLength, exists := visited[visitedPos]; !exists || pathLength <= prevLength {
				queue = append(queue, nextPosition)
				visited[visitedPos] = pathLength
			}
		}

	}
	return shortestPaths
}
