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

// Package stentor provides some shared constants.
package stentor

// Supported hosting platforms.
const (
	//
	HostingGithub = "github"
	HostingGitlab = "gitlab"
)

// Supported markup formats.
const (
	MarkupMD  = "markdown"
	MarkupRST = "rst"
)

// Comment styles that separate the news file's header from the releases.
const (
	CommentMD  = "<!-- stentor output starts -->"
	CommentRST = ".. stentor output starts\n"
)
