package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing file as input")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	var (
		fileReader   = bufio.NewReader(file)
		totalPartOne = float64(0)
		totalPartTwo = float64(0)
	)

	for {
		line, _, err := fileReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatalf("error reading file: %s", err.Error())
		}

		fVal, err := strconv.ParseFloat(string(line), 64)
		if err != nil {
			log.Fatalf("could not convert line to int: %s", err.Error())
		}

		totalPartOne += processLinePartOne(fVal)
		totalPartTwo += processLinePartTwo(fVal)
	}

	fmt.Printf("part one: %.0f\n", totalPartOne)
	fmt.Printf("part two: %.0f\n", totalPartTwo)
}

func processLinePartOne(fVal float64) float64 {
	return math.Floor(fVal/3) - 2
}

func processLinePartTwo(fVal float64) float64 {
	totalFuel := float64(0)

	for {
		fVal = processLinePartOne(fVal)
		if fVal <= 0 {
			break
		}

		totalFuel += fVal
	}

	return totalFuel
}
