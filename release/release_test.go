// Copyright © 2020 The Stentor Authors
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
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name        string
		releaseFunc func(string, string, string) Release
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
				"  [#1](https://github.com/myname/myrepo/issues/1)\n" +
				"\n" +
				"\n" +
				"### Bug Fixes\n" +
				"\n" +
				"- Fix the bug in foo.\n" +
				"  [#2](https://github.com/myname/myrepo/issues/2)\n" +
				"\n" +
				"\n" +
				"### Always Show\n" +
				"\n" +
				"No significant changes.\n" +
				"\n" +
				"\n" +
				"[v0.2.0]: https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0\n" +
				"\n" +
				"\n" +
				"----\n" +
				"\n" +
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
				"  `#1 <https://github.com/myname/myrepo/issues/1>`_\n" +
				"\n" +
				"\n" +
				"Bug Fixes\n" +
				"---------\n" +
				"\n" +
				"- Fix the bug in foo.\n" +
				"  `#2 <https://github.com/myname/myrepo/issues/2>`_\n" +
				"\n" +
				"\n" +
				"Always Show\n" +
				"-----------\n" +
				"\n" +
				"No significant changes.\n" +
				"\n" +
				"\n" +
				".. _v0.2.0: https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0\n" +
				"\n" +
				"\n" +
				"----\n" +
				"\n" +
				"\n",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.releaseFunc("myname/myrepo", "v0.2.0", "v0.1.0")
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
					},
					Title: "Bug Fixes",
				},
				{
					ShowAlways: true,
					Title:      "Always Show",
				},
			}

			tmp, err := templates.New(tt.name)
			require.NoError(err)

			buf := &bytes.Buffer{}
			require.NoError(tmp.Execute(buf, r))

			assert.Equal(tt.want, buf.String())
		})
	}
}

func Test_newRelease(t *testing.T) {
	r := newRelease("myname/myrepo", "v0.2.0", "v0.1.0")
	assert.Equal(t, "myname/myrepo", r.Repository)
}

func Test_newMarkdown(t *testing.T) {
	r := newMarkdown("myname/myrepo", "v0.2.0", "v0.1.0")
	assert.Equal(t, "##", r.Header)
	assert.Equal(t, "###", r.SectionHeader)
}

func Test_newRST(t *testing.T) {
	r := newRST("myname/myrepo", "v0.2.0", "v0.1.0")
	assert.Equal(t, "=", r.Header)
	assert.Equal(t, "-", r.SectionHeader)
}
