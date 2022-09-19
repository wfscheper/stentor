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

// Package fragment parses fragment files.
package fragment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Fragment represents a single change or other news entry.
type Fragment struct {
	// Section is the short name of the section this fragment belongs to.
	Section string
	// Issue is the ID of any issues or pull requests to link to.
	Issue string
	// Text is the content of the change.
	Text string
}

// Deprecated: New returns a Fragment and the short name of the section it goes into.
func New(fn string) (Fragment, error) {
	f, err := Parse(fn)
	if err != nil {
		return Fragment{}, err
	}

	return *f, nil
}

// Parse parses the file fn into a Fragment structure.
//
// A fragment file follows the following naming convention:
// <issues>.<section>[.<summary>].(md|rst).
//
// The summary is optional and is ignored by Parse.
func Parse(fn string) (*Fragment, error) {
	parts := strings.Split(filepath.Base(fn), ".")
	var errMsg string
	switch {
	case len(parts) < 3:
		errMsg = "not enough parts"
	case parts[0] == "":
		errMsg = "empty issue"
	case parts[1] == "":
		errMsg = "empty section"
	}

	if errMsg != "" {
		return nil, fmt.Errorf("not a valid fragment file: %s", errMsg)
	}

	data, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	f := &Fragment{
		Issue:   parts[0],
		Section: parts[1],
		Text:    strings.TrimSpace(string(data)),
	}

	return f, nil
}
