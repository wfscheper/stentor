package fragment

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want Fragment
	}{
		{"ticket.section.md", Fragment{"section", "ticket", "contents"}},
		{"ticket.section.extra-bit.md", Fragment{"section", "ticket", "contents"}},
		{"ticket.section.several.extra.bits.md", Fragment{"section", "ticket", "contents"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpdir, err := ioutil.TempDir("", "stentor-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpdir)

			fn := filepath.Join(tmpdir, tt.name)
			err = ioutil.WriteFile(fn, []byte(`contents`), 0600)
			require.NoError(t, err)

			if got, err := New(fn); assert.NoError(t, err) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestNew_error(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"foo", "not a valid fragment file: not enough parts"},
		{"foo.md", "not a valid fragment file: not enough parts"},
		{".section.md", "not a valid fragment file: empty issue"},
		{"ticket..md", "not a valid fragment file: empty section"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpdir, err := ioutil.TempDir("", "stentor-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpdir)

			fn := filepath.Join(tmpdir, tt.name)
			err = ioutil.WriteFile(fn, []byte(`contents`), 0600)
			require.NoError(t, err)

			_, err = New(fn)
			assert.EqualError(t, err, tt.want)
		})
	}
}
