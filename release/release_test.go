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

package release

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wfscheper/stentor/fragment"
	"github.com/wfscheper/stentor/internal/templates"
	"github.com/wfscheper/stentor/section"
)

func TestSectionTemplate(t *testing.T) {
	tests := []struct {
		name        string
		releaseFunc func(string, string, string) (*Release, error)
		want        string
	}{
		{
			"github-markdown-section",
			newMarkdown,
			"## [v0.2.0] - 2020-01-02\n" +
				"\n" +
				"### Features\n" +
				"\n" +
				"- The foo feature.\n" +
				"\n" +
				"  This is an awesome feature.\n" +
				"  [#1](https://host/myname/myrepo/issues/1)\n" +
				"\n" +
				"\n" +
				"### Bug Fixes\n" +
				"\n" +
				"- Fix the bug in foo.\n" +
				"  [#2](https://host/myname/myrepo/issues/2)\n" +
				"- Multiple other things.\n" +
				"\n" +
				"\n" +
				"### Always Show\n" +
				"\n" +
				"No significant changes.\n" +
				"\n" +
				"\n" +
				"[v0.2.0]: https://host/myname/myrepo/compare/v0.1.0...v0.2.0\n" +
				"\n" +
				"\n" +
				"----\n" +
				"\n",
		},
		{
			"github-rst-section",
			newRST,
			"`v0.2.0`_ - 2020-01-02\n" +
				"======================\n" +
				"\n" +
				"Features\n" +
				"--------\n" +
				"\n" +
				"- The foo feature.\n" +
				"\n" +
				"  This is an awesome feature.\n" +
				"  `#1 <https://host/myname/myrepo/issues/1>`_\n" +
				"\n" +
				"\n" +
				"Bug Fixes\n" +
				"---------\n" +
				"\n" +
				"- Fix the bug in foo.\n" +
				"  `#2 <https://host/myname/myrepo/issues/2>`_\n" +
				"- Multiple other things.\n" +
				"\n" +
				"\n" +
				"Always Show\n" +
				"-----------\n" +
				"\n" +
				"No significant changes.\n" +
				"\n" +
				"\n" +
				".. _v0.2.0: https://host/myname/myrepo/compare/v0.1.0...v0.2.0\n" +
				"\n" +
				"\n" +
				"----\n" +
				"\n",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			r, err := tt.releaseFunc("https://host/myname/myrepo", "v0.2.0", "v0.1.0")
			require.NoError(t, err)

			// assing a fixed date
			r.Date = time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
			r.Sections = []section.Section{
				{
					Fragments: []fragment.Fragment{
						{
							Issue: "1",
							Text:  "The foo feature.\n\nThis is an awesome feature.",
						},
					},
					Title: "Features",
				},
				{
					Fragments: []fragment.Fragment{
						{
							Issue: "2",
							Text:  "Fix the bug in foo.",
						},
						{
							Text: "Multiple other things.",
						},
					},
					Title: "Bug Fixes",
				},
				{
					ShowAlways: true,
					Title:      "Always Show",
				},
			}

			tmp, err := templates.New(tt.name)
			require.NoError(t, err)

			buf := &bytes.Buffer{}
			require.NoError(t, tmp.Execute(buf, r))

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func Test_newRelease(t *testing.T) {
	tests := []struct {
		repo, want, wantError string
	}{
		{"http://myhost/myname/myrepo", "http://myhost/myname/myrepo", ""},
		{"https://myhost/myname/myrepo", "https://myhost/myname/myrepo", ""},
		{"https:///myname/myrepo", "https:///myname/myrepo", ""},
		{"https://myhost/myrepo", "https://myhost/myrepo", ""},
		{"https://myhost/myname/myrepo?branch=foo", "https://myhost/myname/myrepo", ""},
		{"https://myhost/myname/myrepo#afragment", "https://myhost/myname/myrepo", ""},
		{"https://myhost/myname/myrepo?branch=foo#afragment", "https://myhost/myname/myrepo", ""},
		{"myhost/myname/myrepo", "", "invalid URL: no scheme"},
		{"file://myhost/myname/myrepo", "", "invalid URL: only http or https schemes"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.repo, func(t *testing.T) {
			t.Parallel()

			if r, err := newRelease(tt.repo, "v0.2.0", "v0.1.0"); err != nil {
				assert.EqualError(t, err, tt.wantError)
			} else {
				assert.Equal(t, tt.want, r.Repository)
			}
		})
	}
}

func Test_newMarkdown(t *testing.T) {
	if r, err := newMarkdown("https://myhost/myname/myrepo", "v0.2.0", "v0.1.0"); assert.NoError(t, err) {
		assert.Equal(t, "##", r.Header)
		assert.Equal(t, "###", r.SectionHeader)
	}
}

func Test_newRST(t *testing.T) {
	if r, err := newRST("https://myhost/myname/myrepo", "v0.2.0", "v0.1.0"); assert.NoError(t, err) {
		assert.Equal(t, "=", r.Header)
		assert.Equal(t, "-", r.SectionHeader)
	}
}
