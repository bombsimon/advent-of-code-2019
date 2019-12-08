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

const (
	paramModePosition paramMode = iota
	paramModeImmediate
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
	opCodeHalt        = 99
)

type computer struct {
	Name       string
	Phase      int
	Input      int
	InputCount int
	Value      int
	Pointer    int
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

	var (
		permutationsPartOne = permutations([]int{0, 1, 2, 3, 4})
		permutationsPartTwo = permutations([]int{5, 6, 7, 8, 9})
	)

	run(sequence, permutationsPartOne)
	run(sequence, permutationsPartTwo)
}

func run(sequence []int, permutations [][]int) {
	var (
		highestOutput = 0
		highestInput  = []int{}
	)

	for _, perm := range permutations {
		inputs := []computer{
			{Name: "A", Phase: perm[0]},
			{Name: "B", Phase: perm[1]},
			{Name: "C", Phase: perm[2]},
			{Name: "D", Phase: perm[3]},
			{Name: "E", Phase: perm[4]},
		}

		lastOutput := 0

		for i := 0; i < len(inputs); i++ {
			c := inputs[i]

			// Ensure we have a sequence for the permutation.
			if c.Sequence == nil {
				c.Sequence = make([]int, len(sequence))
				copy(c.Sequence, sequence)
			}

			// Update the current computer/amplifier
			c.process()

			// Store the updated computer in index i with it's poinuter and
			// values.
			inputs[i] = c

			// Set it's current value when stopped to next input.
			lastOutput = c.Value

			// Set last output to input for next amplifier, or if we're the last
			// amplifier set it to the first one.
			if i < len(inputs)-1 {
				inputs[i+1].Input = lastOutput
			} else {
				inputs[0].Input = lastOutput
			}

			// If we're the last amplifier but not yet halted, start over.
			if i == len(inputs)-1 && !c.Halted {
				i = -1
			}
		}

		if lastOutput > highestOutput {
			highestOutput = lastOutput
			highestInput = perm
		}
	}

	fmt.Printf("highest thurst '%d' met with '%v'\n", highestOutput, highestInput)
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

	// Get firwst two arguments, parse their input code and the actual op code
	// for larger numbers.
	first, second := c.Sequence[c.Pointer+1], c.Sequence[c.Pointer+2]
	firstMode, secondMode, opCode := parseOpCode(opCode)

	// Store the position (third argument) if there's enough positions left in
	// memory.
	pos := 0
	if len(c.Sequence) > c.Pointer+3 {
		pos = c.Sequence[c.Pointer+3]
	}

	// Convert immediate to positional if within boundaries.
	if firstMode == paramModePosition && len(c.Sequence) >= first && opCode != opCodeStore {
		first = c.Sequence[first]
	}

	if secondMode == paramModePosition && len(c.Sequence) >= second {
		second = c.Sequence[second]
	}

	switch opCode {
	case opCodeAdd:
		c.Sequence[pos] = first + second
		c.Pointer += 4

		c.process()

	case opCodeMultiply:
		c.Sequence[pos] = first * second
		c.Pointer += 4

		c.process()

	case opCodeStore:
		input := c.Phase
		if c.InputCount > 0 {
			input = c.Input
		}

		c.Sequence[first] = input
		c.Pointer += 2
		c.InputCount++

		c.process()

	case opCodeOutput:
		c.Value = first
		c.Pointer += 2

		// Puase process if opcode > 4
		if c.Phase > 4 {
			return
		}

		c.process()

	case opCodeJumpIfTrue:
		if first != 0 {
			c.Pointer = second
			c.process()
		} else {
			c.Pointer += 3
			c.process()
		}

	case opCodeJumpIfFalse:
		if first == 0 {
			c.Pointer = second
			c.process()
		} else {
			c.Pointer += 3
			c.process()
		}

	case opCodeLessThan:
		if first < second {
			c.Sequence[pos] = 1
		} else {
			c.Sequence[pos] = 0
		}

		c.Pointer += 4
		c.process()

	case opCodeEquals:
		if first == second {
			c.Sequence[pos] = 1
		} else {
			c.Sequence[pos] = 0
		}

		c.Pointer += 4
		c.process()

	default:
		panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", c.Pointer, opCode))
	}
}

// https://stackoverflow.com/questions/30226438/generate-all-permutations-in-go
func permutations(arr []int) [][]int {
	var f func([]int, int)

	result := [][]int{}

	f = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)

			result = append(result, tmp)

			return
		}

		for i := 0; i < n; i++ {
			f(arr, n-1)

			if n%2 == 1 {
				arr[i], arr[n-1] = arr[n-1], arr[i]
			} else {
				arr[0], arr[n-1] = arr[n-1], arr[0]
			}
		}
	}

	f(arr, len(arr))

	return result
}

func parseOpCode(opCode int) (paramMode, paramMode, int) {
	var (
		code          = opCode % 100
		modePositions = opCode / 100
		firstMode     = modePositions % 10
		secondMode    = modePositions / 10
	)

	return paramMode(firstMode), paramMode(secondMode), code
}
