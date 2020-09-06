package fragment

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tmpdir := t.TempDir()
	fn := filepath.Join(tmpdir, "ticket.section.md")
	err := ioutil.WriteFile(fn, []byte(`contents`), 0600)
	require.NoError(t, err)

	f, s, err := New(fn)
	require.NoError(t, err)

	assert.Equal(t, "ticket", f.Issue)
	assert.Equal(t, "contents", f.Text)
	assert.Equal(t, "section", s)
}

func TestNew_error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"ticket.section.more.md", "'ticket.section.more.md' is not a valid fragment file: too many parts."},
		{"ticket.section", "'ticket.section' is not a valid fragment file: not enough parts"},
		{"ticket..md", "'ticket..md' is not a valid fragment file: empty section"},
		{".section.md", "'.section.md' is not a valid fragment file: empty issue"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpdir := t.TempDir()
			fn := filepath.Join(tmpdir, tt.name)
			err := ioutil.WriteFile(fn, []byte(`contents`), 0600)
			require.NoError(t, err)

			_, _, err = New(fn)
			assert.EqualError(t, err, tt.want)
		})
	}
}
