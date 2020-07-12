package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStentor_displayVersion(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "stentor-")
	require.NoError(t, err)
	defer os.RemoveAll(tmpdir)

	out := &bytes.Buffer{}
	s := New(tmpdir, []string{}, []string{}, ioutil.Discard, out)
	s.displayVersion()

	assert.Equal(t, "stentor dev built from unknown on unknown\n", out.String())
}
