#!/usr/bin/env python3
"""
Advent of Code 2019 - Day 3
"""

import sys


def main():
    """
    Main will read the input file and extract the two first rows.
    """
    with open(sys.argv[1]) as input_file:
        lines = input_file.read().splitlines()

    row_one, row_two = lines[0].split(","), lines[1].split(",")

    coordinates_one = mark(row_one)
    coordinates_two = mark(row_two)

    # Common keys:
    # https://stackoverflow.com/questions/1317410/finding-matching-keys-in-two-large-dictionaries-and-doing-it-fast
    common_keys = set(coordinates_one).intersection(set(coordinates_two))

    all_distances = [manhattan_distance((0, 0), ck) for ck in common_keys]

    print("shortest distance from 0,0: {}".format(min(all_distances)))

    # Get the number of steps it took to reach each point in one and two, sort
    # the keys and add them together to see which one has the lowest total.
    steps_one = [coordinates_one[k] for k in sorted(common_keys)]
    steps_two = [coordinates_two[k] for k in sorted(common_keys)]

    print(
        "the shortest total intersection is {}".format(
            min([steps_one[x] + steps_two[x] for x in range(len(steps_one))])
        )
    )


def mark(steps):
    """
    Mark all the coordinates extracted byt traversing the steps in each
    direction.
    """
    x_coordiante, y_coordiante = 0, 0
    coordinates = {}
    total_steps = 0

    def add_y(x, y):
        return x, y + 1

    def dec_y(x, y):
        return x, y - 1

    def add_x(x, y):
        return x + 1, y

    def dec_x(x, y):
        return x - 1, y

    actions = {"U": dec_y, "D": add_y, "L": dec_x, "R": add_x}

    for step in steps:
        direction, length = get_diration_and_length(step)

        for _ in range(length):
            total_steps += 1
            x_coordiante, y_coordiante = actions[direction](
                x_coordiante, y_coordiante
            )

            if (x_coordiante, y_coordiante) not in coordinates:
                coordinates[(x_coordiante, y_coordiante)] = 0

            coordinates[(x_coordiante, y_coordiante)] += total_steps

    return coordinates


def get_diration_and_length(step):
    """
    Extract the direction and length.
    """
    return step[0], int(step[1:])


def manhattan_distance(start, end):
    """
    Calculate the shorest distance between two set of coordiantes.
    Reference:
    https://dataaspirant.com/2015/04/11/five-most-popular-similarity-measures-implementation-in-python/
    """
    return sum(abs(a - b) for a, b in zip(start, end))


if __name__ == "__main__":
    main()
