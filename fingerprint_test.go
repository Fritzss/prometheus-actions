package main

import "testing"

func TestBuildFingerprint(t *testing.T) {
	_, err := BuildFingerprint()
	if err != nil {
		t.Error(err)
	}
}
