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
	total := 0
	step2total := 0
	var instructionsNumber [1176][2]int
	input, _ := os.ReadFile("day5input.txt")
	lines := bytes.Split(input, []byte("\n\n"))

	instructions := bytes.Split(lines[0], []byte("\n"))
	bookToCheck := bytes.Split(lines[1], []byte("\n"))

	bookToCheck = bookToCheck[:len(bookToCheck)-1]

	for i, v := range instructions {
		tableOfDigitInstruction := bytes.Split(v, []byte("|"))
		firstPage, _ := strconv.Atoi(string(tableOfDigitInstruction[0]))
		secondPage, _ := strconv.Atoi(string(tableOfDigitInstruction[1]))
		instructionsNumber[i] = [2]int{firstPage, secondPage}
	}

	for _, v := range bookToCheck {
		arrayOfPages := bytes.Split(v, []byte(","))
		pages := make([]int, len(arrayOfPages))

		for i, v := range arrayOfPages {
			digit, _ := strconv.Atoi(string(v))
			pages[i] = digit
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
			finalSort := make([]int, 0, 10)
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

	elapsed := time.Since(start)
	fmt.Println("Part 1 :", total)
	fmt.Println("Part 2 :", step2total)
	fmt.Println("In :", elapsed)
}

type Node struct {
	Name     int
	Children []*Node
}

type Graph struct {
	nodes map[int]*Node
}

func buildGraph(pages []int, instructions [1176][2]int) Graph {
	var graph Graph
	graph.nodes = make(map[int]*Node, 100)

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
	ordered := make([]int, 0, 100)
	roots := make([]Node, 0, 100)
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

func checkBookValidity(pages []int, instructions [1176][2]int) (bool, int) {
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
