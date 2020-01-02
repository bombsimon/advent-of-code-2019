#!/usr/bin/env python3

import sys
import math


def main():
    """
    Main runs the astroid lookup program
    """

    data = []

    # Parse each row into left- and right hand side dictionaries holding
    # materials and amount.
    with open(sys.argv[1]) as f:
        for line in f:
            lhs, rhs = line.strip().split(" => ")

            data.append([parse_material(lhs), parse_material(rhs)])

    # Build a map telling what resources n element requires.
    require_map = {}

    for x in data:
        lhs, rhs = x
        rkey = next(iter(rhs))

        require_map[rkey] = {"amount": rhs[rkey], "requires": {}}

        for k, v in lhs.items():
            require_map[rkey]["requires"][k] = v

    # Part one, look for number of ores for 1 FUEL.
    total_ore = look_for("FUEL", 1, require_map)
    print()
    print("part 1: total ore: {:,}".format(total_ore))

    # Part two, do a binary search to try to find the highest number within the
    # cargo boundaries of FUEL production.
    cargo_ore = 10 ** 12
    max_fuel = 0

    # Assume 1 ORE never produces more than 1 FUEL
    fuel_low, fuel_high = (1, 10 ** 12)

    # Flag if the last result before finding middle was too high.
    last_high = False

    while fuel_low != fuel_high:
        max_fuel = (fuel_low + fuel_high) // 2
        total_ore = look_for("FUEL", max_fuel, require_map)

        if total_ore < cargo_ore:
            fuel_low, fuel_high = max_fuel, fuel_high
            last_high = False
        else:
            fuel_low, fuel_high = fuel_low, max_fuel
            last_high = True

    if last_high:
        max_fuel -= 1

    print()
    print("part 2: max fuel from {:,} ore: {:,}".format(cargo_ore, max_fuel))


def look_for(material, count, require_map, spare={}):
    """
    Look for a given count of a given material.
    """
    if material == "ORE":
        return count

    if material in spare:
        spares = spare[material]

        if spare[material] <= count:
            count -= spares
            spare[material] = 0  # Reset the spares for this material, all used
        else:
            spare[material] -= count

            return 0  # We had spares for all parts.

    material_info = require_map[material]
    required_reactions = math.ceil(count / material_info["amount"])
    total_production = required_reactions * material_info["amount"]
    production_spare = total_production - count

    # Store the spares after producing given material.
    if material not in spare:
        spare[material] = 0

    spare[material] += production_spare

    ore = 0
    for m, c in material_info["requires"].items():
        # No traversal if we don't need anything
        if required_reactions == 0:
            continue

        ore += look_for(m, required_reactions * c, require_map, spare)

    return ore


def parse_material(inp):
    """
    Parse each material and create a map with it's name and the amount.
    """
    material_map = {}

    for x in inp.split(", "):
        amount, name = x.split(" ")
        material_map[name] = int(amount)

    return material_map


if __name__ == "__main__":
    main()
