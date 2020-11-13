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

// Package release provides the Release structure.
package release

import (
	"errors"
	"net/url"
	"time"

	"github.com/wfscheper/stentor"
	"github.com/wfscheper/stentor/config"
	"github.com/wfscheper/stentor/fragment"
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
	// Repository is the URL of the project repository.
	Repository string
	// SectionHeader is the markup character used when writing a section header.
	SectionHeader string
	// Sections is the list of change types in this release.
	Sections []section.Section
	// Version is the version of this release.
	Version string
}

// New returns a Release.
//
// The repo should be a parsable URL.
func New(repo, markup, version, previousVersion string) (*Release, error) {
	switch markup {
	case stentor.MarkupMD:
		return newMarkdown(repo, version, previousVersion)
	case stentor.MarkupRST:
		return newRST(repo, version, previousVersion)
	default:
		return newRelease(repo, version, previousVersion)
	}
}

// SetSections populates the release's sections.
func (r *Release) SetSections(sections []config.Section, fragments []fragment.Fragment) {
	sectionMap := map[string]section.Section{}
	for _, fragment := range fragments {
		s := sectionMap[fragment.Section]
		s.Fragments = append(s.Fragments, fragment)
		sectionMap[fragment.Section] = s
	}

	for _, cfg := range sections {
		if s, ok := sectionMap[cfg.ShortName]; ok {
			if cfg.ShowAlways != nil {
				s.ShowAlways = *cfg.ShowAlways
			}
			s.Title = cfg.Name
			r.Sections = append(r.Sections, s)
		} else if cfg.ShowAlways != nil && *cfg.ShowAlways {
			r.Sections = append(r.Sections, section.Section{
				ShowAlways: *cfg.ShowAlways,
				Title:      cfg.Name,
			})
		}
	}
}

func newRelease(repo, version, previousVersion string) (*Release, error) {
	u, err := url.Parse(repo)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
	case "":
		return nil, errors.New("invalid URL: no scheme")
	default:
		return nil, errors.New("invalid URL: only http or https schemes")
	}

	u.Fragment = ""
	u.RawQuery = ""

	return &Release{
		Date:            time.Now().UTC(),
		Repository:      u.String(),
		PreviousVersion: previousVersion,
		Version:         version,
	}, nil
}

// NewMarkdownRelease returns a Release with markdown style Header and SectionHeader.
func newMarkdown(repo, version, previousVersion string) (*Release, error) {
	r, err := newRelease(repo, version, previousVersion)
	if err != nil {
		return nil, err
	}

	r.Header = "##"
	r.SectionHeader = "###"
	return r, nil
}

// NewRSTRelease returns a Release with reStructuredText style Header and SectionHeader.
func newRST(repo, version, previousVersion string) (*Release, error) {
	r, err := newRelease(repo, version, previousVersion)
	if err != nil {
		return nil, err
	}

	r.Header = "="
	r.SectionHeader = "-"
	return r, nil
}
