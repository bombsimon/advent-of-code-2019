package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func main() {
	orbits := readFile()

	f, b := parseMap(orbits)
	d := map[string]int{}

	findDistanceFrom("COM", 0, f, d)

	totalOrbits := 0
	for _, distance := range d {
		totalOrbits += distance
	}

	// Part 1 solution
	fmt.Println("Total orbis", totalOrbits)

	var (
		youBackrefs           = backrefsFor("YOU", b, map[string]struct{}{})
		sanBackrefs           = backrefsFor("SAN", b, map[string]struct{}{})
		commonKeys            = []string{}
		nearestSharedDistance = 0
	)

	// Check common keys found as backrefs from both 'YOU' and 'SAN'.
	for k := range youBackrefs {
		if _, ok := sanBackrefs[k]; ok {
			commonKeys = append(commonKeys, k)
		}
	}

	// Iterate over the common ones and find the distance to the one closest to
	// you/farthest 'COM'.
	for _, k := range commonKeys {
		if nearestSharedDistance == 0 || d[k] > nearestSharedDistance {
			nearestSharedDistance = d[k]
		}
	}

	// Add the distance between 'YOU' and the closest shared + the distance
	// between 'SAN' and the closest shared one to get a delta.
	srcMinusNearest := len(youBackrefs) - nearestSharedDistance
	dstMinusNearest := len(sanBackrefs) - nearestSharedDistance
	delta := srcMinusNearest + dstMinusNearest - 2 // Remove 'YOU' and 'SAN'.

	// Part 2 solution
	fmt.Println("Minimum orbits to move between 'YOU' and 'SAN':", delta)
}

func backrefsFor(planet string, backrefs map[string]string, m map[string]struct{}) map[string]struct{} {
	backref, ok := backrefs[planet]
	if !ok {
		return m
	}

	m[backref] = struct{}{}

	return backrefsFor(backref, backrefs, m)
}

func findDistanceFrom(planet string, distance int, orbits map[string][]string, store map[string]int) {
	store[planet] = distance
	nextDistance := distance + 1

	f, ok := orbits[planet]
	if !ok {
		return
	}

	for _, p := range f {
		findDistanceFrom(p, nextDistance, orbits, store)
	}
}

func parseMap(orbits []string) (map[string][]string, map[string]string) {
	var (
		forward  = map[string][]string{}
		backward = map[string]string{}
	)

	for _, line := range orbits {
		parts := strings.Split(line, ")")
		lhs, rhs := parts[0], parts[1]

		forward[lhs] = append(forward[lhs], rhs)
		backward[rhs] = lhs
	}

	return forward, backward
}

func readFile() []string {
	if len(os.Args) < 2 {
		log.Fatal("missing file as input")
	}

	line, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("could not read file: %s", err.Error())
	}

	return strings.Split(strings.TrimSpace(string(line)), "\n")
}
