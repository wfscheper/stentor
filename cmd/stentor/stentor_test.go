package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wfscheper/stentor/config"
	"github.com/wfscheper/stentor/fragment"
)

func TestStentor_displayVersion(t *testing.T) {
	tmpdir := t.TempDir()

	out := &bytes.Buffer{}
	s := New(tmpdir, []string{}, []string{}, io.Discard, out)
	s.displayVersion()

	assert.Equal(t, "stentor dev built from unknown on unknown\n", out.String())
}

func TestStentor_verifyFragmentSections(t *testing.T) {
	tests := []struct {
		name      string
		sections  []config.Section
		fragments []fragment.Fragment
		want      string
	}{
		{
			name: "valid",
			sections: []config.Section{
				{
					Name:      "Added",
					ShortName: "add",
				},
				{
					Name:      "Removed",
					ShortName: "remove",
				},
			},
			fragments: []fragment.Fragment{
				{
					Section: "add",
				},
				{
					Section: "remove",
				},
			},
		},
		{
			name: "invalid",
			want: "fragment files contained the following invalid section names: " +
				"[invalid]. section names must be one of the following: [add remove]",
			sections: []config.Section{
				{
					Name:      "Added",
					ShortName: "add",
				},
				{
					Name:      "Removed",
					ShortName: "remove",
				},
			},
			fragments: []fragment.Fragment{
				{
					Section: "add",
				},
				{
					Section: "remove",
				},
				{
					Section: "invalid",
				},
			},
		},
		{
			name: "no valid sections",
			want: "fragment files contained the following invalid section names: " +
				"[add remove]. section names must be one of the following: []",
			sections: []config.Section{},
			fragments: []fragment.Fragment{
				{
					Section: "add",
				},
				{
					Section: "remove",
				},
			},
		},
		{
			name: "no fragments",
			want: "",
			sections: []config.Section{
				{
					Name:      "Added",
					ShortName: "add",
				},
				{
					Name:      "Removed",
					ShortName: "remove",
				},
			},
			fragments: []fragment.Fragment{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyFragmentSections(tt.sections, tt.fragments)
			if tt.want == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tt.want)
			}
		})
	}
}
