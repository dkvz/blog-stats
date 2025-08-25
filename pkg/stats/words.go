package stats

import (
	"regexp"
	"strings"
)

var paReg = regexp.MustCompile(`</?p|a>`)

func WordCount(content string) int {
	// Remove all <p> and </p>
	// Remove all <a href="sdfksd" whatever> and </a>
	// Remove img and svg, some tags may not have closing element
	// Remove most tags and their inner content
	// We can keep blockquote contents
	res := paReg.ReplaceAllString(content, "")

	// Remove consecutive spaces

	// Count spaces and add 1
	return strings.Count(res, " ")
}
