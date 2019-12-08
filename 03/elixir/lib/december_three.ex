defmodule DecemberThree do
  @moduledoc """
  Third Advent of Code task.
  """

  @doc """
  The solution for part one.
  """
  def part_one(file) do
    {_, contents} = File.read(file)

    [ first, second ] = contents
                       |> String.trim()
                       |> String.split("\n")
                       |> Enum.map(&draw_wire/1)

    first
    |> Map.delete("0,0")
    |> Enum.reduce(%{}, fn x, acc ->
      {k, v} = x

      case Map.has_key?(second, k) do
        true ->
          %{k => v}
          |> Map.merge(acc)
        false ->
          acc
      end
    end)
    |> Enum.map(fn {k, _v} -> manhattan_distance(k) end)
    |> Enum.min
  end

  @doc """
  The solution for part two.
  """
  def part_two(file) do
    {_, contents} = File.read(file)

    [ first, second ] = contents
                       |> String.trim()
                       |> String.split("\n")
                       |> Enum.map(&draw_wire/1)

    common = first
    |> Map.delete("0,0")
    |> Enum.reduce(%{}, fn x, acc ->
      {k, v} = x

      case Map.has_key?(second, k) do
        true ->
          %{k => v}
          |> Map.merge(acc)
        false ->
          acc
      end
    end)
    |> Enum.map(fn {k, _v} ->
      Map.get(first, k) + Map.get(second, k)
    end)
    |> Enum.min
  end

  @doc """
  Calculate manhattan distance, always from 0,0 because that's where I start.
  abs(x0-x1) + abs(y0-y1)
  """
  def manhattan_distance(xy) do
    [x, y] = xy
             |> String.split(",")
             |> Enum.map(fn x -> x |> String.to_integer end)

    Kernel.abs(0-x) + Kernel.abs(0-y)
  end

  @doc """
  Convert each step to a map of direction and length, then pass it to
  to_coordinates.
  """
  def draw_wire(steps) do
    steps
    |> String.split(",")
    |> Enum.map(fn a ->
      convert_steps(a)
    end)
    |> to_coordinates
  end

  @doc """
  Convert each step instruct to a map of direction and length.
  """
  def convert_steps(step) do
    direction =
      step
      |> String.at(0)

    length =
      step
      |> String.slice(1, String.length(step))
      |> String.to_integer

    %{"direction" => direction, "length" => length}
  end

  @doc """
  Fold each instruction and merge all coordinates and the distance to each
  coordinate in a map.

  ## Examples

      iex> [
      ...>   %{"direction" => "U", "length" => 3},
      ...>   %{"direction" => "R","length" => 5}
      ...> ] |> DecemberThree.to_coordinates
      %{
        "0,0" => 0,
        "0,1" => 1,
        "0,2" => 2,
        "0,3" => 3,
        "1,3" => 4,
        "2,3" => 5,
        "3,3" => 6,
        "4,3" => 7,
        "5,3" => 8
      }

  """
  def to_coordinates(instructions) do
    instructions
    |> Enum.reduce(%{}, fn instruction, acc ->
      # IO.puts "Iteration #{Map.get(acc, "steps")}: x,y: #{Map.get(acc, "x", 0)},#{Map.get(acc, "y", 0)}"

      update(
        %{}, instruction["direction"],
        Map.get(acc, "x", 0), Map.get(acc, "y", 0),
        Map.get(acc, "steps", 0),
        Map.get(acc, "steps", 0) + Map.get(instruction, "length")
      )
      |> Map.merge(acc, fn _k, v1, v2 ->
        # Always keep v1 values
        v1
      end)
    end)
    |> Map.delete("x")
    |> Map.delete("y")
    |> Map.delete("steps")
  end

  @doc """
  Caulculate the distance to each coordinate. Start by sending the current x and
  y coordinates combined with the current number of steps. The function will run
  until the number of steps reaches stop.

  ## Examples

      iex> DecemberThree.update(%{}, "D", 5, 8, 3, 6)
      %{
        "5,5" => 6,
        "5,6" => 5,
        "5,7" => 4,
        "5,8" => 3,
        "steps" => 6,
        "x" => 5,
        "y" => 5
      }
  """
  def update(co, direction, x, y, current, stop) do
    case current do
      x when x > stop ->
        co
      _ ->
        # Store the current position and step.
        co = co
        |> Map.put("x", x)
        |> Map.put("y", y)
        |> Map.put("steps", current)

        case direction do
          "U" ->
            co
            |> Map.put("#{x},#{y}", current)
            |> update(direction, x, y+1, current+1, stop)
          "D" ->
            co
            |> Map.put("#{x},#{y}", current)
            |> update(direction, x, y-1, current+1, stop)
          "L" ->
            co
            |> Map.put("#{x},#{y}", current)
            |> update(direction, x-1, y, current+1, stop)
          "R" ->
            co
            |> Map.put("#{x},#{y}", current)
            |> update(direction, x+1, y, current+1, stop)
        end
    end
  end
end
