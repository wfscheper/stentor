package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"testing"
)

func TestStentor_displayVersion(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "stentor-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	out := &bytes.Buffer{}
	s := New(tmpdir, []string{}, []string{}, ioutil.Discard, out)
	if got, want := s.displayVersion(), 0; got != want {
		t.Errorf("displayVersion returned %d, want %d", got, want)
	}

	versionInfo := fmt.Sprintf(`stentor
  version     : dev
  build date  : none
  git hash    : none
  go version  : %s
  go compiler : %s
  platform    : %s/%s
`, runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH)

	if got, want := out.String(), versionInfo; got != want {
		t.Errorf("dipslayVersion wrote\n%s\nwant\n%s", got, want)
	}
}
