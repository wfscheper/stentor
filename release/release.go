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
	"time"

	"github.com/wfscheper/stentor/section"
)

// Release represents the data used to generate a release entry in a stentor-managed news file.
type Release struct {
	// Date is the date of the release.
	Date time.Time
	// Header is the markup character used when writing the release header.
	Header string
	// PreviousVersion is the version before this release.
	PreviousVersion string
	// Repository is the user/repo portion of the git repository URL.
	Repository string
	// SectionHeader is the markup character used when writing a section header.
	SectionHeader string
	// Sections is the list of change types in this release.
	Sections []section.Section
	// Version is the version of this release.
	Version string
}

// NewRelease returns a simple Release.
//
// The caller is responsible for defining the Header and SectionHeader.
func New(repo string) *Release {
	return &Release{
		Date:       time.Now().UTC(),
		Repository: repo,
	}
}

// NewMarkdownRelease returns a Release with markdown style Header and SectionHeader.
func NewMarkdown(repo string) *Release {
	r := New(repo)
	r.Header = "##"
	r.SectionHeader = "###"
	return r
}

// NewRSTRelease returns a Release with reStructuredText style Header and SectionHeader.
func NewRST(repo string) *Release {
	r := New(repo)
	r.Header = "-"
	r.SectionHeader = "^"
	return r
}
