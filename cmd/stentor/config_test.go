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

package main

import (
	"strings"
	"testing"

	"github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"pgregory.net/rapid"
)

func TestConfig_marshal(t *testing.T) {
	is := assert.New(t)

	if data, err := toml.Marshal(Config{}); is.NoError(err) {
		is.Equal("", string(data))
	}

	u := Config{
		FragmentDir:     "fragments",
		HeaderTemplate:  "header",
		Hosting:         "hosting",
		Markup:          "markup",
		NewsFile:        "news",
		Repository:      "repo",
		SectionTemplate: "section",
		Sections: []SectionConfig{
			{
				Name:       "Name",
				ShortName:  "name",
				ShowAlways: func(b bool) *bool { return &b }(true),
			},
		},
	}

	wantTOML := `
# Stentor configuration
[stentor]
  fragment_dir = "fragments"
  header_template = "header"
  hosting = "hosting"
  markup = "markup"
  news_file = "news"
  repository = "repo"
  section_template = "section"

  [[stentor.sections]]
    name = "Name"
    short_name = "name"
    show_always = true
`

	var v Config
	if err := toml.Unmarshal([]byte(wantTOML), &tomlConfig{&v}); is.NoError(err) {
		is.Equal(u, v)
	}

	// marshal u and compare to input
	if data, err := toml.Marshal(tomlConfig{&u}); is.NoError(err) {
		is.Equal(wantTOML, string(data))
	}
}

func Test_parseConfig(t *testing.T) {
	defaultConfig := Config{
		FragmentDir: ".stentor.d",
		Hosting:     "github",
		Markup:      "markdown",
		NewsFile:    "CHANGELOG.md",
		Sections: []SectionConfig{
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
		},
	}

	is := assert.New(t)

	if c, err := parseConfig([]byte("")); is.NoError(err) {
		is.Equal(defaultConfig, c)
	}

	// bad toml
	y := []byte(`
[stentor]
foo = "bar"
`)
	_, err := parseConfig(y)
	is.EqualError(err, "undecoded keys: [\"stentor.foo\"]")
}

func Test_validateConfig(t *testing.T) {
	t.Parallel()

	genHosting := rapid.SampledFrom([]string{"github", "gitlab"})
	genMarkup := rapid.SampledFrom([]string{"markdown", "rst"})
	genRepository := rapid.Just("foo/bar")
	is := assert.New(t)

	t.Run("invalid hosting", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Repository: genRepository.Draw(t, "repository").(string),
			Hosting:    rapid.String().Draw(t, "hosting").(string),
			Markup:     genMarkup.Draw(t, "markup").(string),
		}
		is.EqualError(ValidateConfig(c), ErrBadHosting.Error())
	}))

	t.Run("invalid markup", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    genHosting.Draw(t, "hosting").(string),
			Markup:     rapid.String().Draw(t, "markup").(string),
			Repository: genRepository.Draw(t, "repository").(string),
		}
		is.EqualError(ValidateConfig(c), ErrBadMarkup.Error())
	}))

	t.Run("invalid repository", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Repository: rapid.String().Filter(func(s string) bool {
				return strings.Count(s, "/") != 1
			}).Draw(t, "repository").(string),
			Hosting: genHosting.Draw(t, "hosting").(string),
			Markup:  genMarkup.Draw(t, "markup").(string),
		}
		want := ErrBadRepository
		if c.Repository == "" {
			want = ErrMissingRepository
		}
		is.EqualError(ValidateConfig(c), want.Error())
	}))

	t.Run("no sections", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    genHosting.Draw(t, "hosting").(string),
			Markup:     genMarkup.Draw(t, "markup").(string),
			Repository: genRepository.Draw(t, "repository").(string),
		}
		is.EqualError(ValidateConfig(c), ErrBadSections.Error())
	}))
}
