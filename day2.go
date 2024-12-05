package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

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
