package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildFingerprint(t *testing.T) {
	_, err := BuildFingerprint()
	assert.NoError(t, err)
}
