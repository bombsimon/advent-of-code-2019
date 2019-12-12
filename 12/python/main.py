#!/usr/bin/env python3

import sys
import time
from pprint import pprint


def main():
    """
    Main runs the astroid lookup program
    """

    positions = []

    with open(sys.argv[1]) as f:
        for line in f:
            line = line[1:-2].split(" ")
            positions.append(
                {
                    x.split("=")[0]: int(x.split("=")[1].strip(","))
                    for x in line
                }
            )

    apply(positions)


def apply(positions):
    debug = False
    steps = 1000
    velocity = [{"x": 0, "y": 0, "z": 0} for _ in range(len(positions))]

    for _ in range(1000):
        for i in range(len(positions)):
            for comparison in range(len(positions)):
                if i == comparison:
                    continue

                for axis in ["x", "y", "z"]:
                    if positions[i][axis] < positions[comparison][axis]:
                        velocity[i][axis] += 1
                    elif positions[i][axis] > positions[comparison][axis]:
                        velocity[i][axis] -= 1

        if debug:
            print("velocity after cycle")
            pprint(velocity)

        for i, _ in enumerate(velocity):
            for axis in ["x", "y", "z"]:
                positions[i][axis] += velocity[i][axis]

        if debug:
            print("position after velocity")
            show(positions)

    total_sum = 0
    for i, _ in enumerate(positions):
        pot = 0
        kin = 0
        moon_sum = 0

        for axis in ["x", "y", "z"]:
            pot += abs(positions[i][axis])
            kin += abs(velocity[i][axis])
            moon_sum = pot * kin

            positions[i]["energy"] = moon_sum

        total_sum += positions[i]["energy"]

    # NOT: 2151
    pprint(positions)
    print(total_sum)


def show(positions):
    for moon in positions:
        print(
            "<x={:<4d} y={:<4d} z={:d}>".format(
                moon["x"], moon["y"], moon["z"]
            )
        )

    print()


if __name__ == "__main__":
    main()
