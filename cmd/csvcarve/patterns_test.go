package main

import (
	"testing"
)

func TestMatchPatterns(t *testing.T) { //nolint:funlen
	t.Parallel()

	headers := map[string]int{
		"greeting":  0,
		"substance": 1,
		"personal":  2,
		"time":      3,
		"relation":  4,
	}

	tests := map[string]struct {
		input       []string
		row         []string
		headers     map[string]int
		expected    bool
		expectedErr bool
	}{
		"matchString": {
			input:       []string{"greeting==hello"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			headers:     headers,
			expected:    true,
			expectedErr: false,
		},
		"matchRegex": {
			input:       []string{"substance==^c"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			headers:     headers,
			expected:    true,
			expectedErr: false,
		},
		"noMatchString": {
			input:       []string{"substance!=coffee"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			headers:     headers,
			expected:    false,
			expectedErr: false,
		},
		"noMatchRegex": {
			input:       []string{"relation!=^f"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			headers:     headers,
			expected:    false,
			expectedErr: false,
		},
		"badRegex": {
			input:       []string{"relation!=BBB(((?!BBB).)*)EEE"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			headers:     headers,
			expected:    false,
			expectedErr: true,
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			patterns, err := buildPatterns(testCase.input)
			if testCase.expectedErr && err != nil {
				return
			}

			if !testCase.expectedErr && err != nil {
				t.Fatalf("did not expect the test to fail but it did: %v", err)
			}

			matched := matchPattern(patterns, testCase.headers, testCase.row)

			if matched != testCase.expected {
				t.Fatalf("expected %v got %v", testCase.expected, matched)
			}
		})
	}
}
