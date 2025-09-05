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

Another possibility:
- Below 2000
- 2000 to 15000
- After 15000

-> Looks like I should replace 2000 with 1500.

## Early results

### No ignored articles
- Single factor: 0.120968
- Spread is 339.8901565484333
- Avg relative distance: 0.32710758070489915

- 2 factors
    - [0,2000] -> 0.105185
    - [2000,end] -> 0.123875
- Spread is 330.89133272733824
- Avg relative distance: 0.29922121826178866
    - Becomes 0.119952397029 with the first set of ignored articles

- 2 factors
    - [0,1500] -> 0.098807
    - [1500,end] -> 0.123875
- Spread is 330.82980966497405
- Avg relative distance: 0.2877032767174073
    - Becomes 0.1207980 with first set of ignored articles

- 3 factors
    - [0,2000] -> 0.105185
    - [2000,15000] -> 0.125479
    - [15000,end] -> 0.120256
- Spread is 346.81824677530494
- Avg relative distance: 0.29891266607055283
    - Becomes 0.11990945 with the first set of ignored articles

- Linear reg going through 0: 0.123463648676233*x 
- Spread is 330.58284107741406
- Avg relative distance: 0.3365037252608689
    - Becomes 0.1121213 with first set of ignored articles

- Linear reg: 0.123918557180491*x -26.1618180135265
- Spread is 329.8185921134582 but all the other metrics look better
- Avg relative distance: 0.1756664550494245
    - Becomes 0.1337925 with first set of ignored articles
- The extremes are way less extremes for full linear reg

### Ignored articles
Listing the outliers comma separated:
```
79,65,45,100,124,115,75,85
```
The last two are "large" articles.

We could just always remove:
```
79,65,45
```

#### New factors with ignored articles removed
Both plot and verify are done with the articles from the first list above removed.

- 1 factor: 0.124176
- Spread: 341.113727081086
- Avg relative distance: 0.1120250455537777

- 3 factors
    - [0,1500] -> 0.132460
    - [1500,15000] -> 0.124841
    - [15000,end] -> 0.120379
- Spread: 356.9821782
- Avg relative distance: 0.111

- Linear reg going through 0: 0.12351554373884
- Spread: 340.42297682384236
- Avg relative distance: 0.1121584

- Linear reg: 0.124012374330768*x -29.1540734021505
- Spread: 339.562641
- Avg relative distance: 0.13695656192296043

#### Without first major outlier

- Linear reg, Probably the best predictor
- y = 0.123919 x + -26.189935 (extremely close to previous values)
- Spread: 331.11352158905333
- Avg relative distance: 0.1558176278364528

- Triple factor
    - [0,2500] -> 0.112252
    - [2500,15000] -> 0.125722
    - [15000,end] -> 0.120256
- Spread: 348.218075
- Avg relative distance: 0.16170156106967315

