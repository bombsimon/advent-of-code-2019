package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// * ParamModePosition  -> Use the value found at slice[n]
// * ParamModeImmediate -> Use the value n
// * ParamModeRelative  -> Use the value found at slice[BASE+n]
const (
	paramModePosition = iota
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

// nolint: gochecknoglobals
var jumpMap = map[int]int{
	opCodeAdd:         4,
	opCodeMultiply:    4,
	opCodeStore:       2,
	opCodeOutput:      2,
	opCodeJumpIfTrue:  3,
	opCodeJumpIfFalse: 3,
	opCodeLessThan:    4,
	opCodeEquals:      4,
	opCodeAdjustBase:  2,
}

type computer struct {
	Input    int
	Output   []int
	Pointer  int
	Base     int
	Sequence []int
	Halted   bool
}

func main() {
	var (
		line, _        = ioutil.ReadFile(os.Args[1])
		stringSequence = strings.Split(string(line), ",")
		sequence       = make([]int, len(stringSequence))
	)

	for i := range sequence {
		sequence[i], _ = strconv.Atoi(strings.TrimSpace(stringSequence[i]))
	}

	fmt.Println("part one:", run(sequence, 1))
	fmt.Println("part two:", run(sequence, 2))
	run(sequence, 2)
}

func run(sequence []int, input int) string {
	c := computer{Input: input}

	c.Sequence = make([]int, len(sequence))
	copy(c.Sequence, sequence)

	// Update the current computer/amplifier
	c.process()

	// This is what you get without generics.
	return fmt.Sprint(strings.Trim(strings.ReplaceAll(fmt.Sprint(c.Output), " ", ","), "[]"))
}

func (c *computer) process() {
	opCode := c.Sequence[c.Pointer] % 100

	// Halt code found, stop processing.
	if opCode == opCodeHalt {
		c.Halted = true
		return
	}

	getPointer := func(argumentPosition int) int {
		var (
			pointer       = 0
			modePositions = c.Sequence[c.Pointer] / 100
			positions     = map[int]int{
				1: modePositions % 10,
				2: modePositions % 100 / 10,
				3: modePositions % 1000 / 100,
			}
		)

		switch positions[argumentPosition] {
		case paramModePosition:
			pointer = c.Sequence[c.Pointer+argumentPosition]
		case paramModeImmediate:
			pointer = c.Pointer + argumentPosition
		case paramModeRelative:
			pointer = c.Sequence[c.Pointer+argumentPosition] + c.Base
		}

		if pointer >= len(c.Sequence) {
			c.Sequence = append(c.Sequence, make([]int, pointer-len(c.Sequence)+1)...)
		}

		return pointer
	}

	sequenceFor := func(pos int) int {
		return c.Sequence[getPointer(pos)]
	}

	switch opCode {
	case opCodeAdd:
		c.Sequence[getPointer(3)] = sequenceFor(1) + sequenceFor(2)

	case opCodeMultiply:
		c.Sequence[getPointer(3)] = sequenceFor(1) * sequenceFor(2)

	case opCodeStore:
		c.Sequence[getPointer(1)] = c.Input

	case opCodeOutput:
		c.Output = append(c.Output, sequenceFor(1))

	case opCodeJumpIfTrue:
		if c.Sequence[getPointer(1)] != 0 {
			c.Pointer = sequenceFor(2) - jumpMap[opCodeJumpIfTrue]
		}

	case opCodeJumpIfFalse:
		if c.Sequence[getPointer(1)] == 0 {
			c.Pointer = sequenceFor(2) - jumpMap[opCodeJumpIfFalse]
		}

	case opCodeLessThan:
		if c.Sequence[getPointer(1)] < c.Sequence[getPointer(2)] {
			c.Sequence[getPointer(3)] = 1
		} else {
			c.Sequence[getPointer(3)] = 0
		}

	case opCodeEquals:
		if c.Sequence[getPointer(1)] == c.Sequence[getPointer(2)] {
			c.Sequence[getPointer(3)] = 1
		} else {
			c.Sequence[getPointer(3)] = 0
		}

	case opCodeAdjustBase:
		c.Base += sequenceFor(1)

	default:
		panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", c.Pointer, opCode))
	}

	c.Pointer += jumpMap[opCode]
	c.process()
}
