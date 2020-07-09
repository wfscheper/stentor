package main

import (
	"bytes"
	"io/ioutil"
	"os"
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
	s.displayVersion()

	if got, want := out.String(), "stentor dev built from unknown on unknown\n"; got != want {
		t.Errorf("dipslayVersion wrote %q, want %q", got, want)
	}
}
