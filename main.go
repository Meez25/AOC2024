package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strconv"
)

func main() {
	dayFive()
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

func dayFour() {
	input, _ := os.ReadFile("day4input.txt")
	lines := bytes.Split(input, []byte("\n"))

	lines = lines[:len(lines)-1]

	match := 0

	for i := 0; i < len(lines); i++ {
		for y := 0; y < len(lines[i]); y++ {
			if string(lines[i][y]) == "X" {
				// Check left
				if y > 2 && string(lines[i][y-1]) == "M" && string(lines[i][y-2]) == "A" && string(lines[i][y-3]) == "S" {
					fmt.Println("line:", i, "row:", y)
					match++
				}
				// Check right
				if y < len(lines[i])-3 && string(lines[i][y+1]) == "M" && string(lines[i][y+2]) == "A" && string(lines[i][y+3]) == "S" {
					fmt.Println("line:", i, "row:", y)
					match++
				}
				// Check up
				if i > 2 && string(lines[i-1][y]) == "M" && string(lines[i-2][y]) == "A" && string(lines[i-3][y]) == "S" {
					fmt.Println("line:", i, "row:", y)
					match++
				}
				// Check down
				if i < len(lines)-3 && string(lines[i+1][y]) == "M" && string(lines[i+2][y]) == "A" && string(lines[i+3][y]) == "S" {
					fmt.Println("line:", i, "row:", y)
					match++
				}
				// Check up right
				if i > 2 && y < len(lines[i])-3 && string(lines[i-1][y+1]) == "M" && string(lines[i-2][y+2]) == "A" && string(lines[i-3][y+3]) == "S" {
					match++
					fmt.Println("line:", i, "row:", y)
				}
				// Check up left
				if i > 2 && y > 2 && string(lines[i-1][y-1]) == "M" && string(lines[i-2][y-2]) == "A" && string(lines[i-3][y-3]) == "S" {
					fmt.Println("line:", i, "row:", y)
					match++
				}
				// Check down right
				if i < len(lines)-3 && y < len(lines[i])-3 && string(lines[i+1][y+1]) == "M" && string(lines[i+2][y+2]) == "A" && string(lines[i+3][y+3]) == "S" {
					match++
					fmt.Println("line:", i, "row:", y)
				}
				// Check down left
				if i < len(lines)-3 && y > 2 && string(lines[i+1][y-1]) == "M" && string(lines[i+2][y-2]) == "A" && string(lines[i+3][y-3]) == "S" {
					match++
					fmt.Println("line:", i, "row:", y)
				}
			}
		}
	}

	matchPart2 := 0

	for i := 0; i < len(lines); i++ {
		for y := 0; y < len(lines[i]); y++ {
			if string(lines[i][y]) == "A" {
				// Check up right
				if i > 0 && y > 0 && i < len(lines)-1 && y < len(lines[i])-1 {
					if (string(lines[i-1][y-1]) == "M" && string(lines[i+1][y+1]) == "S") || (string(lines[i-1][y-1]) == "S" && string(lines[i+1][y+1]) == "M") {
						if (string(lines[i+1][y-1]) == "M" && string(lines[i-1][y+1]) == "S") || (string(lines[i+1][y-1]) == "S" && string(lines[i-1][y+1]) == "M") {
							matchPart2++
						}
					}
				}
			}
		}
	}

	fmt.Println("part1: ", match)
	fmt.Println("part2: ", matchPart2)
}

func dayFive() {
	total := 0
	step2total := 0
	var instructionsNumber [][]int
	input, _ := os.ReadFile("day5input.txt")
	lines := bytes.Split(input, []byte("\n\n"))

	instructions := bytes.Split(lines[0], []byte("\n"))
	bookToCheck := bytes.Split(lines[1], []byte("\n"))

	bookToCheck = bookToCheck[:len(bookToCheck)-1]

	for _, v := range instructions {
		tableOfDigitInstruction := bytes.Split(v, []byte("|"))
		firstPage, _ := strconv.Atoi(string(tableOfDigitInstruction[0]))
		secondPage, _ := strconv.Atoi(string(tableOfDigitInstruction[1]))
		instructionsNumber = append(instructionsNumber, []int{firstPage, secondPage})
	}

	for _, v := range bookToCheck {
		var pages []int
		arrayOfPages := bytes.Split(v, []byte(","))

		for _, v := range arrayOfPages {
			digit, _ := strconv.Atoi(string(v))
			pages = append(pages, digit)
		}

		// Step 1

		isCorrect, middleNumber1 := checkBookValidity(pages, instructionsNumber)
		if isCorrect {
			total = total + middleNumber1
		} else {

			// Step 2
			graph := buildGraph(pages, instructionsNumber)
			order := graph.TopologicalSort()
			// Remove useless nodes from book
			var finalSort []int
			for _, v := range order {
				if slices.Contains(pages, v) {
					finalSort = append(finalSort, v)
				}
			}
			if len(finalSort) > 0 {
				step2total = step2total + finalSort[len(finalSort)/2]
			}
		}
	}

	fmt.Println(total)
	fmt.Println(step2total)
}

type Node struct {
	Name     int
	Children []*Node
}

type Graph struct {
	nodes map[int]*Node
}

func buildGraph(pages []int, instructions [][]int) Graph {
	var graph Graph
	graph.nodes = make(map[int]*Node)

	// Create nodes
	for _, page := range pages {
		graph.nodes[page] = &Node{Name: page}
	}

	// Add edges only if both pages exist
	for _, v := range instructions {
		if _, exists1 := graph.nodes[v[0]]; exists1 {
			if _, exists2 := graph.nodes[v[1]]; exists2 {
				graph.nodes[v[0]].Children = append(graph.nodes[v[0]].Children, graph.nodes[v[1]])
			}
		}
	}
	return graph
}

func (g Graph) TopologicalSort() []int {
	var ordered []int
	var roots []Node
	incomingEdges := make(map[int]int)

	// Count incoming edges
	for _, node := range g.nodes {
		for _, child := range node.Children {
			incomingEdges[child.Name]++
		}
	}

	// Find initial roots
	for name, node := range g.nodes {
		if incomingEdges[name] == 0 {
			roots = append(roots, *node)
		}
	}

	// Process nodes
	for len(roots) > 0 {
		// Remove and process current root
		node := roots[len(roots)-1]
		roots = roots[:len(roots)-1]
		ordered = append(ordered, node.Name)

		// Process children
		for _, child := range node.Children {
			incomingEdges[child.Name]--
			if incomingEdges[child.Name] == 0 {
				roots = append(roots, *child)
			}
		}
	}

	return ordered
}

func checkBookValidity(book []int, instructions [][]int) (bool, int) {
	var pages = book

	for i, v := range pages {
		// If both number from the instructions are in the book, check if the right side of the instruction doesn't exist in the following numbers, if so, the book is invalid
		for _, instruction := range instructions {
			if v == instruction[0] {
				if !slices.Contains(pages[i:], instruction[1]) && slices.Contains(pages, instruction[0]) && slices.Contains(pages, instruction[1]) {
					// Try switching both values from the pages
					return false, 0
				}
			}
		}
	}

	middlePage := pages[len(pages)/2]

	return true, middlePage
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
