package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func dayThree() {
	input, _ := os.ReadFile("day3input.txt")
	lines := bytes.Split(input, []byte("\n"))

	total := 0

	combined := ""

	for _, line := range lines {
		combined = combined + string(line)
	}

	tableOfGoodStuff := extractGoodStuff(combined)

	for _, e := range tableOfGoodStuff {
		total = total + calculateFromString(e)
	}

	fmt.Println(total)
}

func extractGoodStuff(input string) []string {
	var output []string
	ignore := false
	skipTo := 0
	_ = ignore

	for i, letter := range input {
		if letter == 'd' {
			if input[i:i+7] == "don't()" {
				if !ignore {
					output = append(output, input[skipTo:i])
				}
				ignore = true
				skipTo = i + 7
			}

			if input[i:i+4] == "do()" {
				if !ignore {
					output = append(output, input[skipTo:i])
				}
				ignore = false
				skipTo = i + 4
			}
		}

		if i < skipTo {
			continue
		}

		if i == len(input)-1 {
			if !ignore {
				output = append(output, input[skipTo:i])
			}
		}

	}

	return output

}

func calculateFromString(input string) int {
	total := 0
	exp := regexp.MustCompile(`mul\(\d{1,3},\d{1,3}\)`)
	expToGetNumbers := regexp.MustCompile(`(\d{1,3}),(\d{1,3})`)

	tableWithMatch := exp.FindAllString(input, -1)

	for _, v := range tableWithMatch {
		tableOfDigit := expToGetNumbers.FindAllStringSubmatch(v, -1)
		firstDigit, _ := strconv.Atoi(tableOfDigit[0][1])
		secondDigit, _ := strconv.Atoi(tableOfDigit[0][2])
		total = total + (firstDigit * secondDigit)
	}
	return total
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
