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
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wfscheper/stentor"
	"pgregory.net/rapid"
)

func TestWriteFragments(t *testing.T) {
	is := assert.New(t)
	must := require.New(t)

	tmpdir := t.TempDir()

	wd, err := os.Getwd()
	must.NoError(err)

	must.NoError(os.Chdir(tmpdir))
	defer os.Chdir(wd) // nolint:errcheck // defer func

	fn := filepath.Join(tmpdir, "file")
	must.NoError(os.WriteFile(fn, []byte("some text\n\n.. stentor output starts\n\nsome more text\n"), 0600))

	if err := WriteRelease(fn, stentor.CommentRST, []byte("added data\n"), true); is.NoError(err) {
		data, err := os.ReadFile(fn)
		must.NoError(err)
		is.Equal("some text\n\n.. stentor output starts\nadded data\n\nsome more text\n", string(data))
	}
}

func TestWriteFragments_no_comment(t *testing.T) {
	tmpdir := t.TempDir()

	wd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.Chdir(tmpdir))
	defer os.Chdir(wd) // nolint:errcheck // defer func

	fn := filepath.Join(tmpdir, "file")
	require.NoError(t, os.WriteFile(fn, []byte("some text\nsome more text\n"), 0600))
	require.EqualError(t, WriteRelease(fn, stentor.CommentMD, []byte("added data\n"), true), "no start comment found")
}

func Test_copyIntoFile(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		nt := newsfileGen().Draw(t, "newsfile").(newsfileTest)

		var src io.Reader
		if nt.includeComment {
			src = bytes.NewBufferString(nt.header + nt.startComment + nt.trailer)
		} else {
			src = bytes.NewBufferString(nt.header + nt.trailer)
		}

		dst := &bytes.Buffer{}
		err := copyIntoFile(dst, src, []byte(nt.startComment), []byte(nt.data), nt.keepHeader)
		if nt.includeComment {
			if err != nil {
				t.Fatalf("copyIntoFile() returned an error: %v", err)
			}

			want := nt.data + nt.trailer
			if nt.keepHeader {
				want = nt.header + nt.startComment + want
			}
			if got := dst.String(); got != want {
				t.Errorf("copyIntoFile() wrote\n%s\n\nwant\n%s", got, want)
			}
		} else if err == nil || err.Error() != "no start comment found" {
			t.Fatalf("copyIntoFile() returned an error, %v, wanted, %q", err, "no start comment found")
		}
	})
}

type newsfileTest struct {
	startComment, header, trailer, data string
	keepHeader, includeComment          bool
}

func newsfileGen() *rapid.Generator {
	return rapid.Custom(func(t *rapid.T) newsfileTest {
		return newsfileTest{
			startComment: rapid.SampledFrom([]string{
				stentor.CommentMD,
				stentor.CommentRST,
			}).Draw(t, "startComment").(string),
			// we pick these sizes to force the start comment out past a single read
			header:         rapid.StringN(512, 1024, -1).Draw(t, "header").(string),
			trailer:        rapid.StringN(512, 1024, -1).Draw(t, "trailer").(string),
			data:           rapid.String().Draw(t, "data").(string),
			keepHeader:     rapid.Bool().Draw(t, "keepHeader").(bool),
			includeComment: rapid.Bool().Draw(t, "includeComment").(bool),
		}
	})
}
