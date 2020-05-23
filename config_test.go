package stentor

import (
	"reflect"
	"strings"
	"testing"

	"pgregory.net/rapid"
)

func Test_parseConfig(t *testing.T) {
	defaultConfig := &Config{
		FragmentDir: ".stentor.d",
		Hosting:     "github",
		Markup:      "markdown",
	}

	c, err := parseConfig("")
	if err != nil {
		t.Errorf("parseConfig(%q) returned an error: %v", "", err)
	}
	if got, want := c, defaultConfig; !reflect.DeepEqual(got, want) {
		t.Errorf("parseConfig(%q) returned %+v, want %+v", "", got, want)
	}

	// bad yaml
	y := "hosting: foo\n\tmarkup: markdown\n"
	if _, err := parseConfig(y); err == nil {
		t.Errorf("parseConfig(%q) returned nil", y)
	}
}

func Test_validateConfig(t *testing.T) {
	t.Parallel()

	genHosting := rapid.SampledFrom([]string{"github", "gitlab"})
	genMarkup := rapid.SampledFrom([]string{"markdown", "rst"})
	genRepository := rapid.Just("foo/bar")

	t.Run("invalid hosting", rapid.MakeCheck(func(t *rapid.T) {
		c := &Config{
			Repository: genRepository.Draw(t, "repository").(string),
			Hosting:    rapid.String().Draw(t, "hosting").(string),
			Markup:     genMarkup.Draw(t, "markup").(string),
		}
		if got, want := validateConfig(c), errBadHosting; got != want {
			t.Errorf("validateConfig(%+v) returned %v, want %v", c, got, want)
		}
	}))

	t.Run("invalid markup", rapid.MakeCheck(func(t *rapid.T) {
		c := &Config{
			Hosting:    genHosting.Draw(t, "hosting").(string),
			Markup:     rapid.String().Draw(t, "markup").(string),
			Repository: genRepository.Draw(t, "repository").(string),
		}
		if got, want := validateConfig(c), errBadMarkup; got != want {
			t.Errorf("validateConfig(%+v) returned %v, want %v", c, got, want)
		}
	}))

	t.Run("invalid repository", rapid.MakeCheck(func(t *rapid.T) {
		c := &Config{
			Repository: rapid.String().Filter(func(s string) bool {
				return strings.Count(s, "/") != 1
			}).Draw(t, "repository").(string),
			Hosting: genHosting.Draw(t, "hosting").(string),
			Markup:  genMarkup.Draw(t, "markup").(string),
		}
		want := errBadRepository
		if c.Repository == "" {
			want = errMissingRepository
		}
		if got := validateConfig(c); got != want {
			t.Errorf("validateConfig(%+v) returned %v, want %v", c, got, want)
		}
	}))
}
