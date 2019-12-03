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

    print(min(all_distances))


def mark(steps):
    """
    Mark all the coordinates extracted byt traversing the steps in each
    direction.
    """
    x_coordiante, y_coordiante = 0, 0
    coordinates = {}

    for step in steps:
        direction, length = get_diration_and_length(step)

        if direction == "U":
            for _ in range(length):
                x_coordiante -= 1
                coordinates[(x_coordiante, y_coordiante)] = True

        if direction == "D":
            for _ in range(length):
                x_coordiante += 1
                coordinates[(x_coordiante, y_coordiante)] = True

        if direction == "L":
            for _ in range(length):
                y_coordiante -= 1
                coordinates[(x_coordiante, y_coordiante)] = True

        if direction == "R":
            for _ in range(length):
                y_coordiante += 1
                coordinates[(x_coordiante, y_coordiante)] = True

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
