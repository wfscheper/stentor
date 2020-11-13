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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		want Fragment
	}{
		{"ticket.section.md", Fragment{"section", "ticket", "contents"}},
		{"ticket.section.extra-bit.md", Fragment{"section", "ticket", "contents"}},
		{"ticket.section.several.extra.bits.md", Fragment{"section", "ticket", "contents"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpdir, err := ioutil.TempDir("", "stentor-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpdir)

			fn := filepath.Join(tmpdir, tt.name)
			err = ioutil.WriteFile(fn, []byte(`contents`), 0600)
			require.NoError(t, err)

			if got, err := Parse(fn); assert.NoError(t, err) {
				assert.Equal(t, tt.want, *got)
			}
		})
	}
}

func TestNew_error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"foo", "not a valid fragment file: not enough parts"},
		{"foo.md", "not a valid fragment file: not enough parts"},
		{".section.md", "not a valid fragment file: empty issue"},
		{"ticket..md", "not a valid fragment file: empty section"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "stentor-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpdir)

			fn := filepath.Join(tmpdir, tt.name)
			err = ioutil.WriteFile(fn, []byte(`contents`), 0600)
			require.NoError(t, err)

			_, err = Parse(fn)
			assert.EqualError(t, err, tt.want)
		})
	}
}
