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
	"embed"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates
var fs embed.FS

var (
	funcMap = template.FuncMap{
		"indent": indent,
		"repeat": repeat,
		"sum":    sum,
	}
)

// New returns the named template
func New(name string) (*template.Template, error) {
	data, err := fs.ReadFile("templates/" + name)
	if err != nil {
		return nil, err
	}

	return template.New(name).Funcs(funcMap).Parse(string(data))
}

// Parse returns the template parsed from file fn
func Parse(fn string) (*template.Template, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	return template.New(filepath.Base(fn)).Funcs(funcMap).Parse(string(data))
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
