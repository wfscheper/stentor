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

package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"github-markdown-section"},
		{"github-rst-section"},
		{"gitlab-markdown-section"},
		{"gitlab-rst-section"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.name)
			assert.NoError(t, err)
		})
	}
}

func TestNew_error(t *testing.T) {
	_, err := New("notexist")
	require.Error(t, err)
}

const customTemplate = `Custom template.

{{ "repeat" | repeat 2 }}

{{ "The next two lines\nshould be indented\ntwo spaces." | indent 2 }}

{{ sum 2 3 }}
`

func TestParse(t *testing.T) {
	tmp := t.TempDir()

	fn := filepath.Join(tmp, "test.template")
	require.NoError(t, os.WriteFile(fn, []byte(customTemplate), 0600))

	tmpl, err := Parse(fn)
	require.NoError(t, err)

	buf := &bytes.Buffer{}
	require.NoError(t, tmpl.Execute(buf, "ignore me"))

	want := `Custom template.

repeatrepeat

The next two lines
  should be indented
  two spaces.

5
`
	require.Equal(t, want, buf.String())
}

func TestParse_error(t *testing.T) {
	_, err := Parse("notexist")
	require.Error(t, err)
}
