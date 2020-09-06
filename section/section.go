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
