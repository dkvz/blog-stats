package stats

import (
	"fmt"
	"regexp"
	"strings"
)

var paReg = regexp.MustCompile(`</?(p|h\d|a|i|b|small|strike|sub|sup|abbr|blockquote).*?>`)

// To use at the end to clean up the remaining tags
var tagsReg = regexp.MustCompile(`<.+?>.+?</.+?>`)

func WordCount(content string) int {
	// Remove all <p> and </p>
	// Remove all <a href="sdfksd" whatever> and </a>
	// Remove img and svg, some tags may not have closing element
	// Remove most tags and their inner content
	// We can keep blockquote contents
	res := strings.TrimSpace(paReg.ReplaceAllString(content, ""))
	res = tagsReg.ReplaceAllString(res, "")
	res = strings.ReplaceAll(res, "\n", " ")
	res = strings.ReplaceAll(res, "&nbsp;", " ")

	fmt.Println(res)

	// Remove consecutive spaces

	// Count spaces and add 1
	return strings.Count(res, " ") + 1
}
