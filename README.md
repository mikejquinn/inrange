This is a very small cmd line utility for filtering a stream of numbers by a given numerical range.

Given the file `numbers.txt`:

    -10
    -12
    3
    6.9
    0
    -1
    4
    -8

Then:

    #> cat numbers.txt | inrange '[-8,4)'
    3
    0
    -1
