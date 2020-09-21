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

package fragment

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Fragment represents a single change or other news entry.
//
// A fragment file follows the following naming convention:
// <section>.<summary>.<issue>.(md|rst). Note that the summary is ignored by
// stentor, and is just present to describe the fragment for people browsing the
// fragments.
type Fragment struct {
	// Section is the short name of the section this fragment belongs to.
	Section string
	// Issue is the ID of any issues or pull requests to link to.
	Issue string
	// Text is the content of the change.
	Text string
}

// New returns a Fragment and the short name of the section it goes into.
func New(fn string) (Fragment, error) {
	var f Fragment

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
		return f, fmt.Errorf("not a valid fragment file: %s", errMsg)
	}

	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return f, err
	}

	f = Fragment{
		Issue:   parts[0],
		Section: parts[1],
		Text:    strings.TrimSpace(string(data)),
	}

	return f, err
}
