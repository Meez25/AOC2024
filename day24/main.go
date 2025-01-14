package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

//go:embed inputday24.txt
var inputFile []byte

func main() {
	closureMap := make(map[string]func() bool)
	startingValues, instructions := parseFile(inputFile)
	for _, value := range startingValues {
		name := string(value[0:3])
		number, _ := strconv.Atoi(string(value[5]))
		closureMap[name] = func(val int) func() bool {
			return func() bool {
				return val == 1
			}
		}(number)
	}
	for _, value := range instructions {
		spaceSplitted := bytes.Split(value, []byte(" "))
		v1 := string(spaceSplitted[0])
		v2 := string(spaceSplitted[2])
		operator := string(spaceSplitted[1])
		output := string(spaceSplitted[4])

		switch operator {
		case "XOR":
			closureMap[output] = func(in1, in2 string) func() bool {
				return func() bool {
					return closureMap[in1]() != closureMap[in2]()
				}
			}(v1, v2)
		case "OR":
			closureMap[output] = func(in1, in2 string) func() bool {
				return func() bool {
					return closureMap[in1]() || closureMap[in2]()
				}
			}(v1, v2)
		case "AND":
			closureMap[output] = func(in1, in2 string) func() bool {
				return func() bool {
					return closureMap[in1]() && closureMap[in2]()
				}
			}(v1, v2)
		}
	}
	result := make(map[string]bool)
	for k, v := range closureMap {
		result[k] = v()
	}

	var resultNumber uint64
	for i := 0; i < len(result); i++ {
		key := fmt.Sprintf("z%02d", i)
		if result[key] {
			resultNumber |= 1 << i
		}
	}

	fmt.Printf("Number: %d\n", resultNumber)

	adder := parseInstructions(instructions)
	res := make(map[string]struct{})

	// Start with all z nodes
	var q []string
	for i := range 46 {
		q = append(q, fmt.Sprintf("z%02d", i))
	}

	onZs := true
	for len(q) > 0 {
		n := len(q)
		for range n {
			out := q[0]
			q = q[1:]
			if gate, ok := adder.circuit[out]; ok {
				if onZs && !adder.rule1(out, gate) {
					res[out] = struct{}{}
				} else {
					if r := adder.rule2(out, gate); len(r) > 0 {
						res[r] = struct{}{}
					}
				}
				q = append(q, gate.in1, gate.in2)
			}
		}
		onZs = false
	}

	r := make([]string, 0, len(res))
	for k := range res {
		r = append(r, k)
	}
	sort.Strings(r)

	fmt.Println(strings.Join(r, ","))
}

const (
	Nil = iota
	AND
	OR
	XOR
)

type Gate struct {
	in1, in2 string
	typ      int
}

type Adder struct {
	circuit map[string]Gate
}

func parseInstructions(instructions [][]byte) Adder {
	circuit := map[string]Gate{}

	for _, value := range instructions {
		spaceSplitted := bytes.Split(value, []byte(" "))
		in1 := string(spaceSplitted[0])
		op := string(spaceSplitted[1])
		in2 := string(spaceSplitted[2])
		out := string(spaceSplitted[4])

		// Normalize inputs (smaller one first)
		if in1 > in2 {
			in1, in2 = in2, in1
		}

		var typ int
		switch op {
		case "AND":
			typ = AND
		case "OR":
			typ = OR
		case "XOR":
			typ = XOR
		}

		circuit[out] = Gate{in1, in2, typ}
	}

	return Adder{circuit}
}

func (a Adder) prevGate(x string) int {
	return a.circuit[x].typ
}

func (a Adder) prevX(x string) string {
	return a.circuit[x].in1
}

func (a Adder) rule1(out string, gate Gate) bool {
	// All z's come from XOR gates, except z45
	if out == "z45" {
		if gate.typ != OR {
			return false
		}
	} else if gate.typ != XOR {
		return false
	}

	// z00 comes from XOR gate with x00 and y00
	if out == "z00" {
		if gate.typ != XOR || (gate.in1 != "x00" || gate.in2 != "y00") {
			return false
		}
		return true
	}

	// No other z's come 'directly' after x's and y's
	return !strings.HasPrefix(gate.in1, "x") && !strings.HasPrefix(gate.in2, "y")
}

func (a Adder) rule2(out string, gate Gate) string {
	switch gate.typ {
	case OR:
		// OR gates's inputs are always AND gates outputs
		if a.prevGate(gate.in1) != AND {
			return gate.in1
		}
		if a.prevGate(gate.in2) != AND {
			return gate.in2
		}
	case AND:
		// Except for x00, there's never 2 AND gates in a row
		if a.prevGate(gate.in1) == AND && a.prevX(gate.in1) != "x00" {
			return gate.in1
		}
		if a.prevGate(gate.in2) == AND && a.prevX(gate.in2) != "x00" {
			return gate.in2
		}
	case XOR:
		// XOR gates that come from OR and XOR gates always output to z's
		prev1 := a.prevGate(gate.in1)
		prev2 := a.prevGate(gate.in2)
		if prev1 == XOR && prev2 == OR || prev2 == XOR && prev1 == OR {
			if !strings.HasPrefix(out, "z") {
				return out
			}
		}
	}
	return ""
}

func parseFile(file []byte) ([][]byte, [][]byte) {
	table := bytes.Split(bytes.TrimSpace(file), []byte("\n\n"))
	formattedStartingValues := bytes.Split(table[0], []byte("\n"))
	formattedInstructions := bytes.Split(table[1], []byte("\n"))
	return formattedStartingValues, formattedInstructions
}
