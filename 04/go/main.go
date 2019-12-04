package main

import (
	"fmt"
	"math"
)

const (
	// TRIED: 491639, 1058, 1049
	codeStart = 246515
	codeStop  = 739105
	// codeStart = 20
	// codeStop  = 40
)

func main() {
	combinations := 0

	for i := codeStart; i <= codeStop; i++ {
		if increaseAndHasPair(i) {
			combinations++
		}
	}

	fmt.Println("total combinations:", combinations)
}

func increaseAndHasPair(number int) bool {
	var hasPair bool

	for i := countDigits(number); i > 1; i-- {
		first, second := digitAtPosition(number, i), digitAtPosition(number, i-1)

		if first == second {
			hasPair = true
		}

		if second < first {
			return false
		}
	}

	return hasPair
}

// https://stackoverflow.com/questions/46753736/extract-digits-at-a-certain-position-in-a-number
func digitAtPosition(num, place int) int {
	r := num % int(math.Pow(10, float64(place)))

	return r / int(math.Pow(10, float64(place-1)))
}

func countDigits(i int) int {
	count := 0

	for i != 0 {
		i /= 10
		count++
	}

	return count
}
