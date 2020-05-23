package templates

import (
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"github-markdown-header"},
		{"github-markdown-section"},
		{"github-rst-header"},
		{"github-rst-section"},
	}

	t.Parallel()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.name)
			if err != nil {
				t.Errorf("New returned an error: %v", err)
			}
		})
	}
}
