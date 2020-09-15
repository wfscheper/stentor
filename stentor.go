package stentor

// Supported hosting platforms.
const (
	HostingGithub = "github"
	HostingGitlab = "gitlab"
)

// Supported markup formats.
const (
	MarkupMD  = "markdown"
	MarkupRST = "rst"
)

const (
	CommentMD  = "<!-- stentor output starts -->\n"
	CommentRST = ".. stentor output starts\n"
)
