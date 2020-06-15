package main

import (
	"fmt"
	"strings"

	"github.com/prometheus/common/model"
)

const (
	labelSetFormat = "LABEL_%s_%d=%s"
)

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func LabelSetEnviron(id int, labelSet model.LabelSet) []string {
	var env []string
	for labelKey, labelValue := range labelSet {
		key := strings.ToUpper(string(labelKey))
		kv := fmt.Sprintf(labelSetFormat, key, id, labelValue)
		env = append(env, kv)
	}
	return env
}

func LabelSetSliceEnviron(labelSetSlice []model.LabelSet) []string {
	var env []string
	for id, labelSet := range labelSetSlice {
		labelSetEnv := LabelSetEnviron(id, labelSet)
		env = append(env, labelSetEnv...)
	}
	return env
}
