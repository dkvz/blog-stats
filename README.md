# Blog statistics and stuff

I actually don't do anything with my blog stats. At all. Nothing.

I started this project to try and find a relation between the word count of my articles and their length in character count as an overcomplicated math problem when I could just have the backend count words and cache that value in database.

As always, doing useless things to learn. Although, aren't all the things we do useless to the course of time and the universe?

## References
- [SQlite lib docs](https://practicalgobook.net/posts/go-sqlite-no-cgo/)

## Problems
- The ratio is too outlandish when there are too many images in the content, because we use the length with everything in it and the word count with all of the HTML cleaned up.
    - Consider removing image divs from the length both in here and on the JS client
    - Include the legends in the word count

## TODO
- [ ] Don't forget to use the special length function from words.go for computing UTF-16 equivalent string length that JS uses
- [x] Use readonly or writeonly channels when possible - I always forget about that
- [x] I also need the article ID for word count stats, helps debugging strange results
- [x] I need to test word count with HTML comments at the start
