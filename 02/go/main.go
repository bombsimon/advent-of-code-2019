package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	opCodeAdd         = 1
	opCodeMultiply    = 2
	opCodeIgnore      = 99
	sequenceSkipSteps = 3
	partTwoMaxValue   = 99
	partTwoOutput     = 19690720
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

	sequence[1] = 12
	sequence[2] = 2

	fmt.Println("part one:", partOne(sequence, 12, 2))
	fmt.Println("part two:", partTwo(sequence))
}

func partTwo(originalSequence []int) int {
	for i := range make([]struct{}, partTwoMaxValue) {
		for j := range make([]struct{}, partTwoMaxValue) {
			result := partOne(originalSequence, i, j)
			if result == partTwoOutput {
				v, _ := strconv.Atoi(fmt.Sprintf("%02d%02d", i, j))
				return v
			}
		}
	}

	return 0
}

func partOne(originalSequence []int, noun, verb int) int {
	sequence := make([]int, len(originalSequence))
	copy(sequence, originalSequence)

	sequence[1] = noun
	sequence[2] = verb

	for i := 0; i < len(sequence); i++ {
		if sequence[i] == opCodeIgnore {
			break
		}

		if i >= len(sequence)-sequenceSkipSteps {
			continue
		}

		var (
			opCode = sequence[i]
			n1Pos  = sequence[i+1]
			n2Pos  = sequence[i+2]
			pos    = sequence[i+3]
			n1     = sequence[n1Pos]
			n2     = sequence[n2Pos]
		)

		// Invalid sequence.
		if pos >= len(sequence) {
			return 0
		}

		switch opCode {
		case opCodeAdd:
			sequence[pos] = n1 + n2
			i += sequenceSkipSteps
		case opCodeMultiply:
			sequence[pos] = n1 * n2
			i += sequenceSkipSteps
		case opCodeIgnore:
			continue
		default:
			panic("unknown instruction")
		}
	}

	return sequence[0]
}
