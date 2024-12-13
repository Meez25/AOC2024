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
	X    int64
	Y    int64
}

type prize struct {
	X int64
	Y int64
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
		firstButton := button{name: "A", X: stringToInt64(results[0][1]), Y: stringToInt64(results[0][2])}
		results = regexp.FindAllStringSubmatch(string(lines[1]), -1)
		secondButton := button{name: "B", X: stringToInt64(results[0][1]), Y: stringToInt64(results[0][2])}
		results = regexp.FindAllStringSubmatch(string(lines[2]), -1)
		prize := prize{X: stringToInt64(results[0][1]), Y: stringToInt64(results[0][2])}

		fmt.Println(firstButton, secondButton, prize)

		// result := solvePuzzle(firstButton, secondButton, prize)
		result := solvePuzzleP2(firstButton, secondButton, prize)
		fmt.Println(result)
		grandTotal += result
	}

	fmt.Println(grandTotal)

}

func stringToInt(input string) int {
	asInt, _ := strconv.Atoi(input)
	return asInt
}

func stringToInt64(input string) int64 {
	asInt, _ := strconv.Atoi(input)
	return int64(asInt)
}

// func solvePuzzle(firstButton, secondButton button, prize prize) int {
// 	// Push A cost 3 token, and push B cost 1 token
// 	optimizedPrice := 0
// 	initialNode := NodeDay13{X: 0, Y: 0, cost: 0}
// 	visited := make(map[NodeDay13]bool, 0)
// 	visited[initialNode] = true
// 	queue := make([]NodeDay13, 0)
// 	queue = append(queue, initialNode)
//
// 	for len(queue) > 0 {
// 		node := queue[0]
// 		queue = queue[1:]
// 		if node.X == prize.X && node.Y == prize.Y {
// 			return node.cost
// 		}
// 		if node.pressedATimes >= 100 || node.pressedBTimes >= 100 || node.Y > prize.Y || node.X > prize.X {
// 			continue
// 		}
// 		newNodeButtonA := NodeDay13{X: node.X + firstButton.X, Y: node.Y + firstButton.Y, cost: node.cost + 3, pressedATimes: node.pressedATimes + 1}
// 		newNodeButtonB := NodeDay13{X: node.X + secondButton.X, Y: node.Y + secondButton.Y, cost: node.cost + 1, pressedBTimes: node.pressedBTimes + 1}
// 		if !visited[newNodeButtonB] {
// 			queue = append(queue, newNodeButtonB)
// 			visited[newNodeButtonB] = true
// 		}
// 		if !visited[newNodeButtonA] {
// 			queue = append(queue, newNodeButtonA)
// 			visited[newNodeButtonA] = true
// 		}
// 	}
//
// 	return optimizedPrice
// }

func solvePuzzleP2(firstButton, secondButton button, prize prize) int {
	// Add the offset to prize coordinates
	prize.X += 10000000000000
	prize.Y += 10000000000000

	// Check possibility with GCD
	GCDX := GCD(firstButton.X, secondButton.X)
	GCDY := GCD(firstButton.Y, secondButton.Y)

	if prize.X%GCDX != 0 {
		return 0
	}
	if prize.Y%GCDY != 0 {
		return 0
	}

	eq1 := firstButton.X * secondButton.Y
	eq2 := firstButton.Y * secondButton.X
	coefA := eq1 - eq2
	rightSide := prize.X*secondButton.Y - prize.Y*secondButton.X

	if coefA == 0 {
		return 0
	}

	if rightSide%coefA != 0 {
		return 0
	}
	a := rightSide / coefA

	remainingX := prize.X - firstButton.X*a
	if remainingX%secondButton.X != 0 {
		return 0
	}
	b := remainingX / secondButton.X

	if firstButton.Y*a+secondButton.Y*b != prize.Y {
		return 0
	}

	if a < 0 || b < 0 {
		return 0
	}

	return int(3*a + b)
}

func GCD(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
