package main

import (
	"fmt"
	"testing"
)

func TestIsWhitespace(t *testing.T) {

	var tests = []struct {
		input string
		want  bool
	}{
		{" ", true},
		{"abc", false},
		{"   abc", false},
		{"\n\t\v\f\r ", true},
		{"\n\t\v\f\r abc", false},
		{"abc \n\t\v\f\r 123", false},
		{"a b c 1 2 3", false},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s", test.input)

		t.Run(testName, func(t *testing.T) {
			got := IsWhitespace(test.input)

			if got != test.want {
				t.Errorf("Got %t, wanted %t", got, test.want)
			}
		})
	}
}

func TestParseFindAndReplaceSymbol(t *testing.T) {

	var tests = []struct {
		input string
		name  string
	}{
		{"<Patients>", "Patients"},
		{"<DateOfBirth transform=yearsElapsed", "DateOfBirth"},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s", test.input)

		t.Run(testName, func(t *testing.T) {
			name, _ := ParseFindAndReplaceSymbol(test.input)

			if name != test.name {
				t.Errorf("Got %s, wanted %s", name, test.name)
			}
		})
	}
}
