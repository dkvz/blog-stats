package stats

import (
	"fmt"
	"regexp"
	"strings"
)

// Extract image legends to put them at the end (so they
// don't get matched by the following regexes):
var legendReg = regexp.MustCompile(`(?U)<[a-zA-Z-]+\s+class=.?image-legend.?>([^<]*)</[a-zA-Z-]+>`)

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
var tagsReg = regexp.MustCompile(`(?sU)<[^>]*>[^<]*</[^>]*>`)

// A bit redundant but the previous regex does not match
// single or self-closing tags
var simpleTagsReg = regexp.MustCompile(`<[^>]*>`)

var consSpaceReg = regexp.MustCompile(`\s{2,}`)

func WordCount(content *string) int {
	// Extract image legends to re-insert them later
	legendMatches := legendReg.FindAllStringSubmatch(*content, -1)
	legends := "\n"
	for _, m := range legendMatches {
		legends += fmt.Sprintf("%s\n", m[1])
	}

	res := *content + legends
	res = paReg.ReplaceAllString(res, "")

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
