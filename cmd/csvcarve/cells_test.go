package main

import (
	"testing"
)

func TestMatchCells(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := map[string]struct {
		input       []string
		row         []string
		expected    bool
		expectedErr bool
	}{
		"matchString": {
			input:       []string{"0==hello"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			expected:    true,
			expectedErr: false,
		},
		"matchRegex": {
			input:       []string{"1==^c"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			expected:    true,
			expectedErr: false,
		},
		"noMatchString": {
			input:       []string{"1==dog"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			expected:    false,
			expectedErr: false,
		},
		"noMatchRegex": {
			input:       []string{"0==^h", "1!=^c"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			expected:    false,
			expectedErr: false,
		},
		"badRegex": {
			input:       []string{"0==BBB(((?!BBB).)*)EEE"},
			row:         []string{"hello", "coffee", "my", "old", "friend"},
			expected:    false,
			expectedErr: true,
		},
	}

	for name, testCase := range tests {
		testCase := testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()
			cells, err := buildCells(testCase.input)
			if testCase.expectedErr && err != nil {
				return
			}

			if !testCase.expectedErr && err != nil {
				t.Fatalf("did not expect the test to fail but it did: %v", err)
			}

			matched, err := matchCells(cells, testCase.row)

			if matched != testCase.expected {
				t.Fatalf("expected %v got %v", testCase.expected, matched)
			}

			if testCase.expectedErr && err != nil {
				t.Fatal("expected to fail but didn't")
			}
		})
	}
}
