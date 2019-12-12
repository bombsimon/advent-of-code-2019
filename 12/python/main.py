#!/usr/bin/env python3

import sys
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

    part_one(positions)


def part_one(positions, times=1000, debug=False):
    """
    Apply the velocity to each moon the number of times given
    """
    velocity = [{"x": 0, "y": 0, "z": 0} for _ in range(len(positions))]

    for _ in range(times):
        for i, _ in enumerate(positions):
            for comparison, _ in enumerate(positions):
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

    show(positions, velocity)

    print(total_sum)


def show(positions, velocity=None):
    """
    Show the state of passed positions
    """
    for i, _ in enumerate(positions):
        vel = "<unknown>"

        if velocity is not None:
            vel = "<x={:<4d} y={:<4d} z={:d}>".format(
                velocity[i]["x"], velocity[i]["y"], velocity[i]["z"]
            )

        print(
            "<x={:<4d} y={:<4d} z={:d}> | {}".format(
                positions[i]["x"], positions[i]["y"], positions[i]["z"], vel
            )
        )

    print()


if __name__ == "__main__":
    main()
