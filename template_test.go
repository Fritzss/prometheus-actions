package main

import "testing"

func TestTemplateFuncMap(t *testing.T) {
	data := map[string]interface{}{
		"arr": []string{
			"1",
			"2",
			"3",
		},
	}
	tests := map[string]string{
		`{{ "abcdef" | replace "abc" }}`:  "def",
		`{{ "" | default "abc" }}`:        "abc",
		`{{ default "abc" 0 }}`:           "0",
		`{{ length "abc" }}`:              "3",
		`{{ length .arr }}`:               "3",
		`{{ length 123 }}`:                "0",
		`{{ lower "ABC" }}`:               "abc",
		`{{ upper "abc" }}`:               "ABC",
		`{{ urlencode "?abcd=1&efg=2" }}`: "%3Fabcd%3D1%26efg%3D2",
		`{{ trim " abc " }}`:              "abc",
		`{{ yesno "yes" "no" true }}`:     "yes",
		`{{ yesno "yes" "no" false }}`:    "no",
	}
	for in, out := range tests {
		result, err := generateTemplate(in, "test", data)
		if err != nil {
			t.Error(err)
		}
		if result != out {
			t.Errorf("Must be %s, but got %s", out, result)
		}
	}
}
