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

package newsfile

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wfscheper/stentor"
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

	if err := WriteFragments(fn, stentor.CommentRST, []byte("\nadded data")); is.NoError(err) {
		data, err := ioutil.ReadFile(fn)
		must.NoError(err)
		is.Equal("some text\n\n.. stentor output starts\n\nadded data\nsome more text\n", string(data))
	}
}

func Test_copyIntoFile(t *testing.T) {
	t.Parallel()

	t.Run("comment exists", rapid.MakeCheck(func(t *rapid.T) {
		startComment := rapid.SampledFrom([]string{
			stentor.CommentMD,
			stentor.CommentRST,
		}).Draw(t, "startComment").(string)
		header := rapid.StringN(512, 1024, -1).Draw(t, "header").(string)
		trailer := rapid.StringN(512, 1024, -1).Draw(t, "trailer").(string)

		srcString := header + startComment + trailer
		data := rapid.String().Draw(t, "data").(string)

		src := bytes.NewBufferString(srcString)
		dst := &bytes.Buffer{}

		err := copyIntoFile(dst, src, startComment, []byte(data))
		if err != nil {
			assert.True(t, !strings.Contains(srcString, "\n"+startComment+"\n"))
		} else {
			assert.Equal(t, header+startComment+data+trailer, dst.String())
		}
	}))

	t.Run("no comment exists", rapid.MakeCheck(func(t *rapid.T) {
		startComment := rapid.SampledFrom([]string{stentor.CommentMD, stentor.CommentRST}).
			Draw(t, "startComment").(string)
		header := rapid.StringN(512, 1024, -1).Draw(t, "header").(string)
		trailer := rapid.StringN(512, 1024, -1).Draw(t, "trailer").(string)
		src := bytes.NewBufferString(header + trailer)
		dst := &bytes.Buffer{}
		data := rapid.String().Draw(t, "data").(string)

		err := copyIntoFile(dst, src, startComment, []byte(data))
		assert.EqualError(t, err, "no start comment found")
	}))
}
