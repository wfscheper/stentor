// Copyright © 2020 The Stentor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package templates

import (
	"strings"
	"text/template"

	rice "github.com/GeertJohan/go.rice"
)

//go:generate ../../tools/bin/rice embed-go

var (
	funcMap = template.FuncMap{
		"indent": indent,
		"repeat": repeat,
		"sum":    sum,
	}
)

// New returns the named template
func New(name string) (*template.Template, error) {
	box, err := rice.FindBox("templates")
	if err != nil {
		return nil, err
	}

	templateStr, err := box.String(name)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(funcMap).Parse(templateStr)
}

// Parse returns the template parsed from file fn
func Parse(fn string) (*template.Template, error) {
	return template.ParseFiles(fn)
}

// template functions

// indent pads every line in s after the first with n spaces.
//
// This transforms:
// "Line1\nLine2" into "Line1\n  Line2".
func indent(n int, s string) string {
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if i != 0 && line != "" {
			lines[i] = strings.Repeat(" ", n) + line
		}
	}
	return strings.Join(lines, "\n")
}

// repeat returns the string s repeated n times.
func repeat(n int, s string) string {
	return strings.Repeat(s, n)
}

// sum returns the sum of its arguments
func sum(ns ...int) (i int) {
	for _, n := range ns {
		i += n
	}
	return i
}
