package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
)

func main() {
	dayTwo()
}

func dayOne() {
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

	fmt.Println(sumDistance)

	for _, valueLeft := range leftDigits {
		found := 0
		for _, valueRight := range rightDigits {
			if valueRight == valueLeft {
				found = found + 1
			}
		}
		difference = difference + valueLeft*found
	}

	fmt.Println(difference)
}

func dayTwo() {
	input, _ := os.ReadFile("day2input.txt")
	lines := bytes.Split(input, []byte("\n"))

	numberOfSafeLevels := 0
	numberOfSafeLevelsTwo := 0

	for _, v := range lines[:len(lines)-1] {
		var levels []int
		separatedValues := bytes.Split(v, []byte(" "))
		for i := range separatedValues {
			number, _ := strconv.Atoi(string(separatedValues[i]))
			levels = append(levels, number)
		}
		if isLevelSafe(levels) {
			numberOfSafeLevels = numberOfSafeLevels + 1
		}
		if isLevelSafePartTwo(levels) {
			numberOfSafeLevelsTwo = numberOfSafeLevelsTwo + 1
		}
	}

	fmt.Println(numberOfSafeLevels)
	fmt.Println(numberOfSafeLevelsTwo)
}

func isLevelSafePartTwo(input []int) bool {
	lenght := len(input)

	if isLevelSafe(input) {
		return true
	}
	fmt.Println("BASE", input)

	for i := 0; i < lenght; i++ {
		subslice := OriginalRemoveIndex(input, i)
		fmt.Println("subslice", subslice, i)

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
	newarray := []int{}
	for i := range arr {
		if i != pos {
			newarray = append(newarray, arr[i])
		}
	}
	return newarray
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}
