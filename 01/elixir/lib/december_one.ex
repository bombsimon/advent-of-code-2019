defmodule DecemberOne do
  @moduledoc """
  First Advent of Code task
  """

  @doc """
  Part one of the task.
  """
  def part_one(file) do
    {_, contents} = File.read(file)

    contents
    |> String.split("\n", trim: true)
    |> Enum.map(&String.to_integer/1)
    |> Enum.map(&floor_of_divided_by_three_minus_two/1)
    |> Enum.sum()
  end

  @doc """
  Part two of the task.
  """
  def part_two(file) do
    {_, contents} = File.read(file)

    contents
    |> String.split("\n", trim: true)
    |> Enum.map(&String.to_integer/1)
    |> Enum.map(&floor_of_divided_by_three_minus_two_until_zero/1)
    |> Enum.sum()
  end

  @doc """
  Divie by three and reduce two, round down and truncate to integer.
  it.

  ## Examples

      iex> DecemberOne.floor_of_divided_by_three_minus_two(1969)
      654
      iex> DecemberOne.floor_of_divided_by_three_minus_two(100756)
      33583
  """
  def floor_of_divided_by_three_minus_two(val) do
    (val / 3 - 2)
    |> Float.floor()
    |> Kernel.trunc()
  end

  @doc """
  Divie by three and reduce two until the result is 0 or less, add the sum of
  it.

  ## Examples

      iex> DecemberOne.floor_of_divided_by_three_minus_two_until_zero(100756)
      50346
  """
  def floor_of_divided_by_three_minus_two_until_zero(val) do
    case floor_of_divided_by_three_minus_two(val) do
      x when x > 0 ->
        x + floor_of_divided_by_three_minus_two_until_zero(x)

      _ ->
        0
    end
  end
end
