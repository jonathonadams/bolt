package router

import (
	"testing"
)

type NormalizeUrlAssertion struct {
	input    string
	expected string
}

func TestNormalizeUrl(t *testing.T) {

	tests := []NormalizeUrlAssertion{
		{input: "//", expected: ""},
		{input: "test", expected: "test"},
		{input: "test//", expected: "test"},
		{input: "test/test", expected: "test/test"},
		{input: "test//test", expected: "test/test"},
		{input: "/test//test", expected: "test/test"},
		{input: "/test/test", expected: "test/test"},
		{input: "/test/test/", expected: "test/test"},
		{input: "///test/test/", expected: "test/test"},
		{input: "///test///////test/////", expected: "test/test"},
	}

	for _, assert := range tests {
		actual := normalizeURL(assert.input)
		if actual != assert.expected {
			t.Errorf("Url was incorrect! Expected: %s, Actual: %s", assert.expected, actual)
		}
	}
}

func TestCreatePathSegments(t *testing.T) {
	var input string
	var actual []string

	input = ""
	actual = *createPathSegments(input)

	if len(actual) != 0 {
		t.Errorf("Path segments were not validated correctly!")
	}

	input = "/"
	actual = *createPathSegments(input)

	if len(actual) != 0 {
		t.Errorf("Path segments were not validated correctly!")
	}
}
