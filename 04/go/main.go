package main

import (
	"fmt"
	"math"
)

const (
	codeStart = 246515
	codeStop  = 739105
)

func main() {
	var (
		combinations          = 0
		combinationsNoTriplet = 0
	)

	for i := codeStart; i <= codeStop; i++ {
		hasDoubles, hasDoublesNoTriplets := increaseAndHasPair(i)

		if hasDoubles {
			combinations++
		}

		if hasDoublesNoTriplets {
			combinationsNoTriplet++
		}
	}

	fmt.Println("total combinations with pair:", combinations)
	fmt.Println("total combinations with pair not in a triplet", combinationsNoTriplet)
}

func increaseAndHasPair(number int) (bool, bool) {
	var (
		exludedPos = map[int]struct{}{}
		pairPos    = map[int]struct{}{}
	)

	for i := countDigits(number); i > 1; i-- {
		first, second := digitAtPosition(number, i), digitAtPosition(number, i-1)

		// The digits decreases so this is no candidate.
		if second < first {
			return false, false
		}

		// If we've got more than 2 digits left, extract the third for
		// comparison.
		third := -1
		if i > 2 {
			third = digitAtPosition(number, i-2)
		}

		// Three digits in a row are the same, store their positions.
		if first == second && second == third {
			exludedPos[i] = struct{}{}
			exludedPos[i-1] = struct{}{}
			exludedPos[i-2] = struct{}{}
		}

		// Two digits in a row are the same, store their position.
		if first == second {
			pairPos[i] = struct{}{}
			pairPos[i-1] = struct{}{}
		}
	}

	// We had at least one double if there are any pair posisions.
	hasDoubles := len(pairPos) > 0

	// Check the positions of digits in a triplet, remove them from the pair
	// positions.
	for k := range pairPos {
		if _, ok := exludedPos[k]; ok {
			delete(pairPos, k)
		}
	}

	return hasDoubles, len(pairPos) > 0
}

// Get a digit from a specific position. Reference:
// https://stackoverflow.com/questions/46753736/extract-digits-at-a-certain-position-in-a-number
func digitAtPosition(num, place int) int {
	r := num % int(math.Pow(10, float64(place)))

	return r / int(math.Pow(10, float64(place-1)))
}

// Count the number of digits. Not really needed with only fixed width numbers
// but useful for testing/TDD.
func countDigits(i int) int {
	count := 0

	for i != 0 {
		i /= 10
		count++
	}

	return count
}
