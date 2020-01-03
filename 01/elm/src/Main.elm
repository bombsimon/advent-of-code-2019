module Main exposing (..)

import Html exposing (Html)


solution1 : String -> Int
solution1 input =
    input
        |> parseInput
        |> List.foldl (\x acc -> acc + (x // 3) - 2) 0


solution2 : String -> Int
solution2 input =
    input
        |> parseInput
        |> List.map (\int -> int |> fuelRequired)
        |> List.foldl (+) 0



{--
Convert the input (string with newlines) to a list list of integers
--}


parseInput : String -> List Int
parseInput input =
    input
        |> String.split "\n"
        |> List.filter (\x -> not (x |> String.isEmpty))
        |> List.map
            (\int ->
                int
                    |> String.toInt
                    |> Maybe.withDefault 0
            )



{--
Process fuel recursive until it's below 0.
--}


fuelRequired : Int -> Int
fuelRequired int =
    let
        newFuel =
            int // 3 - 2
    in
    if newFuel < 0 then
        0

    else
        newFuel + (newFuel |> fuelRequired)


main : Html Never
main =
    Html.div []
        [ Html.p []
            [ Html.text "Part 1: "
            , Html.text (String.fromInt (solution1 puzzleInput))
            ]
        , Html.p []
            [ Html.text "Part 2: "
            , Html.text (String.fromInt (solution2 puzzleInput))
            ]
        ]


puzzleInput : String
puzzleInput =
    """
109364
144584
87498
130293
91960
117563
91730
138879
144269
89058
89982
115609
114728
85422
111803
148524
130035
107558
138936
95622
58042
50697
86848
123301
123631
143125
76434
78004
91115
89062
58465
141127
139993
80958
104184
145131
87438
74385
102113
97392
105986
58600
147156
54377
61409
73552
87138
63168
149602
111776
113191
80137
145985
145177
73192
141631
132979
52565
126574
92005
134655
115894
89175
127328
139873
50072
78814
134750
120848
132950
126523
58206
70885
85482
70889
100029
68447
95111
79896
138650
83079
83112
117762
57223
138122
145193
85251
103331
134501
77023
148189
141341
75994
67024
137767
86260
58705
58771
60684
79655
"""
