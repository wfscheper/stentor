package templates

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "github-markdown-section",
		FileModTime: time.Unix(1604548150, 0),

		Content: string("{{- /*\n    Copyright © 2020 The Stentor Authors\n    Licensed under the Apache License, Version 2.0 (the \"License\");\n    you may not use this file except in compliance with the License.\n    You may obtain a copy of the License at\n\n        http://www.apache.org/licenses/LICENSE-2.0\n\n    Unless required by applicable law or agreed to in writing, software\n    distributed under the License is distributed on an \"AS IS\" BASIS,\n    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n    See the License for the specific language governing permissions and\n    limitations under the License.\n*/ -}}\n{{- $repository := .Repository -}}\n{{- $sectionHeader := .SectionHeader -}}\n\n{{ .Header }} [{{ .Version }}] - {{ .Date.Format \"2006-01-02\" }}\n{{- range .Sections -}}\n{{- if or .Fragments .ShowAlways }}\n\n{{ $sectionHeader }} {{ .Title }}\n\n{{ range .Fragments -}}\n- {{ .Text | indent 2 }}{{ if .Issue }}\n  [#{{ .Issue }}]({{ $repository }}/issues/{{ .Issue }}){{ end }}\n{{ else -}}\n{{ if .ShowAlways -}}\nNo significant changes.\n{{ end -}}\n{{ end -}}\n{{ end -}}\n{{ else }}\n\nNo significant changes.\n{{- end }}\n\n[{{ .Version }}]: {{ .Repository }}/compare/{{ .PreviousVersion }}...{{ .Version }}\n\n\n----\n\n"),
	}
	file3 := &embedded.EmbeddedFile{
		Filename:    "github-rst-section",
		FileModTime: time.Unix(1604548275, 0),

		Content: string("{{- /*\n    Copyright © 2020 The Stentor Authors\n    Licensed under the Apache License, Version 2.0 (the \"License\");\n    you may not use this file except in compliance with the License.\n    You may obtain a copy of the License at\n\n        http://www.apache.org/licenses/LICENSE-2.0\n\n    Unless required by applicable law or agreed to in writing, software\n    distributed under the License is distributed on an \"AS IS\" BASIS,\n    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n    See the License for the specific language governing permissions and\n    limitations under the License.\n*/ -}}\n{{- $repository := .Repository -}}\n{{- $sectionHeader := .SectionHeader -}}\n{{- $date := .Date.Format \"2006-01-02\" -}}\n\n`{{ .Version }}`_ - {{ $date }}\n{{ .Header | repeat (sum (len .Version) (len $date) 6) }}\n{{- range .Sections -}}\n{{- if or .Fragments .ShowAlways }}\n\n{{ .Title }}\n{{ $sectionHeader | repeat (len .Title) }}\n\n{{ range .Fragments -}}\n- {{ .Text | indent 2 }}{{ if .Issue }}\n  `#{{ .Issue }} <{{ $repository }}/issues/{{ .Issue }}>`_{{ end }}\n{{ else -}}\n{{ if .ShowAlways -}}\nNo significant changes.\n{{ end -}}\n{{ end -}}\n{{ end -}}\n{{ else }}\n\nNo significant changes.\n{{- end }}\n\n.. _{{ .Version }}: {{ .Repository }}/compare/{{ .PreviousVersion }}...{{ .Version }}\n\n\n----\n\n"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "gitlab-markdown-section",
		FileModTime: time.Unix(1604548280, 0),

		Content: string("{{- /*\n    Copyright © 2020 The Stentor Authors\n    Licensed under the Apache License, Version 2.0 (the \"License\");\n    you may not use this file except in compliance with the License.\n    You may obtain a copy of the License at\n\n        http://www.apache.org/licenses/LICENSE-2.0\n\n    Unless required by applicable law or agreed to in writing, software\n    distributed under the License is distributed on an \"AS IS\" BASIS,\n    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n    See the License for the specific language governing permissions and\n    limitations under the License.\n*/ -}}\n{{- $repository := .Repository -}}\n{{- $sectionHeader := .SectionHeader -}}\n\n{{ .Header }} [{{ .Version }}] - {{ .Date.Format \"2006-01-02\" }}\n{{- range .Sections -}}\n{{- if or .Fragments .ShowAlways }}\n\n{{ $sectionHeader }} {{ .Title }}\n\n{{ range .Fragments -}}\n- {{ .Text | indent 2 }}{{ if .Issue }}\n  [#{{ .Issue }}]({{ $repository }}/-/issues/{{ .Issue }}){{ end }}\n{{ else -}}\n{{ if .ShowAlways -}}\nNo significant changes.\n{{ end -}}\n{{ end -}}\n{{ end -}}\n{{ else }}\n\nNo significant changes.\n{{- end }}\n\n[{{ .Version }}]: {{ .Repository }}/-/compare/{{ .PreviousVersion }}...{{ .Version }}\n\n\n----\n\n"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "gitlab-rst-section",
		FileModTime: time.Unix(1604548289, 0),

		Content: string("{{- /*\n    Copyright © 2020 The Stentor Authors\n    Licensed under the Apache License, Version 2.0 (the \"License\");\n    you may not use this file except in compliance with the License.\n    You may obtain a copy of the License at\n\n        http://www.apache.org/licenses/LICENSE-2.0\n\n    Unless required by applicable law or agreed to in writing, software\n    distributed under the License is distributed on an \"AS IS\" BASIS,\n    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n    See the License for the specific language governing permissions and\n    limitations under the License.\n*/ -}}\n{{- $repository := .Repository -}}\n{{- $sectionHeader := .SectionHeader -}}\n{{- $date := .Date.Format \"2006-01-02\" -}}\n\n`{{ .Version }}`_ - {{ $date }}\n{{ .Header | repeat (sum (len .Version) (len $date) 6) }}\n{{- range .Sections -}}\n{{- if or .Fragments .ShowAlways }}\n\n{{ .Title }}\n{{ $sectionHeader | repeat (len .Title) }}\n\n{{ range .Fragments -}}\n- {{ .Text | indent 2 }}{{ if .Issue }}\n  `#{{ .Issue }} <{{ $repository }}/-/issues/{{ .Issue }}>`_{{ end }}\n{{ else -}}\n{{ if .ShowAlways -}}\nNo significant changes.\n{{ end -}}\n{{ end -}}\n{{ end -}}\n{{ else }}\n\nNo significant changes.\n{{- end }}\n\n.. _{{ .Version }}: {{ .Repository }}/-/compare/{{ .PreviousVersion }}...{{ .Version }}\n\n\n----\n\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1604111653, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "github-markdown-section"
			file3, // "github-rst-section"
			file4, // "gitlab-markdown-section"
			file5, // "gitlab-rst-section"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`templates`, &embedded.EmbeddedBox{
		Name: `templates`,
		Time: time.Unix(1604111653, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"github-markdown-section": file2,
			"github-rst-section":      file3,
			"gitlab-markdown-section": file4,
			"gitlab-rst-section":      file5,
		},
	})
}
