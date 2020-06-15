package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStandardizeSpaces(t *testing.T) {
	input := `line 1
	line 2
	line 3`
	out := StandardizeSpaces(input)
	assert.Equal(t, "line 1 line 2 line 3", out)

	input = `
	(
        node_filesystem_free{instance="localhost", mountpoint="/var/lib/docker"} /
        node_filesystem_size{instance="localhost", mountpoint="/var/lib/docker"}
	) * 100 < 100
	`
	out = StandardizeSpaces(input)
	assert.Equal(t, `( node_filesystem_free{instance="localhost", mountpoint="/var/lib/docker"} / node_filesystem_size{instance="localhost", mountpoint="/var/lib/docker"} ) * 100 < 100`, out)
}
