package main

import (
	"bufio"
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

	partOne(sequence)
}

func partOne(originalSequence []int) {
	sequence := make([]int, len(originalSequence))
	copy(sequence, originalSequence)

	process(sequence, 0)
}

func process(sequence []int, ptr int) {
	opCode := sequence[ptr]

	// Halt code found, stop processing.
	if opCode == opCodeHalt {
		return
	}

	// Pointer is outside of list range.
	if ptr > len(sequence) {
		return
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
		process(sequence, ptr+4)
	case opCodeMultiply:
		sequence[pos] = first * second
		process(sequence, ptr+4)
	case opCodeStore:
		sequence[first] = readInput()
		process(sequence, ptr+2)
	case opCodeOutput:
		fmt.Println(first)
		process(sequence, ptr+2)
	case opCodeJumpIfTrue:
		if first != 0 {
			process(sequence, second)
		} else {
			process(sequence, ptr+3)
		}
	case opCodeJumpIfFalse:
		if first == 0 {
			process(sequence, second)
		} else {
			process(sequence, ptr+3)
		}
	case opCodeLessThan:
		if first < second {
			sequence[pos] = 1
		} else {
			sequence[pos] = 0
		}

		process(sequence, ptr+4)
	case opCodeEquals:
		if first == second {
			sequence[pos] = 1
		} else {
			sequence[pos] = 0
		}

		process(sequence, ptr+4)
	default:
		panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", ptr, opCode))
	}
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

func readInput() int {
	var readCode int

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("input: ")

		text, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		code, err := strconv.Atoi(strings.TrimSpace(text))
		if err != nil {
			continue
		}

		readCode = code

		break
	}

	return readCode
}
