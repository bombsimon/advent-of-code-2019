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
		highestOutput = 0
		highestInput  = []int{}
	)

	for _, perm := range permutations([]int{0, 1, 2, 3, 4}) {
		var (
			a = perm[0]
			b = perm[1]
			c = perm[2]
			d = perm[3]
			e = perm[4]
		)

		inputs := [][]int{
			{a, 0},
			{b, 0},
			{c, 0},
			{d, 0},
			{e, 0},
		}

		lastOutput := 0
		for i, input := range inputs {
			lastOutput = run(sequence, input)

			if i < len(inputs)-1 {
				inputs[i+1][1] = lastOutput
			}
		}

		if lastOutput > highestOutput {
			highestOutput = lastOutput
			highestInput = perm
		}
	}

	fmt.Printf("highest thurst '%d' met with '%v'\n", highestOutput, highestInput)
}

func run(originalSequence, input []int) int {
	sequence := make([]int, len(originalSequence))
	copy(sequence, originalSequence)

	var (
		startPointer = 0
		output       = 0
	)

	return process(sequence, startPointer, input, output)
}

func process(sequence []int, ptr int, input []int, output int) int {
	opCode := sequence[ptr]

	// Halt code found, stop processing.
	if opCode == opCodeHalt {
		return output
	}

	// Pointer is outside of list range.
	if ptr > len(sequence) {
		return output
	}

	// Get firwst two arguments, parse their input code and the actual op code
	// for larger numbers.
	first, second, pos := sequence[ptr+1], sequence[ptr+2], sequence[ptr+3]
	firstMode, secondMode, opCode := parseOpCode(opCode)

	// Convert immediate to positional if within boundaries.
	if firstMode == paramModePosition && len(sequence) >= first && opCode != opCodeStore {
		first = sequence[first]
	}

	if secondMode == paramModePosition && len(sequence) >= second {
		second = sequence[second]
	}

	switch opCode {
	case opCodeAdd:
		sequence[pos] = first + second

		return process(sequence, ptr+4, input, output)

	case opCodeMultiply:
		sequence[pos] = first * second

		return process(sequence, ptr+4, input, output)

	case opCodeStore:
		// Pick first and shift it.
		in := input[0]

		if len(input) > 1 {
			input = input[1:]
		}

		sequence[first] = in

		return process(sequence, ptr+2, input, output)

	case opCodeOutput:
		return process(sequence, ptr+2, input, first)

	case opCodeJumpIfTrue:
		if first != 0 {
			return process(sequence, second, input, output)
		}

		return process(sequence, ptr+3, input, output)

	case opCodeJumpIfFalse:
		if first == 0 {
			return process(sequence, second, input, output)
		}

		return process(sequence, ptr+3, input, output)

	case opCodeLessThan:
		if first < second {
			sequence[pos] = 1
		} else {
			sequence[pos] = 0
		}

		return process(sequence, ptr+4, input, output)

	case opCodeEquals:
		if first == second {
			sequence[pos] = 1
		} else {
			sequence[pos] = 0
		}

		return process(sequence, ptr+4, input, output)

	default:
		panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", ptr, opCode))
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
