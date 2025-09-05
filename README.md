# Blog statistics and stuff

I actually don't do anything with my blog stats. At all. Nothing.

I started this project to try and find a relation between the word count of my articles and their length in character count as an overcomplicated math problem when I could just have the backend count words and cache that value in database.

As always, doing useless things to learn. Although, aren't all the things we do useless to the course of time and the universe?

## How I use this thing
TODO

Right now I manually save the factors or formulas to be used for the word count prediction. Would be nice the have some official "last state of things" with a special format that can be tested against the current DB for the quality of predictions.

## References
- [SQlite lib docs](https://practicalgobook.net/posts/go-sqlite-no-cgo/)
- [Go-echarts online doc and examples](https://go-echarts.github.io/)

## Style issues
My project structure is a bit out there. It would be better to provide structs with methods in packages but I often export one or two functions instead.

Might refactor that someday if I really got nothing else to do.

## TODO
- [ ] Plot mode should probably be called "analyse" at this point
- [ ] I could add some plot to verify mode as well
- [ ] Verify mode should also compute the std dev and how much it is compared to the actual value (some percentage of "error")
- [x] Use readonly or writeonly channels when possible - I always forget about that
- [x] I also need the article ID for word count stats, helps debugging strange results
- [x] I need to test word count with HTML comments at the start
