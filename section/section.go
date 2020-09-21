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

package section

import "github.com/wfscheper/stentor/fragment"

// Section represents a collection of changes. Features, bug fixes, etc.
type Section struct {
	// Fragments is the list of changes of this section type in the release.
	Fragments []fragment.Fragment
	// ShowAlways is a boolean indicating if this section should be included in the
	// news file even if there are no fragments.
	ShowAlways bool
	// Title is the string written to the news file for this section.
	Title string
}
