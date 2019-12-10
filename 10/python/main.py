#!/usr/bin/env python3

import sys


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
        x1, y1 = coordinates
        all_coords = {}

        for i, _ in enumerate(astroid_map):
            for j, _ in enumerate(astroid_map[i]):
                seen_slopes = {}
                coordinates_on_line = bresenhams_line_algorithm(
                    (x1, y1), (i, j)
                )

                x0, y0 = coordinates_on_line[0]
                other = coordinates_on_line[1:]

                # All slopes
                slopes = [
                    (y - y0) / (x - x0) if x != x0 else x for x, y in other
                ]

                for idx, slope in enumerate(slopes):
                    # If we've seen an astroid with this slope for these
                    # coordiantes it's blocking the rest.
                    if slope in seen_slopes:
                        continue

                    fx, fy = other[idx]
                    if astroid_map[fx][fy] == "#":
                        all_coords[other[idx]] = True
                        seen_slopes[slope] = other[idx]

        sees = len(all_coords)

        if sees > best["sees"] or best["sees"] == 0:
            best = {"coordinates": coordinates, "sees": sees}

    print(
        "most stars to be seen: {}, from {}".format(
            best["sees"], best["coordinates"]
        )
    )


def mark_astroids(astroid_map):
    """
    Mark all coordiantes in the grid with an astroid (# sign)
    """
    astroids = []

    for row, _ in enumerate(astroid_map):
        for col, _ in enumerate(astroid_map[row]):
            if astroid_map[row][col] == "#":
                astroids.append((row, col))

    return astroids


def bresenhams_line_algorithm(start, end):
    """Bresenham's Line Algorithm
    Get all points between start (x0,y0) and end (x1,y1)
    https://en.wikipedia.org/wiki/Bresenham%27s_line_algorithm
    http://www.roguebasin.com/index.php?title=Bresenham%27s_Line_Algorithm#Python
    """
    # Setup initial conditions
    start_x, start_y = start
    end_x, end_y = end
    delta_x = end_x - start_x
    delta_y = end_y - start_y

    # Determine how steep the line is
    is_steep = abs(delta_y) > abs(delta_x)

    # Rotate line
    if is_steep:
        start_x, start_y = start_y, start_x
        end_x, end_y = end_y, end_x

    # Swap start and end points if necessary and store swap state
    swapped = False
    if start_x > end_x:
        start_x, end_x = end_x, start_x
        start_y, end_y = end_y, start_y
        swapped = True

    # Recalculate differentials
    delta_x = end_x - start_x
    delta_y = end_y - start_y

    # Calculate error
    error = int(delta_x / 2.0)
    ystep = 1 if start_y < end_y else -1

    # Iterate over bounding box generating points between start and end
    y = start_y
    points = []

    for x in range(start_x, end_x + 1):
        coord = (y, x) if is_steep else (x, y)
        points.append(coord)
        error -= abs(delta_y)
        if error < 0:
            y += ystep
            error += delta_x

    # Reverse the list if the coordinates were swapped
    if swapped:
        points.reverse()

    return points


if __name__ == "__main__":
    main()
