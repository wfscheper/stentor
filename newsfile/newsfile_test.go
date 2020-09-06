package newsfile

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"pgregory.net/rapid"
)

func TestWriteFragments(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	tmpdir, err := ioutil.TempDir("", "stentor-")
	must.NoError(err)

	wd, err := os.Getwd()
	must.NoError(err)

	must.NoError(os.Chdir(tmpdir))
	defer os.Chdir(wd) // nolint:errcheck // defer func

	fn := filepath.Join(tmpdir, "file")
	err = ioutil.WriteFile(fn, []byte("some text\n\n.. stentor output starts\n\nsome more text\n"), 0600)
	must.NoError(err)

	if err := WriteFragments(fn, []byte(startCommentRST), []byte("\nadded data")); is.NoError(err) {
		data, err := ioutil.ReadFile(fn)
		must.NoError(err)
		is.Equal("some text\n\n.. stentor output starts\n\nadded data\nsome more text\n", string(data))
	}
}

func Test_copyIntoFile(t *testing.T) {
	t.Parallel()

	t.Run("comment exists", rapid.MakeCheck(func(t *rapid.T) {
		startComment := rapid.SampledFrom([]string{startCommentMarkdown, startCommentRST}).
			Draw(t, "startComment").(string)
		header := rapid.StringN(512, 1024, -1).Draw(t, "header").(string)
		trailer := rapid.StringN(512, 1024, -1).Draw(t, "trailer").(string)
		src := bytes.NewBufferString(header + startComment + trailer)
		dst := &bytes.Buffer{}
		data := rapid.String().Draw(t, "data").(string)

		if err := copyIntoFile(dst, src, []byte(startComment), []byte(data)); assert.NoError(t, err) {
			assert.Equal(t, header+startComment+data+trailer, dst.String())
		}
	}))

	t.Run("no comment exists", rapid.MakeCheck(func(t *rapid.T) {
		startComment := rapid.SampledFrom([]string{startCommentMarkdown, startCommentRST}).
			Draw(t, "startComment").(string)
		header := rapid.StringN(512, 1024, -1).Draw(t, "header").(string)
		trailer := rapid.StringN(512, 1024, -1).Draw(t, "trailer").(string)
		src := bytes.NewBufferString(header + trailer)
		dst := &bytes.Buffer{}
		data := rapid.String().Draw(t, "data").(string)

		err := copyIntoFile(dst, src, []byte(startComment), []byte(data))
		assert.EqualError(t, err, "no start comment found")
	}))
}
