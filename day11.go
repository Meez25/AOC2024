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

	// stonesList := make([]int, 0)

	stonesMap := make(map[int]int, 0)

	for _, v := range lines {
		asInt, _ := strconv.Atoi(string(v))
		// stonesList = append(stonesList, asInt)
		stonesMap[asInt]++
	}

	// fmt.Println("PART 1")
	// for i := 0; i < 25; i++ {
	// 	stonesList = blink(stonesList)
	// 	fmt.Println(i, len(stonesList))
	// }

	fmt.Println("PART 2")

	count := 0
	for i := 0; i < 75; i++ {
		stonesMap = blinkMap(stonesMap)
		// fmt.Println(i, len(stonesMap))
	}

	for _, v := range stonesMap {
		count += v
	}
	// fmt.Println(stonesMap)
	fmt.Println(count)

}

func blink(stonesList []int) []int {

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
		}
		stonesList[i] = stone * 2024
		i++
	}
	return stonesList
}

func blinkMap(stonesMap map[int]int) map[int]int {
	secondMap := make(map[int]int, len(stonesMap))

	for k, v := range stonesMap {
		if k == 0 {
			secondMap[1] = v + secondMap[1]
			continue
		}
		if lenItoa(k)%2 == 0 {
			toInsert := cutIntInHalf(k)
			secondMap[toInsert[0]] = v + secondMap[toInsert[0]]
			secondMap[toInsert[1]] = v + secondMap[toInsert[1]]
			continue
		} else {
			secondMap[k*2024] += v
		}
	}
	return secondMap
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
