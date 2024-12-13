package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

type button struct {
	name string
	X    int
	Y    int
}

type prize struct {
	X int
	Y int
}

type NodeDay13 struct {
	cost          int
	X             int
	Y             int
	pressedATimes int
	pressedBTimes int
}

func dayThirteen() {
	input, _ := os.ReadFile("day13input.txt")
	trimmed := bytes.TrimSpace(input)
	blocks := bytes.Split(trimmed, []byte("\n\n"))

	regexp := regexp.MustCompile(`(\d+).*?(\d+)`)

	grandTotal := 0

	for _, block := range blocks {
		lines := bytes.Split(block, []byte("\n"))
		results := regexp.FindAllStringSubmatch(string(lines[0]), -1)
		firstButton := button{name: "A", X: stringToInt(results[0][1]), Y: stringToInt(results[0][2])}
		results = regexp.FindAllStringSubmatch(string(lines[1]), -1)
		secondButton := button{name: "B", X: stringToInt(results[0][1]), Y: stringToInt(results[0][2])}
		results = regexp.FindAllStringSubmatch(string(lines[2]), -1)
		prize := prize{X: stringToInt(results[0][1]), Y: stringToInt(results[0][2])}

		fmt.Println(firstButton, secondButton, prize)

		result := solvePuzzle(firstButton, secondButton, prize)
		fmt.Println(result)
		grandTotal += result
	}

	fmt.Println(grandTotal)

}

func stringToInt(input string) int {
	asInt, _ := strconv.Atoi(input)
	return asInt
}

func solvePuzzle(firstButton, secondButton button, prize prize) int {
	// Push A cost 3 token, and push B cost 1 token
	optimizedPrice := 0
	initialNode := NodeDay13{X: 0, Y: 0, cost: 0}
	visited := make(map[NodeDay13]bool, 0)
	visited[initialNode] = true
	queue := make([]NodeDay13, 0)
	queue = append(queue, initialNode)

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.X == prize.X && node.Y == prize.Y {
			return node.cost
		}
		if node.pressedATimes >= 100 || node.pressedBTimes >= 100 || node.Y > prize.Y || node.X > prize.X {
			continue
		}
		newNodeButtonA := NodeDay13{X: node.X + firstButton.X, Y: node.Y + firstButton.Y, cost: node.cost + 3, pressedATimes: node.pressedATimes + 1}
		newNodeButtonB := NodeDay13{X: node.X + secondButton.X, Y: node.Y + secondButton.Y, cost: node.cost + 1, pressedBTimes: node.pressedBTimes + 1}
		if !visited[newNodeButtonB] {
			queue = append(queue, newNodeButtonB)
			visited[newNodeButtonB] = true
		}
		if !visited[newNodeButtonA] {
			queue = append(queue, newNodeButtonA)
			visited[newNodeButtonA] = true
		}
	}

	return optimizedPrice
}
