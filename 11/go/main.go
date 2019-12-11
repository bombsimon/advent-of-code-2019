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

const (
	black = "█"
	white = "░"
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
	Input         int
	Output        []int
	Pointer       int
	Base          int
	Sequence      []int
	Halted        bool
	PauseAtOutput bool
}

type direction int

const (
	directionUp = iota
	directionLeft
	directionDown
	directionRight
)

type robot struct {
	X         int
	Y         int
	Direction direction
	Seen      map[string]struct{}
	Grid      [][]string
	Computer  computer
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

	r1 := newRobot(sequence)
	r1.run(1)
	fmt.Println("number of panes painted at least once:", len(r1.Seen))

	r2 := newRobot(sequence)
	r2.run(2)
	r2.show()
}

func newRobot(sequence []int) *robot {
	var (
		size = 90
		grid = make([][]string, size)
	)

	r := robot{
		X:         75,
		Y:         40,
		Direction: directionUp,
		Seen:      map[string]struct{}{},
		Grid:      grid,
		Computer: computer{
			Sequence:      make([]int, len(sequence)),
			PauseAtOutput: true,
		},
	}

	copy(r.Computer.Sequence, sequence)

	for i := range r.Grid {
		r.Grid[i] = make([]string, size)

		for j := range r.Grid[i] {
			r.Grid[i][j] = "."
		}
	}

	return &r
}

func (r *robot) run(part int) {
	for {
		// Part two starts at white square
		if part == 2 && r.Computer.Pointer == 0 {
			r.Grid[r.X][r.Y] = white
		}

		// All is black by default, only change input if painted white.
		if r.Grid[r.X][r.Y] == white {
			r.Computer.Input = 1
		} else {
			r.Computer.Input = 0
		}

		// Fetch two codes
		r.Computer.process()
		r.Computer.process()

		if r.Computer.Halted {
			break
		}

		colorToDraw, directionToMove := r.Computer.Output[0], r.Computer.Output[1]

		r.draw(colorToDraw)
		r.turn(directionToMove)

		// Reset output
		r.Computer.Output = []int{}
	}
}

func (r *robot) turn(turnDirection int) {
	switch turnDirection {
	case 0:
		switch r.Direction {
		case directionUp:
			r.Y--
		case directionLeft:
			r.X++
		case directionDown:
			r.Y++
		case directionRight:
			r.X--
		}

		r.Direction++

		if r.Direction > 3 {
			r.Direction = directionUp
		}

	case 1:
		switch r.Direction {
		case directionUp:
			r.Y++
		case directionLeft:
			r.X--
		case directionDown:
			r.Y--
		case directionRight:
			r.X++
		}

		r.Direction--

		if r.Direction < 0 {
			r.Direction = directionRight
		}
	}
}

func (r *robot) draw(color int) {
	r.Seen[fmt.Sprintf("%d,%d", r.X, r.Y)] = struct{}{}

	switch color {
	case 0:
		r.Grid[r.X][r.Y] = black
	case 1:
		r.Grid[r.X][r.Y] = white
	}
}

func (r *robot) arrow() string {
	switch r.Direction {
	case directionUp:
		return "^"
	case directionLeft:
		return "<"
	case directionDown:
		return "v"
	case directionRight:
		return ">"
	}

	return "O"
}

func (r *robot) show() {
	fmt.Printf("current: %d,%d (%d)\n", r.X, r.Y, r.Direction)

	val := r.Grid[r.X][r.Y]

	r.Grid[r.X][r.Y] = r.arrow()

	for _, row := range r.Grid {
		fmt.Println(strings.Trim(strings.ReplaceAll(fmt.Sprint(row), " ", ""), "[]"))
	}

	r.Grid[r.X][r.Y] = val
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

		if c.PauseAtOutput {
			c.Pointer += jumpMap[opCodeOutput]
			return
		}

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
