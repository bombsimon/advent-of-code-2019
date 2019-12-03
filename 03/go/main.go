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

type coordinateMap map[coordinates]int

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
		fewestSteps      int
	)

	m1 := mark(r1)
	m2 := mark(r2)

	bothFunc := func(ma, mb coordinateMap) coordinateMap {
		both := coordinateMap{}

		for k := range ma {
			if _, ok := mb[k]; ok {
				both[k] = ma[k] + mb[k]
			}
		}

		return both
	}

	intersections := bothFunc(m1, m2)

	for c := range intersections {
		distances = append(distances, manhattanDistance(c))

		totalSteps := intersections[c]

		if fewestSteps == 0 || totalSteps < fewestSteps {
			fewestSteps = totalSteps
		}
	}

	for _, d := range distances {
		if shortestDistance == 0 || d < shortestDistance {
			shortestDistance = d
		}
	}

	fmt.Println("shortest distance", shortestDistance)
	fmt.Println("fewest steps", fewestSteps)
}

func mark(steps []string) map[coordinates]int {
	type updateCoordianteFunc func(x, y int) (int, int)

	var (
		x, y          = 0, 0
		totalSteps    = 0
		c             = coordinateMap{}
		updateFuncMap = map[string]updateCoordianteFunc{
			Up: func(x, y int) (int, int) {
				y--
				return x, y
			},
			Down: func(x, y int) (int, int) {
				y++
				return x, y
			},
			Left: func(x, y int) (int, int) {
				x--
				return x, y
			},
			Right: func(x, y int) (int, int) {
				x++
				return x, y
			},
		}
	)

	for _, step := range steps {
		direction, length := getDirectionAndLength(step)

		for range make([]struct{}, length) {
			totalSteps++

			x, y = updateFuncMap[direction](x, y)

			co := coordinates{x, y}
			if _, ok := c[co]; !ok {
				c[co] = 0
			}

			c[co] += totalSteps
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
