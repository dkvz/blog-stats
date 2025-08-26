package stats

import (
	"fmt"
	"regexp"
	"strings"
)

var paReg = regexp.MustCompile(
	`</?(p|h\d|a|i|b|small|strike|sub|sup|abbr|blockquote|ul|ol|li|strong|em|del).*?>`,
)

// To use at the end to clean up the remaining tags
// var tagsReg = regexp.MustCompile(`<.+?>.+?</.+?>`)
// Will also destroy image legends but I can live with that
// var tagsReg = regexp.MustCompile(`(?s)<[^>]*>[^>]*<[^>]*>`)
// s flag was needed for "." to also match line feeds
var tagsReg = regexp.MustCompile(`(?s)<[^>]*>.*<[^>]*>`)

// A bit redundant but the previous regex does not match
// single or self-closing tags
var simpleTagsReg = regexp.MustCompile(`<[^>]*>`)

var consSpaceReg = regexp.MustCompile(`\s{2,}`)

func WordCount(content *string) int {
	// Remove all <p> and </p>
	// Remove all <a href="sdfksd" whatever> and </a>
	// Remove img and svg, some tags may not have closing element
	// Remove most tags and their inner content
	// We can keep blockquote contents
	res := paReg.ReplaceAllString(*content, "")
	res = tagsReg.ReplaceAllString(res, "")

	res = simpleTagsReg.ReplaceAllString(res, "")

	res = strings.ReplaceAll(res, "\n", " ")
	res = strings.ReplaceAll(res, "&nbsp;", " ")

	// Remove consecutive spaces
	res = strings.TrimSpace(consSpaceReg.ReplaceAllString(res, " "))

	fmt.Printf("Cleaned up before word count:\n\n%s\n", res)

	// Count spaces
	return strings.Count(res, " ") + 1
}
