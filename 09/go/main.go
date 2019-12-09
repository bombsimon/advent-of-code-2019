package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

// paramMode represents the parameter mode which can be positional (the value at
// the position) or immediate (the actual value).
type paramMode int

// * ParamModePosition -> Use the value found at slice[n]
// * ParamModeImmediate -> Use the value n
// * ParamModeRelative -> Use the value found at slice[BASE+n]
const (
	paramModePosition paramMode = iota
	paramModeImmediate
	paramModeRelative
)

const (
	opCodeAdd         = 1
	opCodeMultiply    = 2
	opCodeStore       = 3
	opCodeOutput      = 4
	opCodeJumpIfTrue  = 5
	opCodeJumpIfFalse = 6
	opCodeLessThan    = 7
	opCodeEquals      = 8
	opCodeAdjustBase  = 9
	opCodeHalt        = 99
)

type computer struct {
	Name       string
	Phase      int
	Input      int
	InputCount int
	Output     []int
	Pointer    int
	Base       int
	Sequence   []int
	Halted     bool
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing file as input")
	}

	line, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	var (
		stringSequence = strings.Split(string(line), ",")
		sequence       = make([]int, len(stringSequence))
	)

	for i := range sequence {
		sv := strings.TrimSpace(stringSequence[i])
		v, _ := strconv.Atoi(sv)
		sequence[i] = v
	}

	fmt.Println("part one")
	run(sequence, 1)

	fmt.Println("part two")
	run(sequence, 2)
}

func run(sequence []int, phase int) {
	c := computer{
		Phase: phase,
	}

	if c.Sequence == nil {
		c.Sequence = make([]int, len(sequence))
		copy(c.Sequence, sequence)
	}

	// Update the current computer/amplifier
	c.process()

	// This is what you get without generics.
	fmt.Println(strings.Trim(strings.ReplaceAll(fmt.Sprint(c.Output), " ", ","), "[]"))
}

func (c *computer) process() {
	opCode := c.Sequence[c.Pointer]

	// Halt code found, stop processing.
	if opCode == opCodeHalt {
		c.Halted = true
		return
	}

	// Pointer is outside of list range.
	if c.Pointer > len(c.Sequence) {
		return
	}

	expand := func(s []int, n int) []int {
		if len(s)-1 > n {
			return s
		}

		var toAdd = make([]int, n-len(s)+1)

		s = append(s, toAdd...)

		return s
	}

	// Ensure we have at least three extra positions after the pointer to be
	// able to assign all variables at once.
	c.Sequence = expand(c.Sequence, c.Pointer+3)

	// Get firwst two arguments, parse their input code and the actual op code
	// for larger numbers.
	first, second, pos := c.Sequence[c.Pointer+1], c.Sequence[c.Pointer+2], c.Sequence[c.Pointer+3]
	firstMode, secondMode, thirdMode, opCode := parseOpCode(opCode)

	c.Sequence = expand(c.Sequence, first)
	c.Sequence = expand(c.Sequence, second)
	c.Sequence = expand(c.Sequence, pos)

	switch firstMode {
	case paramModePosition:
		if first >= 0 {
			first = c.Sequence[first]
		}
	case paramModeImmediate:
		break
	case paramModeRelative:
		// If we're storing a value, just add the destination to where.
		if opCode == opCodeStore {
			first = c.Base + first
		} else {
			first = c.Sequence[c.Base+first]
		}
	}

	switch secondMode {
	case paramModePosition:
		if second >= 0 {
			second = c.Sequence[second]
		}
	case paramModeImmediate:
		break
	case paramModeRelative:
		second = c.Sequence[c.Base+second]
	}

	switch thirdMode {
	case paramModePosition, paramModeImmediate:
		break
	case paramModeRelative:
		pos += c.Base
	}

	switch opCode {
	case opCodeAdd:
		c.Sequence[pos] = first + second
		c.Pointer += 4

	case opCodeMultiply:
		c.Sequence[pos] = first * second
		c.Pointer += 4

	case opCodeStore:
		input := c.Phase
		if c.InputCount > 0 {
			input = c.Input
		}

		c.Sequence[first] = input
		c.Pointer += 2
		c.InputCount++

	case opCodeOutput:
		c.Output = append(c.Output, first)
		c.Pointer += 2

		// Puase process if opcode > 4
		if c.Phase > 4 {
			return
		}

	case opCodeJumpIfTrue:
		if first != 0 {
			c.Pointer = second
		} else {
			c.Pointer += 3
		}

	case opCodeJumpIfFalse:
		if first == 0 {
			c.Pointer = second
		} else {
			c.Pointer += 3
		}

	case opCodeLessThan:
		if first < second {
			c.Sequence[pos] = 1
		} else {
			c.Sequence[pos] = 0
		}

		c.Pointer += 4

	case opCodeEquals:
		if first == second {
			c.Sequence[pos] = 1
		} else {
			c.Sequence[pos] = 0
		}

		c.Pointer += 4

	case opCodeAdjustBase:
		c.Base += first
		c.Pointer += 2

	default:
		panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", c.Pointer, opCode))
	}

	c.process()
}

func parseOpCode(opCode int) (paramMode, paramMode, paramMode, int) {
	var (
		code          = opCode % 100
		modePositions = opCode / 100
		firstMode     = modePositions % 10
		secondMode    = modePositions % 100 / 10
		thirdMode     = modePositions % 1000 / 100
	)

	return paramMode(firstMode), paramMode(secondMode), paramMode(thirdMode), code
}
