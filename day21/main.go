package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"strconv"
	"time"
)

//go:embed day21input.txt
var inputFile []byte

type Node struct {
	value     string
	direction string
	x         int
	y         int
}

type NodeVisitedSet struct {
	value string
	x     int
	y     int
}

var directionalKeypad = [][]string{
	{"", "^", "A"},
	{"<", "v", ">"},
}

var numericKeypad = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{"", "0", "A"},
}

type direction struct {
	direction [2]int
	arrow     string
}

var directions = []direction{
	{direction: [2]int{0, -1}, arrow: "^"}, // UP
	{direction: [2]int{0, 1}, arrow: "v"},  // DOWN
	{direction: [2]int{-1, 0}, arrow: "<"}, // LEFT
	{direction: [2]int{1, 0}, arrow: ">"},  // RIGHT
}

func main() {
	start := time.Now()
	total := 0
	cache := make(map[cacheKey]int, 0)
	for _, line := range bytes.Split(bytes.TrimSpace(inputFile), []byte("\n")) {
		results := convertNumsToDir(string(line))
		minCost := 0
		for _, result := range results {
			cost := convertToShortest(result, 25, cache)
			if cost < minCost || minCost == 0 {
				minCost = cost
			}
		}
		digitPart, _ := strconv.Atoi(string(line[:3]))
		total += digitPart * minCost
	}
	fmt.Println("Part 2", total, "in", time.Since(start))
}

type cacheKey struct {
	input string
	depth int
}

func convertToShortest(input string, depth int, cache map[cacheKey]int) int {
	key := cacheKey{input, depth}
	if val, ok := cache[key]; ok {
		return val
	}
	if depth == 0 {
		cache[key] = len(input)
		return len(input)
	}
	input = "A" + input
	// Check the min cost of children
	inputCost := 0
	for i := 0; i < len(input)-1; i++ {
		// For each 2 digits children of the INPUT
		possibilities := numKeyDirToDirection(directionalKeypad, string(input[i]), string(input[i+1]))
		minCost := 0
		for _, possibility := range possibilities {
			cost := convertToShortest(possibility, depth-1, cache)
			if cost < minCost || minCost == 0 {
				minCost = cost
			}
		}
		inputCost += minCost
	}

	cache[key] = inputCost
	return inputCost
}

func convertDirToDir(input []string) []string {
	lists := make([]string, 0)
	for _, word := range input {
		results := make([]string, 0)

		current := "A"

		for i := range word {
			result := numKeyDirToDirection(directionalKeypad, current, string(word[i]))
			if len(results) == 0 {
				results = result
				current = string(word[i])
				continue
			}

			toReplace := make([]string, 0)

			for i := range result {
				for y := range results {
					toReplace = append(toReplace, results[y]+result[i])
				}
			}
			results = toReplace

			current = string(word[i])
		}
		lists = append(lists, results...)
	}
	return lists
}

func convertNumsToDir(input string) []string {
	results := make([]string, 0)

	current := "A"

	for i := range input {
		result := numKeyPadToDirection(numericKeypad, current, string(input[i]))
		if len(results) == 0 {
			results = result
			current = string(input[i])
			continue
		}

		toReplace := make([]string, 0)

		for i := range result {
			for y := range results {
				toReplace = append(toReplace, results[y]+result[i])
			}
		}
		results = toReplace

		current = string(input[i])
	}
	return results
}

func numKeyDirToDirection(directionalKeypad [][]string, startButton, wanted string) []string {
	// Find all path to the wanted char
	// Return a list of possibilities to navigate using arrow keys
	queue := make([]Node, 0)
	visited := make(map[string]bool)
	var start Node

	for y := range directionalKeypad {
		for x := range directionalKeypad[y] {
			if directionalKeypad[y][x] == startButton {
				start = Node{value: startButton, direction: "", x: x, y: y}
			}
		}
	}
	queue = append(queue, start)
	visited[start.direction] = true
	possibilities := make([]string, 0)
	shortestPath := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.value == wanted {
			if shortestPath == 0 || len(current.direction) <= shortestPath {
				shortestPath = len(current.direction)
				possibilities = append(possibilities, current.direction+"A")
			} else {
				return possibilities
			}
		}
		for _, direction := range directions {
			nextX := current.x + direction.direction[0]
			nextY := current.y + direction.direction[1]
			if nextX > len(directionalKeypad[0])-1 || nextX < 0 || nextY > len(directionalKeypad)-1 || nextY < 0 {
				continue
			}
			nextNode := Node{value: directionalKeypad[nextY][nextX], direction: current.direction + direction.arrow, x: nextX, y: nextY}
			if nextNode.value == "" {
				continue
			}
			if visited[nextNode.direction] == true {
				continue
			}
			queue = append(queue, nextNode)
			visited[nextNode.direction] = true
		}
	}
	return possibilities
}

func numKeyPadToDirection(numericPad [][]string, startButton, wanted string) []string {
	// Find all path to the wanted char
	// Return a list of possibilities to navigate using arrow keys
	queue := make([]Node, 0)
	visited := make(map[string]bool)
	var start Node

	for y := range numericPad {
		for x := range numericPad[y] {
			if numericPad[y][x] == startButton {
				start = Node{value: startButton, direction: "", x: x, y: y}
			}
		}
	}
	queue = append(queue, start)
	visited[start.direction] = true
	possibilities := make([]string, 0)
	shortestPath := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current.value == wanted {
			if shortestPath == 0 || len(current.direction) <= shortestPath {
				shortestPath = len(current.direction)
				possibilities = append(possibilities, current.direction+"A")
			} else {
				return possibilities
			}
		}
		for _, direction := range directions {
			nextX := current.x + direction.direction[0]
			nextY := current.y + direction.direction[1]
			if nextX > len(numericKeypad[0])-1 || nextX < 0 || nextY > len(numericKeypad)-1 || nextY < 0 {
				continue
			}
			nextNode := Node{value: numericPad[nextY][nextX], direction: current.direction + direction.arrow, x: nextX, y: nextY}
			if nextNode.value == "" {
				continue
			}
			if visited[nextNode.direction] == true {
				continue
			}
			queue = append(queue, nextNode)
			visited[nextNode.direction] = true
		}
	}
	return possibilities
}
