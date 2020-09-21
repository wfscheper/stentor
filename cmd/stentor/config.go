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

package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/wfscheper/stentor"
)

const (
	DefaultConfigDir = ".stentor.d"
	HostingGithub    = "github"
	HostingGitlab    = "gitlab"
	MarkupMarkdown   = "markdown"
	MarkupRST        = "rst"
)

var (
	ErrBadHosting        = errors.New("hosting must be one of 'github' or 'gitlab'")
	ErrBadMarkup         = errors.New("markup must be one of 'markdown' or 'rst'")
	ErrBadRepository     = errors.New("repository must be in the format <user name>/<repository name>")
	ErrBadSections       = errors.New("must define at least one section")
	ErrMissingRepository = errors.New("repository is required")

	defaultSectionConfig = []SectionConfig{
		{
			Name:      "Security",
			ShortName: "security",
		},
		{
			Name:      "Deprecated",
			ShortName: "deprecate",
		},
		{
			Name:      "Removed",
			ShortName: "remove",
		},
		{
			Name:      "Changed",
			ShortName: "change",
		},
		{
			Name:      "Added",
			ShortName: "feature",
		},
		{
			Name:      "Fixed",
			ShortName: "fix",
		},
	}
)

type tomlConfig struct {
	Stentor *Config `toml:"stentor" comment:"Stentor configuration"`
}

// Config represents the project's configuration for stentor.
type Config struct {
	// Repository is the name of your repository in <username>/<repo name> format.
	Repository string `toml:"repository,omitempty"`
	// FragmentDir is the path to the directory holding the project's news fragments.
	// Defaults to '.stentor.d'.
	FragmentDir string `toml:"fragment_dir,omitempty" yaml:"fragment_dir,omitempty"`
	// Hosting is the source repository host.
	// When Markup is set to markdown, this also determines the markdown flavor.
	// Currently, github and gitlab are supported.
	// Defaults to github.
	Hosting string `toml:"hosting,omitempty"`
	// Markup sets the format of your changelog.
	// Currently, markdown and rst (ReStructuredText) are supported.
	// Defaults to markdown
	Markup string `toml:"markup,omitempty"`
	// Sections define the different news sections.
	// Sections will be listed in the order in which they are defined here.
	Sections []SectionConfig `toml:"sections,omitempty"`
	// HeaderTemplate is the name of the template used to render the header of the news file.
	HeaderTemplate string `toml:"header_template,omitempty"`
	// SectionTemplate is the name of the template used to render the individual sections of the news file.
	SectionTemplate string `toml:"section_template,omitempty"`
	// NewsFile is the name of the file to update
	NewsFile string `toml:"news_file,omitempty"`
}

// ParseBytes parses bytes data into a Config.
func ParseBytes(data []byte) (Config, error) {
	c, err := parseConfig(data)
	if err != nil {
		return Config{}, err
	}
	return c, nil
}

func parseConfig(data []byte) (Config, error) {
	var c Config
	if err := toml.NewDecoder(bytes.NewReader(data)).Strict(true).Decode(&tomlConfig{&c}); err != nil {
		return Config{}, err
	}

	if c.FragmentDir == "" {
		c.FragmentDir = DefaultConfigDir
	}

	if c.Hosting == "" {
		c.Hosting = stentor.HostingGithub
	}

	if c.Markup == "" {
		c.Markup = stentor.MarkupMD
	}

	if c.NewsFile == "" {
		c.NewsFile = "CHANGELOG"
		switch c.Markup {
		case stentor.MarkupMD:
			c.NewsFile += ".md"
		case stentor.MarkupRST:
			c.NewsFile += ".rst"
		default:
			return Config{}, fmt.Errorf("unrecognized markup: %s", c.Markup)
		}
	}

	if len(c.Sections) == 0 {
		c.Sections = defaultSectionConfig
	}

	return c, nil
}

func ValidateConfig(c Config) error {
	// repository must be non-empty
	if c.Repository == "" {
		return ErrMissingRepository
	}
	// and must contain a single / that isn't the first or last character
	if strings.Count(c.Repository, "/") != 1 || c.Repository[0] == '/' || c.Repository[len(c.Repository)-1] == '/' {
		return ErrBadRepository
	}
	// hosting must be github or gitlab
	if c.Hosting != HostingGithub && c.Hosting != HostingGitlab {
		return ErrBadHosting
	}
	// markup must be markdown or rst
	if c.Markup != MarkupMarkdown && c.Markup != MarkupRST {
		return ErrBadMarkup
	}
	// must have at least one section
	if len(c.Sections) < 1 {
		return ErrBadSections
	}
	return nil
}

func (c Config) FragmentFiles() ([]string, error) {
	var glob string
	switch c.Markup {
	case stentor.MarkupMD:
		glob = "*.md"
	case stentor.MarkupRST:
		glob = "*.rst"
	default:
		return nil, fmt.Errorf("unknown markup %s", c.Markup)
	}

	return filepath.Glob(filepath.Join(c.FragmentDir, glob))
}

func (c Config) StartComment() string {
	switch c.Markup {
	case stentor.MarkupMD:
		return stentor.CommentMD
	case stentor.MarkupRST:
		return stentor.CommentRST
	default:
		return ""
	}
}

// Section represents a group of news items in a release.
type SectionConfig struct {
	// Name of the section.
	Name string `toml:"name,omitempty"`
	// ShorName is the string used in a fragment file to indicate what section the fragment is for.
	ShortName string `toml:"short_name,omitempty" yaml:"short_name,omitempty"`
	// ShowAlways is a boolean indicating whether to show the section even if there are no news items.
	// This is a pointer to that we can use omitempty, and still render false values.
	ShowAlways *bool `toml:"show_always,omitempty" yaml:"show_always,omitempty"`
}
