package main

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
)

func dayEleven() {
	input, _ := os.ReadFile("day11input.txt")
	trimmed := bytes.TrimSpace(input)
	lines := bytes.Split(trimmed, []byte(" "))

	stonesList := make([]int, 0)

	for _, v := range lines {
		asInt, _ := strconv.Atoi(string(v))
		stonesList = append(stonesList, asInt)
	}

	for i := 0; i < 35; i++ {
		stonesList = blink(stonesList)
		fmt.Println(i, len(stonesList))
	}

}

func blink(stonesList []int) []int {

	// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
	// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
	// If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.

	i := 0
	j := len(stonesList) - 1
	for i <= j {
		// Apply rules
		stone := stonesList[i]

		if stone == 0 {
			stonesList[i] = 1
			i++
			continue
		}
		if lenItoa(stone)%2 == 0 {
			toInsert := cutIntInHalf(stone)
			stonesList = append(stonesList, toInsert...)
			stonesList = removeWithoutOrder(stonesList, i)
			i++
			continue
		}
		stonesList[i] = stone * 2024
		i++
	}
	return stonesList
}

func lenItoa(i int) int {
	return len(strconv.Itoa(i))
}

func cutIntInHalf(input int) []int {
	asString := strconv.Itoa(input)
	firstPart := asString[:len(asString)/2]
	secondPart := asString[len(asString)/2:]

	asIntFirst, _ := strconv.Atoi(firstPart)
	asIntSecond, _ := strconv.Atoi(secondPart)

	return []int{asIntFirst, asIntSecond}
}

func removeWithoutOrder(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
