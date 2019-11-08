package main

import (
	"bytes"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"text/template"
)

var TemplateFuncMap = template.FuncMap{
	"replace": func(s1 string, s2 string) string {
		defer recovery()

		return strings.Replace(s2, s1, "", -1)
	},
	"default": func(arg interface{}, value interface{}) interface{} {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Slice, reflect.Array, reflect.Map:
			if v.Len() == 0 {
				return arg
			}
		case reflect.Bool:
			if !v.Bool() {
				return arg
			}
		}

		return value
	},
	"length": func(value interface{}) int {
		defer recovery()

		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.Slice, reflect.Array, reflect.Map:
			return v.Len()
		case reflect.String:
			return len([]rune(v.String()))
		}

		return 0
	},
	"lower": func(s string) string {
		defer recovery()

		return strings.ToLower(s)
	},
	"upper": func(s string) string {
		defer recovery()

		return strings.ToUpper(s)
	},
	"urlencode": func(s string) string {
		defer recovery()

		return url.QueryEscape(s)
	},
	"trim": func(s string) string {
		defer recovery()

		return strings.TrimSpace(s)
	},
	"yesno": func(yes string, no string, value bool) string {
		defer recovery()

		if value {
			return yes
		}

		return no
	},
}

func recovery() {
	recover()
}

func GenerateTemplate(templ, name string, data interface{}) (string, error) {
	var templateEng *template.Template
	buf := bytes.NewBufferString("")
	templateEng = template.New(name)
	if messageTempl, err := templateEng.Funcs(TemplateFuncMap).Parse(templ); err != nil {
		return "", fmt.Errorf("failed to parse template for %s: %v", name, err)
	} else if err := messageTempl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %v", name, err)
	}
	return buf.String(), nil
}
