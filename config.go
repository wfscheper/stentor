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
	"errors"
	"strings"

	"gopkg.in/yaml.v2"
)

const (
	defaultConfigDir = ".stentor.d"
	defaultHosting   = "github"
	defaultMarkup    = "markdown"
)

var (
	errBadHosting        = errors.New("hosting must be one of 'github' or 'gitlab'")
	errBadMarkup         = errors.New("markup must be one of 'markdown' or 'rst'")
	errBadRepository     = errors.New("repository must be in the format <user name>/<repository name>")
	errMissingRepository = errors.New("repository is required")
)

// Config represents the project's configuration for stentor.
type Config struct {
	// Repository is the name of your repository in <username>/<repo name> format.
	Repository string
	// FragmentDir is the path to the directory holding the project's news fragments.
	// Defaults to '.stentor.d'.
	FragmentDir string `yaml:"fragments"`
	// Hosting is the source repository host.
	// When Markup is set to markdown, this also determines the markdown flavor.
	// Currently, github and gitlab are supported.
	// Defaults to github.
	Hosting string
	// Markup sets the format of your changelog.
	// Currently, markdown and rst (ReStructuredText) are supported.
	// Defaults to markdown
	Markup string
}

// Parse the string s into a Config.
func Parse(s string) (*Config, error) {
	c, err := parseConfig(s)
	if err != nil {
		return nil, err
	}
	if err := validateConfig(c); err != nil {
		return nil, err
	}
	return c, nil
}

func parseConfig(s string) (*Config, error) {
	c := &Config{}
	err := yaml.Unmarshal([]byte(s), c)
	if err != nil {
		return nil, err
	}
	if c.FragmentDir == "" {
		c.FragmentDir = defaultConfigDir
	}
	if c.Hosting == "" {
		c.Hosting = defaultHosting
	}
	if c.Markup == "" {
		c.Markup = defaultMarkup
	}
	return c, nil
}

func validateConfig(c *Config) error {
	if c.Repository == "" {
		return errMissingRepository
	}
	if strings.Count(c.Repository, "/") != 1 || c.Repository[0] == '/' || c.Repository[len(c.Repository)-1] == '/' {
		return errBadRepository
	}
	if c.Hosting != defaultHosting && c.Hosting != "gitlab" {
		return errBadHosting
	}
	if c.Markup != defaultMarkup && c.Markup != "rst" {
		return errBadMarkup
	}
	return nil
}

// MustParse behaves the same as Parse, but panics if there is an error parsing the config file.
func MustParse(path string) *Config {
	c, err := Parse(path)
	if err != nil {
		panic(err)
	}
	return c
}
