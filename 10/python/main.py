#!/usr/bin/env python3

import math
import os
import sys
import time

SPACE = "◾️"
ROCKET = "🚀"
ASTROID = "🌘"
DESTROYED = "💥"


def main():
    """
    Main runs the astroid lookup program
    """
    astroid_map = []

    with open(sys.argv[1]) as map_file:
        for line in map_file:
            astroid_map.append([x for x in line.strip()])

    marked = mark_astroids(astroid_map)

    best = {"coordinates": (0, 0), "sees": 0}

    # Check the line from all marked coordiantes and the number of astroids
    # visible from that point.
    for coordinates in marked:
        all_coords = get_line_of_sight(astroid_map, coordinates)
        sees = len(all_coords)

        if sees > best["sees"] or best["sees"] == 0:
            best = {"coordinates": coordinates, "sees": sees}

    print(
        "most stars to be seen: {}, from {}".format(
            best["sees"], best["coordinates"]
        )
    )

    shots = 0
    animate = False
    two_hundred = None

    # Remove the best spot to ensure has_astroids report correct.
    best_x, best_y = best["coordinates"]
    astroid_map[best_x][best_y] = SPACE

    while has_astroids(astroid_map):
        all_coords = get_line_of_sight(astroid_map, best["coordinates"])

        sorted_coords = sorted(all_coords.items())
        sorted_coords.reverse()

        for _, data in sorted_coords:
            row, col = data["coordinates"]
            shots += 1

            if shots == 200:
                two_hundred = (col, row)

            astroid_map[row][col] = DESTROYED

            if animate:
                print_map(astroid_map, best["coordinates"])
                time.sleep(0.05)

    print("200th astroid shot down is at {}".format(two_hundred))


def has_astroids(astroid_map):
    """
    Check if there's any astroids left on the map.
    """
    for e in astroid_map:
        try:
            e.index(ASTROID)
            return True
        except ValueError:
            pass

    return False


def print_map(astroid_map, xy=None):
    """
    Print the map. Sorry, not for Windows (change to 'cls')
    """
    os.system("clear")

    if xy is not None:
        x, y = xy
        astroid_map[x][y] = ROCKET

    for x in astroid_map:
        print("  ".join(x))

    print()


def get_line_of_sight(astroid_map, coordinates):
    """
    Get the angle and ditance to all other astroids. If multiple astroids have
    the same angle, chose the closest one (the others are hidden behind it)
    """
    x1, y1 = coordinates
    all_coords = {}

    for x2, row in enumerate(astroid_map):
        for y2, col in enumerate(row):
            if (x1 == x2 and y1 == y2) or col != ASTROID:
                continue

            dx = x2 - x1
            dy = y2 - y1

            angle = math.atan2(dy, dx)
            degrees = math.degrees(angle)
            distance = abs(x1 - x2) + abs(y1 - y2)

            if angle in all_coords:
                if distance > all_coords[angle]["distance"]:
                    continue

            all_coords[angle] = {
                "distance": distance,
                "coordinates": (x2, y2),
                "degrees": degrees,
            }

    return all_coords


def mark_astroids(astroid_map):
    """
    Mark all coordiantes in the grid with an astroid (# sign)
    """
    astroids = []

    for row, _ in enumerate(astroid_map):
        for col, _ in enumerate(astroid_map[row]):
            if astroid_map[row][col] == "#":
                astroid_map[row][col] = ASTROID
                astroids.append((row, col))
            else:
                astroid_map[row][col] = SPACE

    return astroids


if __name__ == "__main__":
    main()
