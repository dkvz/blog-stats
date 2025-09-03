# Research notes

## Stats thoughts
- I can probably derive the error range in minutes (or whatever time unit is being used) using the std dev when making the estimate

## Issue with big code snippets
Articles with big snippets or a lot of code have abnormally low word count to length ratio because we don't count code in the total.

The time estimate will thus end up higher than reality.

- Maybe we should count words in code?
- Course correct the length on the client by subtracting the length of code snippets
- Just live with it
- Extract code and count half of the words (rounded up) in the total -> Would lower the std dev but is this what I want?

**Remember some old articles don't use <pre><code> and just have <code>**.

I think I'll live with it, for article 153 the real value is 11.2 minutes and the estimate would be around 15. It's probably fine since people also sort of read the code. Maybe.

## Issue with lots of images
- Initially not counting the legends in the word count


-> Will add them in.
```
Current stats:
Avg: 0.116979   Min: 0.005051   Max: 0.155844
StdDev: 0.024556        Med: 0.119692

Adding the legends made the StdDev decrease by a minor amount:
Avg: 0.119005   Min: 0.005051   Max: 0.155844
StdDev: 0.023649        Med: 0.122434
```

## Outliers
Should probably remove them from the calculation with command line params or the .env.

I added the `-ignore-ids` param for that purpose.

The largest length outliers are:
- 107
- 141

ignoring them makes it easier to see the rest of the data.

## Computed length was wrong
I forgot that JavaScript uses UTF-16 whereas the Go program uses UTF-8 for everything.

The length for strings is thus inconsistent from one to the other.

I made it consistent by adding a conversion to a slice of 16 bit integers and counting the length of that (see `stats/words.go`).

Before the length computation change, stats were:
```
Word count stats:
Avg: 2134.770370        Min: 1.000000   Max: 20286.000000
StdDev: 3292.293461     Med: 1020.000000


Article length stats:
Avg: 17750.148148       Min: 189.000000 Max: 173626.000000
StdDev: 26939.841124    Med: 8413.000000


Ratio stats:
Avg: 0.118856   Min: 0.005051   Max: 0.153559
StdDev: 0.023456        Med: 0.122196
```

After the change:
```
Word count stats:
Avg: 2134.770370        Min: 1.000000   Max: 20286.000000
StdDev: 3292.293461     Med: 1020.000000


Article length stats:
Avg: 17438.325926       Min: 188.000000 Max: 170366.000000
StdDev: 26434.524024    Med: 8251.000000


Ratio stats:
Avg: 0.120968   Min: 0.005051   Max: 0.156182
StdDev: 0.024052        Med: 0.124370
```

The StdDev increased a bit and the length decreased globally.

## Factor ranges candidates
I'll start with just two:

- Below length 2000
- After that

