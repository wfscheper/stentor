package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStentor_displayVersion(t *testing.T) {
	tmpdir := t.TempDir()

	out := &bytes.Buffer{}
	s := New(tmpdir, []string{}, []string{}, io.Discard, out)
	s.displayVersion()

	assert.Equal(t, "stentor dev built from unknown on unknown\n", out.String())
}
