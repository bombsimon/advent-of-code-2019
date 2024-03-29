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
	numberOfRows = 6
	numberOfCols = 25
	colorBlack   = 0
	colorWhite   = 1
)

// Layer represents a 2D grid with pixels.
type Layer [][]int

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing file as input")
	}

	line, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	var layers = createLayers(line)

	partOne(layers)
	partTwo(layers)
}

func partOne(layers []Layer) {
	var (
		fewestZeros          = 0
		fewestZeroLayerIndex = 0
	)

	for i, layer := range layers {
		zero := findOccurancesInlayerOf(layer, 0)

		if fewestZeros == 0 || zero < fewestZeros {
			fewestZeros = zero
			fewestZeroLayerIndex = i
		}
	}

	ones := findOccurancesInlayerOf(layers[fewestZeroLayerIndex], 1)
	twos := findOccurancesInlayerOf(layers[fewestZeroLayerIndex], 2)

	fmt.Println("part one: ", ones*twos)
}

func partTwo(layers []Layer) {
	var (
		layerToDraw = layers[0]
	)

	for _, layer := range layers {
		for i, row := range layer {
			for j, col := range row {
				// Never re-fill transparent
				if col != colorBlack && col != colorWhite {
					continue
				}

				// If there's no active color in the pixil, add it.
				if layerToDraw[i][j] != colorBlack && layerToDraw[i][j] != colorWhite {
					layerToDraw[i][j] = col
				}
			}
		}
	}

	fmt.Println("rendering the image:")
	printLayerLetter(layerToDraw)
}

func createLayers(line []byte) []Layer {
	var (
		layer      = Layer{}
		layers     = []Layer{}
		currentRow = []int{}
	)

	for i, v := range strings.TrimSpace(string(line)) {
		intVal, _ := strconv.Atoi(string(v))

		currentRow = append(currentRow, intVal)

		if (i+1)%numberOfCols == 0 {
			// Add the row to the current layer, reset the row.
			layer = append(layer, currentRow)
			currentRow = []int{}

			if len(layer)%numberOfRows == 0 {
				// Add the current layer to the list of layers ever 6 rows.
				layers = append(layers, layer)
				layer = Layer{}
			}
		}
	}

	return layers
}

func findOccurancesInlayerOf(layer Layer, n int) int {
	var count = 0

	for _, row := range layer {
		for _, col := range row {
			if col == n {
				count++
			}
		}
	}

	return count
}

func printLayer(layer Layer) {
	for _, row := range layer {
		fmt.Println(row)
	}
}

func printLayerLetter(layer Layer) {
	var (
		black       = "█"
		white       = "░"
		transparent = " "
	)

	for _, row := range layer {
		for _, col := range row {
			switch col {
			case colorBlack:
				fmt.Print(black)
			case colorWhite:
				fmt.Print(white)
			default:
				fmt.Print(transparent)
			}
		}

		fmt.Println("")
	}
}
