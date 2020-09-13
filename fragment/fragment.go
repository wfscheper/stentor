package fragment

import (
	"bytes"
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
	// Text is the content of the change.
	Text string
	// Issue is the ID of any issues or pull requests to link to.
	Issue string
}

// New returns a Fragment and the short name of the section it goes into.
func New(fn string) (Fragment, string, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return Fragment{}, "", err
	}

	// strip whitespace
	data = bytes.TrimSpace(data)

	parts := strings.Split(filepath.Base(fn), ".")

	// error checking
	var errMsg string
	switch {
	case len(parts) > 4:
		errMsg = "too many parts"
	case len(parts) < 4:
		errMsg = "not enough parts"
	case parts[0] == "":
		errMsg = "empty section"
	case parts[2] == "":
		errMsg = "empty issue"
	}

	if errMsg != "" {
		return Fragment{}, "", fmt.Errorf("'%s' is not a valid fragment file: %s", filepath.Base(fn), errMsg)
	}

	return Fragment{string(data), parts[2]}, parts[0], nil
}
