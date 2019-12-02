defmodule DecemberTwo do
  @moduledoc """
  Second Advent of Code task.
  """

  @noun_position 1
  @verb_position 2
  @part_one_noun 12
  @part_one_verb 2
  @part_two_target 19_690_720
  @part_two_max_val 99

  @doc """
  Read file and convert the comma separated list to a map with the index as key.
  """
  def init_map(file) do
    {_, contents} = File.read(file)

    contents
    |> String.trim()
    |> String.split(",", trim: true)
    |> Enum.map(&String.to_integer/1)
    |> Enum.with_index()
    |> Enum.map(fn {val, idx} -> {idx, val} end)
    |> Map.new()
  end

  @doc """
  Part one of the task.
  """
  def part_one(file) do
    init_map(file)
    |> update_noun_and_verb(@part_one_noun, @part_one_verb)
    |> check_code(0)
    |> Map.get(0)
  end

  @doc """
  Part two of the task.
  """
  def part_two(file) do
    to_check = for noun <- 0..@part_two_max_val, verb <- 0..@part_two_max_val, do: {noun, verb}

    init_map(file)
    |> check_noun_and_verb_to(@part_two_target, to_check)
  end

  @doc """
  """
  def check_noun_and_verb_to(sequence, look_for, [{noun, verb} | rest]) do
    result =
      sequence
      |> update_noun_and_verb(noun, verb)
      |> check_code(0)
      |> Map.get(0)

    case result do
      ^look_for ->
        100 * noun + verb

      _ ->
        check_noun_and_verb_to(sequence, look_for, rest)
    end
  end

  @doc """
  Update noun and verb by setting position 1 and 2 of the map.
  """
  def update_noun_and_verb(sequence, noun, verb) do
    sequence
    |> Map.put(@noun_position, noun)
    |> Map.put(@verb_position, verb)
  end

  @doc """
  Check the operation code and perform appropreate action.
  """
  def check_code(sequence, current_position) do
    op_cde = Map.get(sequence, current_position)
    n1_pos = Map.get(sequence, current_position + 1)
    n2_pos = Map.get(sequence, current_position + 2)
    sr_pos = Map.get(sequence, current_position + 3)
    n1_val = Map.get(sequence, n1_pos)
    n2_val = Map.get(sequence, n2_pos)

    case op_cde do
      1 ->
        Map.put(sequence, sr_pos, n1_val + n2_val)
        |> check_code(current_position + 4)

      2 ->
        Map.put(sequence, sr_pos, n1_val * n2_val)
        |> check_code(current_position + 4)

      99 ->
        sequence
    end
  end
end
