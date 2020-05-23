package templates

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "github-markdown-header",
		FileModTime: time.Unix(1590199928, 0),

		Content: string("# Changelog\n\nAll notable changes to this project will be documented in this file.\n\nThe format is based on [Keep a Changelog], and this project adheres to [Semantic Versioning].\n\n[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/\n[Semantic Versioning]: https://semver.org/spec/v2.0.0.html\n\n<!-- stentor output starts -->\n"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "github-markdown-section",
		FileModTime: time.Unix(1590199928, 0),

		Content: string("{{- $r := . -}}\n{{ $r.Header }} [{{ $r.Version }}](https://github.com/{{ $r.Repository }}/compare/{{ $r.PreviousVersion }}...{{ $r.Version }}) - {{ $r.Date.Format \"2006-01-02\" }}\n\n{{ range $r.Sections -}}\n{{ $r.SectionHeader }} {{ .Title }}\n\n{{ if .Fragments -}}\n{{ range .Fragments -}}\n- {{ .Text | indent 2  }}\n  [#{{ .Issue }}](https://github.com/{{ $r.Repository }}/issues/{{ .Issue }})\n{{ end }}\n{{ else if .ShowAlways -}}\nNo significant changes.\n{{ end }}\n{{ end }}\n----\n"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "github-rst-header",
		FileModTime: time.Unix(1590199928, 0),

		Content: string("Changelog\n=========\n\nAll notable changes to this project will be documented in this file.\n\nThe format is based on `Keep a Changelog <https://keepachangelog.com/en/1.0.0/>`_,\nand this project adheres to `Semantic Versioning <https://semver.org/spec/v2.0.0.html>`_.\n\n.. stentor output starts\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "github-rst-section",
		FileModTime: time.Unix(1590199928, 0),

		Content: string("{{- $r := . -}}\n{{- $date := $r.Date.Format \"2006-01-02\" -}}\n`{{ $r.Version }} <https://github.com/{{ $r.Repository }}/compare/{{ $r.PreviousVersion }}...{{ $r.Version }}>`_ - {{ $date }}\n{{ $r.Header | repeat (sum (len $r.Version) (len $date) 6) }}\n\n{{ range $r.Sections -}}\n{{ .Title }}\n{{ $r.SectionHeader | repeat (len .Title) }}\n\n{{ if .Fragments -}}\n{{ range .Fragments -}}\n- {{ .Text | indent 2 }}\n  `#{{ .Issue }} <https://github.com/{{ $r.Repository }}/issues/{{ .Issue }}>`_\n{{ end }}\n{{ else if .ShowAlways -}}\nNo significant changes.\n{{ end }}\n{{ end }}\n----\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1590199928, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "github-markdown-header"
			file3, // "github-markdown-section"
			file4, // "github-rst-header"
			file5, // "github-rst-section"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1590199928, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"github-markdown-header":  file2,
			"github-markdown-section": file3,
			"github-rst-header":       file4,
			"github-rst-section":      file5,
		},
	})
}
