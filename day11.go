package main

import (
	"bytes"
	"fmt"
	"os"
	"slices"
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

	for i := 0; i < 25; i++ {
		stonesList = blink(stonesList)
	}
	fmt.Println(len(stonesList))

}

func blink(stonesList []int) []int {
	copyStone := slices.Clone(stonesList)

	// If the stone is engraved with the number 0, it is replaced by a stone engraved with the number 1.
	// If the stone is engraved with a number that has an even number of digits, it is replaced by two stones. The left half of the digits are engraved on the new left stone, and the right half of the digits are engraved on the new right stone. (The new numbers don't keep extra leading zeroes: 1000 would become stones 10 and 0.)
	// If none of the other rules apply, the stone is replaced by a new stone; the old stone's number multiplied by 2024 is engraved on the new stone.

	i := 0
	for i < len(copyStone) {
		stone := copyStone[i]
		// Apply rules

		if stone == 0 {
			copyStone[i] = 1
			i++
			continue
		}
		if lenItoa(stone)%2 == 0 {
			toInsert := cutIntInHalf(stone)
			copyStone = slices.Insert(copyStone, i, toInsert...)
			copyStone = slices.Delete(copyStone, i+2, i+3)
			i = i + 2
			continue
		}
		copyStone[i] = stone * 2024
		i++
	}

	return copyStone
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
