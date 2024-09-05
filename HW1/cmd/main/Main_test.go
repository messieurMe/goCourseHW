package main

import (
	"testing"
)

func fooTest(t *testing.T) {
	a, b := 4, 10
	expected := 14

	actual := foo(a, b)

	if expected != actual {
		t.Error(
			"Expected $d, but got $",
			expected,
			actual,
		)
	}
}
