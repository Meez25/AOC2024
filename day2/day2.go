package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	input, _ := os.ReadFile("day2input.txt")
	lines := bytes.Split(input, []byte("\n"))

	numberOfSafeLevels := 0
	numberOfSafeLevelsTwo := 0

	for _, v := range lines[:len(lines)-1] {
		levels := make([]int, 0, 10)
		separatedValues := bytes.Split(v, []byte(" "))
		for i := range separatedValues {
			number, _ := strconv.Atoi(string(separatedValues[i]))
			levels = append(levels, number)
		}
		if isLevelSafe(levels) {
			numberOfSafeLevels++
		}
		if isLevelSafePartTwo(levels) {
			numberOfSafeLevelsTwo++
		}
	}

	elapsed := time.Since(start)
	fmt.Println("Part 1:", numberOfSafeLevels)
	fmt.Println("Part 2:", numberOfSafeLevelsTwo)
	fmt.Println("In :", elapsed)
}

func isLevelSafePartTwo(input []int) bool {
	lenght := len(input)

	if isLevelSafe(input) {
		return true
	}

	for i := 0; i < lenght; i++ {
		subslice := OriginalRemoveIndex(input, i)

		isSafe := isLevelSafe(subslice)
		if isSafe {
			return true
		}

	}

	return false
}

func isLevelSafe(input []int) bool {
	isLevelSafe := true
	lenght := len(input)
	var isIncreasing bool

	if input[0] < input[1] {
		isIncreasing = true
	} else if input[0] > input[1] {
		isIncreasing = false
	} else {
		return false
	}

	for i := 0; i < lenght-1; i++ {
		if isIncreasing {
			if input[i] > input[i+1] || input[i+1]-input[i] > 3 || input[i+1]-input[i] < 1 {
				return false
			}
		}
		if !isIncreasing {
			if input[i] < input[i+1] || input[i]-input[i+1] > 3 || input[i]-input[i+1] < 1 {
				return false
			}
		}
	}

	return isLevelSafe
}

func OriginalRemoveIndex(arr []int, pos int) []int {
	newarray := make([]int, 0, 10)
	for i := range arr {
		if i != pos {
			newarray = append(newarray, arr[i])
		}
	}
	return newarray
}
