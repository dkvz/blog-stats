# Research notes

## Issue with big code snippets
Articles with big snippets or a lot of code have abnormally low word count to length ratio because we don't count code in the total.

The time estimate will thus end up higher than reality.

- Maybe we should count words in code?
- Course correct the length on the client by subtracting the length of code snippets
- Just live with it
- Extract code and count half of the words (rounded up) in the total -> Would lower the std dev but is this what I want?

I think I'll live with it, for article 153 the real value is 11.2 minutes and the estimate would be around 15. It's probably fine since people also sort of read the code. Maybe.

## Issue with lots of images
- Initially not counting the legends in the word count

-> Will add them in.

Current stats:
Avg: 0.116979   Min: 0.005051   Max: 0.155844
StdDev: 0.024556        Med: 0.119692
