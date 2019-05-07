package main

import (
	"testing"
)

func TestStandardizeSpaces(t *testing.T) {
	input := `line 1
	line 2
	line 3`
	out := StandardizeSpaces(input)
	if out != "line 1 line 2 line 3" {
		t.Errorf("Failed match string: %s", out)
	}
	input = `
	(
        node_filesystem_free{instance="localhost", mountpoint="/var/lib/docker"} /
        node_filesystem_size{instance="localhost", mountpoint="/var/lib/docker"}
	) * 100 < 100
	`
	out = StandardizeSpaces(input)
	if out != `( node_filesystem_free{instance="localhost", mountpoint="/var/lib/docker"} / node_filesystem_size{instance="localhost", mountpoint="/var/lib/docker"} ) * 100 < 100` {
		t.Errorf("Failed match string: %s", out)
	}
}
