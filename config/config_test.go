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

package config

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
		Sections: []Section{
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
	t.Parallel()

	type parseFunc func([]byte) (Config, error)
	tests := []parseFunc{parseConfig, ParseBytes}

	defaultConfig := Config{
		FragmentDir: ".stentor.d",
		Hosting:     "github",
		Markup:      "markdown",
		NewsFile:    "CHANGELOG.md",
		Sections: []Section{
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

	for _, tf := range tests {
		tf := tf

		t.Run("empty config", func(t *testing.T) {
			t.Parallel()

			if c, err := tf([]byte("")); assert.NoError(t, err) {
				assert.Equal(t, defaultConfig, c)
			}
		})

		t.Run("bad toml", func(t *testing.T) {
			t.Parallel()

			// bad toml
			y := []byte(`
[stentor]
foo = "bar"
`)
			_, err := tf(y)
			assert.EqualError(t, err, "undecoded keys: [\"stentor.foo\"]")
		})
	}
}

func Test_validateConfig(t *testing.T) {
	t.Parallel()

	t.Run("invalid hosting", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    rapid.String().Draw(t, "hosting"),
			Markup:     genMarkup().Draw(t, "markup"),
			Repository: genRepository().Draw(t, "repository"),
		}
		assert.EqualError(t, ValidateConfig(c), ErrBadHosting.Error())
	}))

	t.Run("invalid markup", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    genHosting().Draw(t, "hosting"),
			Markup:     rapid.String().Draw(t, "markup"),
			Repository: genRepository().Draw(t, "repository"),
		}
		assert.EqualError(t, ValidateConfig(c), ErrBadMarkup.Error())
	}))

	t.Run("invalid repository", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    genHosting().Draw(t, "hosting"),
			Markup:     genMarkup().Draw(t, "markup"),
			Repository: rapid.SampledFrom([]string{"file", ""}).Draw(t, "repository"),
		}
		if err := ValidateConfig(c); err != nil {
			if c.Repository == "" {
				if err.Error() != ErrMissingRepository.Error() {
					t.Errorf("expected error %v, got %v", ErrMissingRepository, err)
				}
			} else {
				if !strings.HasPrefix(err.Error(), "invalid repository: ") {
					t.Errorf("expected invalid repository error, got %v", err)
				}
			}
		}
	}))

	t.Run("no sections", rapid.MakeCheck(func(t *rapid.T) {
		c := Config{
			Hosting:    genHosting().Draw(t, "hosting"),
			Markup:     genMarkup().Draw(t, "markup"),
			Repository: genRepository().Draw(t, "repository"),
		}
		assert.EqualError(t, ValidateConfig(c), ErrBadSections.Error())
	}))
}

func genHosting() *rapid.Generator[string]    { return rapid.SampledFrom([]string{"github", "gitlab"}) }
func genMarkup() *rapid.Generator[string]     { return rapid.SampledFrom([]string{"markdown", "rst"}) }
func genRepository() *rapid.Generator[string] { return rapid.Just("https://host/name/repo") }
