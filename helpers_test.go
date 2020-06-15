package main

import (
	"reflect"
	"testing"

	"github.com/prometheus/common/model"
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

func TestLabelSetEnviron(t *testing.T) {
	tests := []struct {
		in  model.LabelSet
		out []string
	}{
		{
			in: model.LabelSet{
				model.LabelName("__name__"): model.LabelValue("up"),
			},
			out: []string{
				"LABEL___NAME___0=up",
			},
		},
		{
			in: model.LabelSet{
				model.LabelName("__name__"): model.LabelValue("up"),
				model.LabelName("instance"): model.LabelValue("127.0.0.1:9090"),
			},
			out: []string{
				"LABEL___NAME___0=up",
				"LABEL_INSTANCE_0=127.0.0.1:9090",
			},
		},
	}
	for _, test := range tests {
		if !reflect.DeepEqual(test.out, LabelSetEnviron(0, test.in)) {
			assert.Equal(t, test.out, LabelSetEnviron(0, test.in))
		}
	}
}
