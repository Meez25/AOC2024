package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"strconv"
)

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
