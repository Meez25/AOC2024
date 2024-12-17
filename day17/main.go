package main

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, _ := os.ReadFile("day17input.txt")
	trimmed := bytes.TrimSpace(input)
	programList := bytes.Split(trimmed, []byte(","))
	program := make([]int, len(programList))
	for i := range programList {
		int, _ := strconv.Atoi(string(programList[i]))
		program[i] = int
	}
	c := computer{A: 0, B: 0, C: 0, pointer: 0, program: program}
	c.start()
	fmt.Println(c.A, c.display())

	var a int64
	for i := len(program) - 1; i >= 0; i-- {
		a <<= 3

		// Try values 0-7 until we find one that produces correct output
		for b := int64(0); b < 8; b++ {
			candidate := a + b
			c := computer{A: candidate, B: 0, C: 0, pointer: 0, program: program}
			c.start()

			// Check if output matches expected program slice
			if matchesProgram(c.output, program[i:]) {
				a = candidate
				break
			}
		}
	}
	fmt.Println("P2 :", a)

}
func matchesProgram(output []int, expected []int) bool {
	if len(output) != len(expected) {
		return false
	}
	for i := range output {
		if output[i] != expected[i] {
			return false
		}
	}
	return true
}

func (c *computer) display() string {
	if len(c.output) == 0 {
		return ""
	}
	display := make([]string, len(c.output))
	for i, v := range c.output {
		display[i] = strconv.Itoa(v)
	}
	return strings.Join(display, ",")
}

func (c *computer) start() {
	instructionmap := map[int]func() error{
		0: c.adv,
		1: c.bxl,
		2: c.bst,
		3: c.jnz,
		4: c.bxc,
		5: c.out,
		6: c.bdv,
		7: c.cdv,
	}

	for c.pointer < len(c.program) {
		err := instructionmap[c.program[c.pointer]]()
		if err != nil {
			break
		}
	}
}

type computer struct {
	A       int64
	B       int64
	C       int64
	pointer int
	program []int
	output  []int
}

func (c computer) getCombo(combo int) (int64, error) {
	switch combo {
	case 0:
		return 0, nil
	case 1:
		return 1, nil
	case 2:
		return 2, nil
	case 3:
		return 3, nil
	case 4:
		return c.A, nil
	case 5:
		return int64(c.B), nil
	case 6:
		return int64(c.C), nil
	case 7:
		return 7, nil
	default:
		return 0, errors.New("Should not happend")
	}

}

func (c *computer) bst() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	modulo := combo % 8
	c.B = modulo
	c.pointer += 2
	return nil
}

func (c *computer) adv() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	result := float64(c.A) / math.Pow(float64(2), float64(combo))
	c.A = int64(result)
	c.pointer += 2
	return nil
}

func (c *computer) bxl() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	c.B = c.B ^ combo
	c.pointer += 2
	return nil
}

func (c *computer) jnz() error {
	if c.A != 0 {
		combo, err := c.getCombo(c.program[c.pointer+1])
		if err != nil {
			return err
		}
		c.pointer = int(combo)
		return nil
	}
	c.pointer += 2
	return nil
}

func (c *computer) bxc() error {
	c.B = c.B ^ c.C
	c.pointer += 2
	return nil
}

func (c *computer) out() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	c.output = append(c.output, int(combo%8))
	c.pointer += 2
	return nil
}

func (c *computer) bdv() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	result := float64(c.A) / math.Pow(float64(2), float64(combo))
	c.B = int64(result)
	c.pointer += 2
	return nil
}

func (c *computer) cdv() error {
	combo, err := c.getCombo(c.program[c.pointer+1])
	if err != nil {
		return err
	}
	result := float64(c.A) / math.Pow(float64(2), float64(combo))
	c.C = int64(result)
	c.pointer += 2
	return nil
}
