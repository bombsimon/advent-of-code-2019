package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
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
	tileEmpty = iota
	tileWall
	tileBlock
	tileHorizontalPaddle
	tileBall
)

const (
	joystickNeutral = 0
	joystickLeft    = -1
	joytickRight    = 1
)

var tileMap = map[int]string{
	tileEmpty:            "  ",
	tileWall:             "🧱",
	tileBlock:            " ◻️",
	tileHorizontalPaddle: "🏓",
	tileBall:             "🎾",
}

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
	PauseAtOutput bool
	Halted        bool
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

	run(sequence, 0)
}

func run(sequence []int, input int) {
	c := computer{Input: input, PauseAtOutput: true}

	c.Sequence = make([]int, len(sequence))
	copy(c.Sequence, sequence)

	var (
		width            = 40
		height           = 25
		row              = make([][]int, height)
		objects          = make(map[int]int)
		display          = 0
		ballPosition     = []int{0, 0}
		prevBallPosition = []int{0, 0}
		paddlePosition   = []int{0, 0}
	)

	// Create a grid for visualization
	for i := range row {
		row[i] = make([]int, width)
	}

	// Set quarters to 2
	c.Sequence[0] = 2

	showState := func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout

		if err := cmd.Run(); err != nil {
			panic(err)
		}

		for i := range row {
			for j := range row[i] {
				fmt.Print(tileMap[row[i][j]])
			}

			fmt.Println("")
		}

		time.Sleep(30 * time.Millisecond)
	}

	// Update the current computer/amplifier
	for {
		c.process()

		if c.Halted {
			break
		}

		if len(c.Output)%3 == 0 {
			x, y, objectID := c.Output[0], c.Output[1], c.Output[2]

			// Draw tiles if we're inbound.
			if x >= 0 && y >= 0 {
				row[y][x] = objectID
			}

			// Calculate number of objects
			objects[objectID]++

			// Update object positions if we're at a paddle or ball position.
			switch objectID {
			case tileHorizontalPaddle:
				paddlePosition = []int{x, y}
			case tileBall:
				ballPosition = []int{x, y}
			}

			// Decide how to move the paddle based on the relation to the ball.
			switch {
			case ballPosition[0] > paddlePosition[0]:
				c.Input = joytickRight
			case ballPosition[0] < paddlePosition[0]:
				c.Input = joystickLeft
			default:
				c.Input = joystickNeutral
			}

			// Set score when given instruction is shown.
			if x == -1 && y == 0 {
				display = objectID
			}

			if prevBallPosition[0] != ballPosition[0] || prevBallPosition[1] != ballPosition[1] {
				showState()

				prevBallPosition = []int{ballPosition[0], ballPosition[1]}
			}

			// Reset the output
			c.Output = []int{}
		}
	}

	fmt.Println("number of blocks", objects[tileBlock])
	fmt.Println("final score", display)
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
