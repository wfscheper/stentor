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
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		want    Fragment
		wantErr string
	}{
		{
			name: "ticket.section.md",
			want: Fragment{"section", "ticket", "contents"},
		},
		{
			name: "ticket.section.extra-bit.md",
			want: Fragment{"section", "ticket", "contents"},
		},
		{
			name: "ticket.section.several.extra.bits.md",
			want: Fragment{"section", "ticket", "contents"},
		},
		{
			name:    "foo",
			wantErr: "not a valid fragment file: not enough parts",
		},
		{
			name:    "foo.md",
			wantErr: "not a valid fragment file: not enough parts",
		},
		{
			name:    ".section.md",
			wantErr: "not a valid fragment file: empty issue",
		},
		{
			name:    "ticket..md",
			wantErr: "not a valid fragment file: empty section",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpdir := t.TempDir()

			fn := filepath.Join(tmpdir, tt.name)
			require.NoError(t, ioutil.WriteFile(fn, []byte(`contents`), 0600))

			got, err := Parse(fn)
			if tt.wantErr == "" {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, *got)
				}
			} else {
				assert.EqualError(t, err, tt.wantErr)
			}
		})
	}
}
