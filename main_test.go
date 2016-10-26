package main

import (
	"os"
	"testing"
)

type TestData struct {
	input      string
	envPattern string
	expected   string
}

func TestRenderingOutput(t *testing.T) {
	os.Setenv("VARIABLE", "value")

	cases := []TestData{
		TestData{"", "${ENV}", ""},
		TestData{"This is some text.", "${ENV}", "This is some text."},
		TestData{"This is some ${VARIABLE} in the text.", "${ENV}", "This is some value in the text."},
		TestData{"This is some ${VARIABLE} in the text.", "$<ENV>", "This is some ${VARIABLE} in the text."},
		TestData{"This is some $<VARIABLE> in the text.", "$<ENV>", "This is some value in the text."},
	}

	for _, c := range cases {
		output := render(c.input, c.envPattern)
		if output != c.expected {
			t.Errorf("Expected %s but got: %s", c.expected, output)
		}
	}
}
