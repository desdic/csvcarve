package main

import (
	"fmt"
	"regexp"
)

var filterPatternRe = regexp.MustCompile(`\s*(?P<field>.*)\s*(?P<operator>==|!=)\s*(?P<pattern>.*$)`)

type Pattern struct {
	CellName    string
	Pattern     string
	ShouldMatch bool
	compiled    *regexp.Regexp
}

func buildPatterns(arr []string) ([]Pattern, error) {
	patterns := []Pattern{}

	for _, v := range arr {
		match := filterPatternRe.FindStringSubmatch(v)

		if len(match) != REGEXLEN {
			return nil, errParse
		}

		shouldmatch := true
		if match[2] == "!=" {
			shouldmatch = false
		}

		userPattern := match[3]

		compiledPattern, err := regexp.Compile(userPattern)
		if err != nil {
			return nil, fmt.Errorf("failed to compile %s, %w", userPattern, err)
		}

		patterns = append(patterns, Pattern{
			CellName:    match[1],
			Pattern:     userPattern,
			ShouldMatch: shouldmatch,
			compiled:    compiledPattern,
		})
	}

	return patterns, nil
}

func matchPattern(patterns []Pattern, headers map[string]int, row []string) bool {
	for _, pattern := range patterns {
		index := headers[pattern.CellName]

		matched := pattern.compiled.Match([]byte(row[index]))

		if pattern.ShouldMatch {
			if !matched {
				return false
			}
		} else {
			if matched {
				return false
			}
		}
	}

	return true
}
