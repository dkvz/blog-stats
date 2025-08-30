package stats

import (
	"regexp"
	"strings"
)

// Identify tags fr which we have to keep the inner content
// I certainly missed a whole bunch of these
var paReg = regexp.MustCompile(
	`</?(p|h\d|a|i|b|small|strike|sub|sup|abbr|span|blockquote|ul|ol|li|strong|em|del).*?>`,
)

// Looks like I need something extra for the comments
// I could probably combine it with the later one but
// it feels safer as its own thing
var comsReg = regexp.MustCompile(`(?sU)<!--.*-->`)

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
	res = comsReg.ReplaceAllString(res, "")
	// fmt.Printf("After pa and coms tags regexes:\n%v", res)
	res = tagsReg.ReplaceAllString(res, "")

	res = simpleTagsReg.ReplaceAllString(res, "")
	// fmt.Printf("After simple tags regex:\n%v", res)

	res = strings.ReplaceAll(res, "\n", " ")
	res = strings.ReplaceAll(res, "&nbsp;", " ")

	// Remove consecutive spaces
	res = strings.TrimSpace(consSpaceReg.ReplaceAllString(res, " "))

	// Count spaces
	return strings.Count(res, " ") + 1
}
