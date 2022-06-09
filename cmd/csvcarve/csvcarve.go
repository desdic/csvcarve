package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const REGEXLEN = 4

var errParse = errors.New("filter must be key != or == pattern")

func buildHeaders(row []string, patterns []Pattern) map[string]int {
	var err error

	headers := make(map[string]int)

	for i := range row {
		headers[row[i]] = i
	}

	for index, pattern := range patterns {
		if _, ok := headers[pattern.CellName]; !ok {
			log.Error().Msgf("Error: %s does not exist in the header", pattern.CellName)

			os.Exit(1)
		}

		patterns[index].compiled, err = regexp.Compile(pattern.Pattern)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to compile %s", pattern.Pattern)

			os.Exit(1)
		}
	}

	return headers
}

func matchCSV(csvfile io.Reader, noHeader bool, patterns []Pattern, cells []Cell) {
	csvReader := csv.NewReader(csvfile)
	csvReader.LazyQuotes = true
	csvReader.TrimLeadingSpace = true

	gotHeader := false

	var headers map[string]int

	for {
		row, err := csvReader.Read()

		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Error().Err(err).Msgf("Failed to read line")

			return
		}

		if !gotHeader && !noHeader {
			gotHeader = true

			headers = buildHeaders(row, patterns)
			fmt.Printf("%s\n", strings.Join(row, ",")) //nolint:forbidigo

			continue
		}

		if !matchPattern(patterns, headers, row) {
			continue
		}

		cellMatched, err := matchCells(cells, row)
		if err != nil {
			log.Error().Err(err).Msg("matchCells failed")

			return
		}

		if !cellMatched {
			continue
		}

		fmt.Printf("%s\n", strings.Join(row, ",")) //nolint:forbidigo
	}
}

func NewCommand() *cobra.Command {
	var (
		filters     []string
		cellfilters []string
		noHeader    bool
		filename    string
	)

	c := &cobra.Command{
		Use: "cmd",
		Run: func(_ *cobra.Command, _ []string) {
			patterns, err := buildPatterns(filters)
			if err != nil {
				log.Error().Err(err).Msg("Failed to build patterns")
				os.Exit(1)
			}

			cells, err := buildCells(cellfilters)
			if err != nil {
				log.Error().Err(err).Msg("Failed to build patterns")
				os.Exit(1)
			}

			if len(patterns) < 1 && len(cells) < 1 {
				log.Error().Err(err).Msg("No patterns defined")
				os.Exit(1)
			}

			if filename != "" {
				csvinput, err := os.Open(filename)
				if err != nil {
					log.Error().Err(err).Msgf("Unable to open %s", filename)
					os.Exit(1)
				}

				defer csvinput.Close()

				matchCSV(csvinput, noHeader, patterns, cells)

				return
			}

			csvinput := bufio.NewReader(os.Stdin)
			matchCSV(csvinput, noHeader, patterns, cells)
		},
	}

	c.Flags().StringVarP(&filename, "filename", "f", "", "")
	c.Flags().StringArrayVarP(&filters, "filter", "p", []string{}, "")
	c.Flags().StringArrayVarP(&cellfilters, "cell", "a", []string{}, "")
	c.Flags().BoolVarP(&noHeader, "noheader", "n", false, "CSV has no header")

	return c
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05",
		},
	)
	if err := NewCommand().Execute(); err != nil {
		log.Error().Err(err).Msg("Command failed")
	}
}
