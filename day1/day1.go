package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	input, _ := os.ReadFile("day1input.txt")
	lines := bytes.Split(input, []byte("\n"))

	var leftDigits []int
	var rightDigits []int

	sumDistance := 0
	difference := 0

	for _, v := range lines[:len(lines)-1] {
		separatedValues := bytes.Split(v, []byte("   "))
		firstDigit, _ := strconv.Atoi(string(separatedValues[0]))
		secondDigit, _ := strconv.Atoi(string(separatedValues[1]))
		leftDigits = append(leftDigits, firstDigit)
		rightDigits = append(rightDigits, secondDigit)
	}

	slices.Sort(leftDigits)
	slices.Sort(rightDigits)

	for i := range leftDigits {
		difference := diff(leftDigits[i], rightDigits[i])
		sumDistance = sumDistance + difference
	}

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", sumDistance, "in", elapsed)
	start = time.Now()

	for _, valueLeft := range leftDigits {
		found := 0
		for _, valueRight := range rightDigits {
			if valueRight == valueLeft {
				found = found + 1
			}
		}
		difference = difference + valueLeft*found
	}

	elapsed = time.Since(start)
	fmt.Println("Part 2 :", difference, "in", elapsed)
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}
