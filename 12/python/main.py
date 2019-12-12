#!/usr/bin/env python3

import sys
import math


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

    step(positions, 1000)
    step(positions)


def step(positions, times=None):
    """
    Apply the velocity to each moon the number of times given
    """

    # r is the cu*r*rrent state of the positions after possible gravity applied.
    r = {"x": [], "y": [], "z": []}

    # To zero is the number of steps it took for the velocity to reach zero
    # again.
    to_zero = {"x": 0, "y": 0, "z": 0}

    # Populate the initial state of the positions given. This will be set as the
    # first index. The second index is used as velocity and starts at 0.
    for pos in positions:
        r["x"].append([pos["x"], 0])
        r["y"].append([pos["y"], 0])
        r["z"].append([pos["z"], 0])

    # The number of times the iteration has been made.
    for c in ["x", "y", "z"]:
        # Reset step count for each coordiante.
        steps = 0

        while True:
            steps += 1

            # Use each position as base.
            for i, _ in enumerate(r[c]):
                # Compare to all other positions.
                for j, _ in enumerate(r[c]):
                    # Exclude ourself when we compare.
                    if i == j:
                        continue

                    # If the moon has lower position than the other moon (j),
                    # increase velocity.
                    if r[c][i][0] < r[c][j][0]:
                        r[c][i][1] += 1

                    # If the moon has higher position than the other moon (j),
                    # decrease velocity.
                    if r[c][i][0] > r[c][j][0]:
                        r[c][i][1] -= 1

            # Apply the velocity for each moon and coordinate.
            all_zero = True
            for i, _ in enumerate(r[c]):
                r[c][i][0] += r[c][i][1]

                # If any of the n moons didn't reach zero we're not back at the
                # initial state.
                if r[c][i][1] != 0:
                    all_zero = False

            # If all are zero we're back at the initial state. This only means
            # that the velocity is back to zero. This means that half the
            # process is done. We can use this value to calculate the iterations
            # by multiplying the least common multiple of the coordinates with
            # two.
            if all_zero:
                to_zero[c] = steps

                # Stop iteration if the number of times to iterage wasn't set.
                if times is None:
                    break

            # Stop iteration if we've iterated as many times as passed.
            if times is not None:
                if steps == times:
                    break

    # Part one
    if times is not None:
        total = 0

        # Assume we have four moons to sum.
        for i in range(4):
            pot = 0
            kin = 0

            # Sum each axis for each moon
            for axis in ["x", "y", "z"]:
                pot += abs(r[axis][i][0])
                kin += abs(r[axis][i][1])

            # Multiply the moons potential and kinetic energy.
            total += pot * kin

        print("The sum of the total energy is {}".format(total))

        return

    # Part two
    # https://stackoverflow.com/a/42517664/2274551
    x2, y2, z2 = to_zero["x"], to_zero["y"], to_zero["z"]

    gcd = math.gcd(x2, y2)
    lcm2 = y2 * z2 // gcd
    lcm3 = x2 * lcm2 // math.gcd(x2, lcm2)

    # Multiply by two to get the final number of iterations needed.
    print("Number of iterations needed: {}".format(lcm3 * 2))


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
