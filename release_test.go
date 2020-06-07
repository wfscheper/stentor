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

package stentor

import (
	"bytes"
	"testing"
	"time"

	"github.com/ianbruene/go-difflib/difflib"
	"github.com/wfscheper/stentor/internal/templates"
)

func TestSectionTemplate(t *testing.T) {
	tests := []struct {
		name        string
		releaseFunc func(string) *Release
		want        string
	}{
		{
			"github-markdown-section",
			NewMarkdownRelease,
			"## [v0.2.0](https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0) - 2020-01-02\n" +
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
				"----\n",
		},
		{
			"github-rst-section",
			NewRSTRelease,
			"`v0.2.0 <https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0>`_ - 2020-01-02\n" +
				"----------------------\n" +
				"\n" +
				"Features\n" +
				"^^^^^^^^\n" +
				"\n" +
				"- The foo feature.\n" +
				"\n" +
				"  This is an awesome feature.\n" +
				"  `#1 <https://github.com/myname/myrepo/issues/1>`_\n" +
				"\n" +
				"\n" +
				"Bug Fixes\n" +
				"^^^^^^^^^\n" +
				"\n" +
				"- Fix the bug in foo.\n" +
				"  `#2 <https://github.com/myname/myrepo/issues/2>`_\n" +
				"\n" +
				"\n" +
				"Always Show\n" +
				"^^^^^^^^^^^\n" +
				"\n" +
				"No significant changes.\n" +
				"\n" +
				"\n" +
				"----\n",
		},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.releaseFunc("myname/myrepo")
			r.Date = time.Date(2020, 1, 2, 3, 4, 5, 6, time.UTC)
			r.PreviousVersion = "v0.1.0"
			r.Version = "v0.2.0"
			r.Sections = []Section{
				{
					Fragments: []Fragment{
						{
							Issue: "1",
							Text:  "The foo feature.\n\nThis is an awesome feature.",
						},
					},
					Title: "Features",
				},
				{
					Fragments: []Fragment{
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
			if err != nil {
				t.Fatal(err)
			}

			buf := &bytes.Buffer{}
			if err := tmp.Execute(buf, r); err != nil {
				t.Fatalf("tmp.Execute returned an error: %v", err)
			}

			if got, want := buf.String(), tt.want; got != want {
				diff := difflib.UnifiedDiff{
					A:        difflib.SplitLines(want),
					B:        difflib.SplitLines(got),
					FromFile: "want",
					ToFile:   "got",
					Context:  3,
				}
				text, err := difflib.GetUnifiedDiffString(diff)
				if err != nil {
					t.Fatal(err)
				}
				t.Errorf("tmp.Execute returned:\n%s", text)
			}
		})
	}
}

func TestNewRelease(t *testing.T) {
	repo := "myname/myrepo"
	r := NewRelease(repo)
	if got, want := r.Repository, repo; got != want {
		t.Errorf("r.Repository == %q, want %q", got, want)
	}
}

func TestNewMarkdownRelease(t *testing.T) {
	r := NewMarkdownRelease("myname/myrepo")
	if got, want := r.Header, "##"; got != want {
		t.Errorf("r.Header == %q, want %q", got, want)
	}
	if got, want := r.SectionHeader, "###"; got != want {
		t.Errorf("r.SectionHeader == %q, want %q", got, want)
	}
}

func TestNewRSTRelease(t *testing.T) {
	r := NewRSTRelease("myname/myrepo")
	if got, want := r.Header, "-"; got != want {
		t.Errorf("r.Header == %q, want %q", got, want)
	}
	if got, want := r.SectionHeader, "^"; got != want {
		t.Errorf("r.SectionHeader == %q, want %q", got, want)
	}
}
