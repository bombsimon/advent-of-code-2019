defmodule DecemberFour do
  @moduledoc """
  Fourth Advent of Code task.
  """

  @start 246_515
  @stop 739_105

  @doc """
  Part one of day 4.
  """
  def part_one do
    list_to_duples_and_triples()
    |> Enum.filter(fn a -> Enum.at(a[:result], 0) end)
    |> Kernel.length()
  end

  @doc """
  Part two of day 4.
  """
  def part_two do
    list_to_duples_and_triples()
    |> Enum.filter(fn a -> Enum.at(a[:result], 0) && !Enum.at(a[:result], 1) end)
    |> Kernel.length()
  end

  @doc """
  Convert each list of digits to a list of two booleans.
    0: Has duplicates
    1: Has triplets overlaping all duplicates
  """
  def list_to_duples_and_triples do
    @start..@stop
    |> Enum.map(&Integer.digits/1)
    |> Enum.map(fn a ->
      %{
        :list => a,
        :result => check_progression(a)
      }
    end)
  end

  @doc """
  Traverse each digit and store the position from the tail for each place a
  duplicate and triplet is found. If a triplet is found, store where both
  preceding digits were found. End by checking if there were any dupes and if
  there were any dupes not overlaped by a triplet.

  ## Examples

  iex> DecemberFour.check_progression([ 1, 2, 1 ])
  [false, false]

  iex> DecemberFour.check_progression([ 1, 1, 2 ])
  [true, false]

  iex> DecemberFour.check_progression([ 1, 1, 1 ])
  [true, true]

  iex> DecemberFour.check_progression([ 1, 1, 2, 2, 2 ])
  [true, false]

  iex> DecemberFour.check_progression([ 1, 1, 1, 2, 2 ])
  [true, false]

  iex> DecemberFour.check_progression([ 1, 1, 1, 2, 2, 2])
  [true, true]
  """
  def check_progression(
        [digit | tail],
        previous \\ 0,
        secondPrevious \\ 0,
        dupes \\ [],
        trips \\ []
      ) do
    dupes =
      if digit == previous do
        [Kernel.length(tail) | dupes]
      else
        dupes
      end

    trips =
      if digit == previous && digit == secondPrevious do
        [Kernel.length(tail) | [Kernel.length(tail) + 1 | trips]]
      else
        trips
      end

    case digit do
      x when x < previous ->
        [false, false]

      _ ->
        if Kernel.length(tail) == 0 do
          hasDupes =
            dupes
            |> Kernel.length() > 0

          hasTrips =
            dupes
            |> Enum.filter(fn el -> !Enum.member?(trips, el) end)
            |> Kernel.length() == 0

          [hasDupes, hasTrips]
        else
          check_progression(tail, digit, previous, dupes, trips)
        end
    end
  end
end
