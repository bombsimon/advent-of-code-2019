package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type coordinates struct {
	x int
	y int
}

// Directions
const (
	Up    = "U"
	Down  = "D"
	Left  = "L"
	Right = "R"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing file as input")
	}

	fileContent, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	var (
		fileLines        = strings.Split(string(fileContent), "\n")
		r1               = strings.Split(fileLines[0], ",")
		r2               = strings.Split(fileLines[1], ",")
		distances        []int
		shortestDistance int
	)

	m1 := mark(r1)
	m2 := mark(r2)

	bothFunc := func(ma, mb map[coordinates]struct{}) map[coordinates]struct{} {
		both := map[coordinates]struct{}{}

		for k := range ma {
			if _, ok := mb[k]; ok {
				both[k] = struct{}{}
			}
		}

		return both
	}

	intersections := bothFunc(m1, m2)

	for c := range intersections {
		distances = append(distances, manhattanDistance(c))
	}

	for _, d := range distances {
		if shortestDistance == 0 || d < shortestDistance {
			shortestDistance = d
		}
	}

	fmt.Println("shortest distance", shortestDistance)
}

func mark(steps []string) map[coordinates]struct{} {
	var (
		x, y = 0, 0
		c    = map[coordinates]struct{}{}
	)

	for _, step := range steps {
		direction, length := getDirectionAndLength(step)

		switch direction {
		case Up:
			for range make([]struct{}, length) {
				y--
				c[coordinates{x, y}] = struct{}{}
			}
		case Down:
			for range make([]struct{}, length) {
				y++
				c[coordinates{x, y}] = struct{}{}
			}
		case Left:
			for range make([]struct{}, length) {
				x--
				c[coordinates{x, y}] = struct{}{}
			}
		case Right:
			for range make([]struct{}, length) {
				x++
				c[coordinates{x, y}] = struct{}{}
			}
		default:
			panic("nope")
		}
	}

	return c
}

func getDirectionAndLength(step string) (string, int) {
	var (
		direction = string(step[0])
		strLength = step[1:]
		length, _ = strconv.Atoi(strLength)
	)

	return direction, length
}

func manhattanDistance(c coordinates) int {
	return int(math.Abs(0-float64(c.x)) + math.Abs(0-float64(c.y)))
}
