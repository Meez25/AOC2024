package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"slices"
	"sort"
	"strings"
)

//go:embed inputday23.txt
var inputFile []byte

func main() {
	graph := buildGraph(inputFile)
	fmt.Println(countSetsOfThree(graph))
	fmt.Println(len(graph))
	max := bronKerboschClassic(graph)
	fmt.Println(strings.Join(max[0], ","))
}

func bronKerboschClassic(graph map[string][]string) [][]string {
	var maximalCliques [][]string

	// Get all vertices
	vertices := make([]string, 0, len(graph))
	for v := range graph {
		vertices = append(vertices, v)
	}

	// Helper function for set intersection
	intersect := func(a, b []string) []string {
		result := make([]string, 0)
		for _, v := range a {
			if slices.Contains(b, v) {
				result = append(result, v)
			}
		}
		return result
	}

	// Helper function for set difference (a - b)
	difference := func(a, b []string) []string {
		result := make([]string, 0)
		for _, v := range a {
			if !slices.Contains(b, v) {
				result = append(result, v)
			}
		}
		return result
	}

	// Recursive Bron-Kerbosch without pivot
	var bronKerbosch func(r, p, x []string)
	bronKerbosch = func(r, p, x []string) {
		// If both P and X are empty, we found a maximal clique
		if len(p) == 0 && len(x) == 0 {
			clique := make([]string, len(r))
			copy(clique, r)
			slices.Sort(clique)
			maximalCliques = append(maximalCliques, clique)
			return
		}

		// Try each vertex in P
		for _, v := range p {
			neighbors := graph[v]

			// Recursive call with:
			// R ∪ {v}
			newR := append([]string{}, r...)
			newR = append(newR, v)

			// P ∩ N(v)
			newP := intersect(p, neighbors)

			// X ∩ N(v)
			newX := intersect(x, neighbors)

			bronKerbosch(newR, newP, newX)

			// Move v from P to X
			p = difference(p, []string{v})
			x = append(x, v)
		}
	}

	// Initial call with empty R, all vertices in P, empty X
	bronKerbosch([]string{}, vertices, []string{})

	// Find maximum cliques
	if len(maximalCliques) == 0 {
		return nil
	}

	// Find the size of the largest clique
	maxSize := len(maximalCliques[0])
	for _, clique := range maximalCliques {
		if len(clique) > maxSize {
			maxSize = len(clique)
		}
	}

	// Return only the maximum cliques
	var result [][]string
	for _, clique := range maximalCliques {
		if len(clique) == maxSize {
			result = append(result, clique)
		}
	}

	return result
}

func countSetsOfThree(graph map[string][]string) int {
	set := make([]string, 0)
	for k := range graph {
		children := getChild(graph, k)
		for _, child := range children {
			grandChildren := getChild(graph, child)
			for _, grandChild := range grandChildren {
				greatGrandChildren := getChild(graph, grandChild)
				if slices.Contains(greatGrandChildren, k) {
					combinaison := []string{k, child, grandChild}
					if k[0] == 't' || child[0] == 't' || grandChild[0] == 't' {
						sort.Strings(combinaison)
						if !slices.Contains(set, strings.Join(combinaison, "")) {
							set = append(set, strings.Join(combinaison, ""))
						}
					}
				}
			}
		}
	}
	return len(set)
}

func getChild(graph map[string][]string, current string) []string {
	return graph[current]
}

func buildGraph(input []byte) map[string][]string {
	graph := make(map[string][]string)
	for _, line := range bytes.Split(bytes.TrimSpace(input), []byte("\n")) {
		extractedPair := bytes.Split(line, []byte("-"))
		pair1 := string(extractedPair[0])
		pair2 := string(extractedPair[1])
		if _, ok := graph[pair1]; !ok {
			graph[pair1] = []string{}
		}
		if _, ok := graph[pair2]; !ok {
			graph[pair2] = []string{}
		}
		graph[pair1] = append(graph[pair1], pair2)
		graph[pair2] = append(graph[pair2], pair1)
	}
	return graph
}
