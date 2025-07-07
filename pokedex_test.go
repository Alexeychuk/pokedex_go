package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  bulbasaur  venuSAUR  charmander",
			expected: []string{"bulbasaur", "venusaur", "charmander"},
		},
		{
			input:    "blastoise",
			expected: []string{"blastoise"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		// Check the length of the actual slice against the expected slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test

		if len(actual) != len(c.expected) {
			t.Errorf("Lengths are not Equal")
		}

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
			if word != expectedWord {
				t.Errorf("Words are not Equal: %s - %s", word, expectedWord)
			}

		}
	}

}
