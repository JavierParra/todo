package apiServer

import (
	"strings"
	"regexp"
	"fmt"
	"strconv"
)

// Matcher compiles path strings into regular expressions to compare with url
// paths. It handles the following path strings:
// - /test/test     - Matches the path exactly
// - /test/:id      - Matches any path that starts with /test and is followed
// by anything not containing a slash (basically: [^\/]+). This will be
// captured and made available in the request.
// - /test/:id(\d+) - Matches any path that starts with /test and is followed
// by anything that matches the regular expression in the parens. This will be
// captured and made available in the request.
type Matcher struct {
	patterns []Pattern
	// An index of the string patterns with ids to handle duplicated patterns.
	ix map[string] string
}

type Pattern struct {
	reg    *regexp.Regexp
	groups []string
}

func (matcher *Matcher) Match(path string) (string, map[string] string) {
	for id, pattern := range matcher.patterns {
		re := pattern.reg

		match := re.FindStringSubmatch(path)

		if len(match) > 0 {
			fmt.Println(re, match, len(match))
			return strconv.FormatInt(int64(id), 16), fillMatches(&pattern, match)
		}
	}

	return "", make(map[string] string, 0)
}

func fillMatches(pattern *Pattern, match []string) map[string] string {
	matches := make(map[string] string, 0)

	for i, key := range pattern.groups {
		matches[key] = match[i+1]
	}

	return matches
}

// RegisterPattern registers a pattern with the matcher and returns an
// identifier or it.
func (matcher *Matcher) RegisterPattern (pattern string) string {
	if matcher.patterns == nil {
		matcher.patterns = make([]Pattern, 0)
		matcher.ix = make(map[string] string, 0)
	}

	if matcher.ix[pattern] != "" {
		return matcher.ix[pattern]
	}

	// This regexp finds all tokens. Either :key or :key(regex)
	re := regexp.MustCompile(`:([a-z]+)(\(.+?\))?`)
	matches := re.FindAllStringSubmatchIndex(pattern, -1)
	pat := compile(pattern, matches)
	id := strconv.FormatInt(int64(len(matcher.patterns)), 16)

	matcher.patterns = append(matcher.patterns, pat)
	matcher.ix[pattern] = id

	return id
}

// compile compiles a submatchIndex into a Pattern
func compile (pattern string, matches [][]int) Pattern {
	var match []int
	ix := 0
	parts := []string{"^"}
	keys := make([]string, 0)

	if len(matches) == 0 {
		parts = append(parts, regexp.QuoteMeta(pattern))
	}

	for len(matches) > 0 {
		match, matches = matches[0], matches[1:]
		var reg string
		key := string(pattern[match[2] : match[3]])

		if match[4] < 0 {
			// default regexp matches everything but a slash
			reg = "([^\\/]+)"
		} else {
			reg = string(pattern[match[4] : match[5]])
		}

		if ix < match[0] {
			parts = append(parts, regexp.QuoteMeta(string(pattern[ix : match[0]])))
		}

		parts = append(parts, reg)
		keys = append(keys, key)
		ix = match[1]
	}

	parts = append(parts, "$")

	fmt.Println("registering", strings.Join(parts, ""))
	return Pattern{regexp.MustCompile(strings.Join(parts, "")), keys}
}
