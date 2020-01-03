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
	Input         int
	Output        []int
	Pointer       int
	Base          int
	Sequence      []int
	Halted        bool
	PauseAtOutput bool
}

const (
	gridSize = 50
)

const (
	wall      = "🚧"
	path      = "◾️"
	unknown   = "🧱"
	oxygen    = "💧"
	robotChar = "🤖"
)

type coordinate struct {
	X int
	Y int
}

type direction int

const (
	directionNorth = iota + 1
	directionSouth
	directionWest
	directionEast
)

type robot struct {
	X             int
	Y             int
	Direction     direction
	Grid          map[coordinate]string
	Computer      computer
	oxygenPos     coordinate
	oxygenStep    int
	oxygenMaxStep int
}

func (d direction) String() string {
	switch d {
	case directionNorth:
		return "north"
	case directionSouth:
		return "south"
	case directionWest:
		return "west"
	case directionEast:
		return "east"
	}

	return "unknown"
}

func newRobot(sequence []int) *robot {
	r := robot{
		X:    gridSize / 2,
		Y:    gridSize / 2,
		Grid: map[coordinate]string{},
		Computer: computer{
			Sequence:      make([]int, len(sequence)),
			PauseAtOutput: true,
			Input:         directionNorth,
		},
	}

	copy(r.Computer.Sequence, sequence)

	return &r
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
	r1.checkNext(sequence, 1)

	r1.show()

	// Put the robot at the oxygen position after completing the first part.
	r1.X, r1.Y = r1.oxygenPos.X, r1.oxygenPos.Y
	r1.oxygenTime(make(map[coordinate]struct{}), 1)

	fmt.Println("part 1: found oxygen after", r1.oxygenStep, "steps")
	fmt.Println("part 2: time to fill", r1.oxygenMaxStep)
}

func (r *robot) checkNext(originalSequence []int, stepsFromStart int) {
	var (
		x   = r.X
		y   = r.Y
		seq = make([]int, len(originalSequence))
	)

	// Copy the original sequence passed to this function to sequence.
	copy(seq, originalSequence)

	// Mark the current position as a path.
	r.Grid[coordinate{X: x, Y: y}] = path

	for _, dir := range []direction{directionNorth, directionSouth, directionWest, directionEast} {
		// Copy the sequence passed so we can reset it for each direction.
		copy(r.Computer.Sequence, seq)

		// Reset X and Y for each direction based on where we started.
		r.X = x
		r.Y = y

		// Set the new direction (just for info) and set the direction as input
		// to the computer.
		r.Direction = dir
		r.Computer.Input = int(dir)

		// Get the new coordinates based of the direction we're looking. If
		// we've already been in that direction, move on.
		nextCoordinates := r.nextCoordiantes()
		if _, ok := r.Grid[nextCoordinates]; ok {
			continue
		}

		// Run one clock in the computer.
		r.Computer.process()

		// Fetch the output from the process.
		moveResult := r.Computer.Output[0]
		r.Computer.Output = []int{}

		switch moveResult {
		case 0:
			r.Grid[nextCoordinates] = wall

		case 1:
			r.Grid[nextCoordinates] = path
			r.X = nextCoordinates.X
			r.Y = nextCoordinates.Y

			// Keep going this direction
			r.checkNext(r.Computer.Sequence, stepsFromStart+1)

		case 2:
			r.oxygenPos = nextCoordinates
			r.oxygenStep = stepsFromStart
			r.Grid[nextCoordinates] = oxygen
		}
	}
}

func (r *robot) oxygenTime(seen map[coordinate]struct{}, maxSteps int) {
	x, y := r.X, r.Y

	for _, dir := range []direction{directionNorth, directionSouth, directionWest, directionEast} {
		r.Direction = dir
		r.X, r.Y = x, y

		nextCoordinates := r.nextCoordiantes()
		if _, ok := seen[nextCoordinates]; ok {
			continue
		}

		typ, ok := r.Grid[nextCoordinates]
		if !ok {
			continue
		}

		if typ != path {
			continue
		}

		if maxSteps > r.oxygenMaxStep {
			r.oxygenMaxStep = maxSteps
		}

		seen[nextCoordinates] = struct{}{}
		r.X = nextCoordinates.X
		r.Y = nextCoordinates.Y

		r.oxygenTime(seen, maxSteps+1)
	}
}

func (r *robot) show() {
	for x := range make([]struct{}, gridSize) {
		for y := range make([]struct{}, gridSize) {
			c := coordinate{X: x, Y: y}
			p := unknown

			if v, ok := r.Grid[c]; ok {
				p = v
			}

			if x == r.X && y == r.Y {
				p = robotChar
			}

			fmt.Print(p)
		}

		fmt.Println("")
	}
}

func (r *robot) nextCoordiantes() coordinate {
	x, y := r.X, r.Y

	switch r.Direction {
	case directionNorth:
		x--
	case directionSouth:
		x++
	case directionWest:
		y--
	case directionEast:
		y++
	}

	return coordinate{X: x, Y: y}
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
