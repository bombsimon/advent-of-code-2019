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
	opCodeIgnore      = 99
	sequenceSkipSteps = 3
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

	for i := 0; i < len(sequence); i++ {
		if sequence[i] == opCodeIgnore {
			break
		}

		if i >= len(sequence)-sequenceSkipSteps {
			continue
		}

		var (
			opCode = sequence[i]
			n1Val  = sequence[i+1]
			n2Val  = sequence[i+2]
			pos    = sequence[i+3]
		)

		paramModeOne, paramModeTwo, _, opCode := parseOpCode(opCode)

		switch opCode {
		case opCodeStore, opCodeOutput:
			// Always immediate for store and output.
		default:
			if paramModeOne == paramModePosition {
				n1Val = sequence[n1Val]
			}

			if paramModeTwo == paramModePosition {
				n2Val = sequence[n2Val]
			}
		}

		switch opCode {
		case opCodeAdd:
			sequence[pos] = n1Val + n2Val
			i += sequenceSkipSteps
		case opCodeMultiply:
			sequence[pos] = n1Val * n2Val
			i += sequenceSkipSteps
		case opCodeStore:
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("input: ")

			text, err := reader.ReadString('\n')
			if err != nil {
				panic(err)
			}

			code, err := strconv.Atoi(strings.TrimSpace(text))
			if err != nil {
				panic(err)
			}

			sequence[n1Val] = code
			i++
		case opCodeOutput:
			fmt.Println(sequence[n1Val])
			i++
		case opCodeJumpIfTrue:
			if n1Val != 0 {
				i = n2Val - 1
			} else {
				// Skip to next instruction
				i += 2
			}
		case opCodeJumpIfFalse:
			if n1Val == 0 {
				i = n2Val - 1
			} else {
				// Skip to next instruction
				i += 2
			}
		case opCodeLessThan:
			if n1Val < n2Val {
				sequence[pos] = 1
			} else {
				sequence[pos] = 0
			}

			i += sequenceSkipSteps
		case opCodeEquals:
			if n1Val == n2Val {
				sequence[pos] = 1
			} else {
				sequence[pos] = 0
			}

			i += sequenceSkipSteps
		case opCodeIgnore:
			break
		default:
			panic(fmt.Sprintf("unknown instruction at position %d, opCode: %d", i, opCode))
		}
	}
}

func parseOpCode(opCode int) (paramMode, paramMode, paramMode, int) {
	codeAsString := fmt.Sprintf("%05d", opCode)

	code, err := strconv.Atoi(codeAsString[4:])
	if err != nil {
		panic(err)
	}

	modeArgOne, err := strconv.Atoi(string(codeAsString[2]))
	if err != nil {
		panic(err)
	}

	modeArgTwo, err := strconv.Atoi(string(codeAsString[1]))
	if err != nil {
		panic(err)
	}

	modeArgThree, err := strconv.Atoi(string(codeAsString[0]))
	if err != nil {
		panic(err)
	}

	return paramMode(modeArgOne), paramMode(modeArgTwo), paramMode(modeArgThree), code
}
