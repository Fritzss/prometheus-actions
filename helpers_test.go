package main

import "testing"

func TestStandardizeSpaces(t *testing.T) {
	input := `line 1
	line 2
	line 3`
	out := StandardizeSpaces(input)
	if out != "line 1 line 2 line 3" {
		t.Errorf("Failed match string: %s", out)
	}
}
