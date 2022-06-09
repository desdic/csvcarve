package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/rs/zerolog/log"
)

var (
	filterCellRe = regexp.MustCompile(`\s*(?P<cellid>\d+)\s*(?P<operator>==|!=)\s*(?P<pattern>.*$)`)
	errNumRows   = errors.New("specified cell id is larger then number of rows")
)

type Cell struct {
	CellID      int
	Pattern     string
	ShouldMatch bool
	compiled    *regexp.Regexp
}

func buildCells(arr []string) ([]Cell, error) {
	cells := []Cell{}

	for _, v := range arr {
		match := filterCellRe.FindStringSubmatch(v)

		if len(match) != REGEXLEN {
			log.Error().Msgf("Last %s", v)

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

		cellid, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, fmt.Errorf("cell id must be a number: %w", err)
		}

		cells = append(cells, Cell{
			CellID:      cellid,
			Pattern:     userPattern,
			ShouldMatch: shouldmatch,
			compiled:    compiledPattern,
		})
	}

	return cells, nil
}

func matchCells(cells []Cell, row []string) (bool, error) {
	for _, cellpattern := range cells {
		if len(row) < cellpattern.CellID {
			return false, errNumRows
		}

		matched := cellpattern.compiled.Match([]byte(row[cellpattern.CellID]))

		if cellpattern.ShouldMatch {
			if !matched {
				return false, nil
			}
		} else {
			if matched {
				return false, nil
			}
		}
	}

	return true, nil
}
