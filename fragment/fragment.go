package fragment

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Fragment represents a single change or other news entry.
//
// A fragment file follows the following naming convention: <issue #>.<section>.(md|rst)
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

	parts := strings.Split(filepath.Base(fn), ".")

	// error checking
	var errMsg string
	switch {
	case len(parts) > 3:
		errMsg = "too many parts"
	case len(parts) < 3:
		errMsg = "not enough parts"
	case parts[0] == "":
		errMsg = "empty issue"
	case parts[1] == "":
		errMsg = "empty section"
	}

	if errMsg != "" {
		return Fragment{}, "", fmt.Errorf("'%s' is not a valid fragment file: %s", filepath.Base(fn), errMsg)
	}

	return Fragment{string(data), parts[0]}, parts[1], nil
}
