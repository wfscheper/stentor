# stentor

*noun*

1. (in the Iliad) a Greek herald with a loud voice.
1. (lowercase) a person having a very loud or powerful voice.
1. (lowercase) a trumpet-shaped ciliate protozoan of the genus Stentor.

Current release: [v0.3.0](https://github.com/wfscheper/stentor/releases/tag/v0.3.0)

## The Problem

Most software projects
maintain some record of the changes made in each release.
For many projects,
this takes the form of a changelog
or release notes.
`stentor` refers to this as a "news file".

Managing news files is annoying.
Many project maintainers discover they can either have
a richly descriptive news file
that is a constant source of merge conflicts
and last minute nightmares of cherry-picking releases,
or they can have an easily generated changelog
that is a word salad of filtered commit messages.

`stentor` solves this problem
by updating your news file with content
sourced from individual files,
called "fragment files".
The fragment file describing a change
is included as part of the commit making the change.
This means you can write your news file in comprehensible prose,
while also avoiding the maintenance pain
of multiple commits modifying the same file.

The fragment files describing unreleased changes
are kept together in a directory known as the
"fragment directory".
Typically,
`stentor`'s configuration file is also in this directory,
which ensures that the directory is retained by source control
even when there are no fragment files.

## Initial Setup

First,
you will need to create a configuration file for `stentor`.
The minimum configuration required
is to set the URL of your repository:

```bash
mkdir ./.stentor.d
printf '[stentor]\nrepository = "https://myhost.example/myrepo"\n' >./.stentor.d/stentor.toml
```

If your project does not have an existing news file,
then you need to create one for `stentor` to update.
A good starting place is [Keep a Changelog].
The `stentor` projet uses a modified version,
which points readers to the fragment directory for unreleaed changes:

```markdown
# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog],
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Changes for the next release can be found in the [".stentor.d" directory](./.stentor.d).

<!-- stentor output starts -->
```

The comment at the end is important:
`stentor` uses this comment
to mark where in the news file it should add content
when updating the news file for a release.

If your project already has a news file,
then simply add the `<!-- stentor output starts -->` comment
where you want stentor to add new content
when you create releases of your project.

At this point,
you are ready to start adding fragment files.

## Configuration

There are a few settings you may need to adjust,
depending on your project.
The first is what `stentor` refers to as "hosting".
This setting controls whether links are generated in
"GitHub" or "Gitlab" style.
GitHub is the default,
so you only need to set this
if your project is hosted on Gitlab:

```toml
[stentor]
hosting = "gitlab"
```

The second option you may need to set
is the flavor of markup used by your project's news file.
`stentor` currently supports [Markdown](https://www.markdownguide.org/)
and [reStructuredText](https://www.sphinx-doc.org/en/master/usage/restructuredtext/basics.html) (rST).
Markdown is the default,
so you only need to set this if your project uses rST:

```toml
[stentor]
markup = "rst"
```

## Creating Updates

Most projects name their news file `CHANGELONG.md` or `CHANGELOG.rst`.
If your project uses a different file name,
you can set this using the "news_file" config option:

```toml
[stentor]
news_file = "NEWS"
```

By default,
fragment files are stored in the `.stentor.d` directory,
If you want to use a different directory,
you can set this using the "fragment_dir" config option:

```toml
[stentor]
fragment_dir = "NEWS.d"
```

`stentor` expects fragment file names
to include important details about the change they describe:

* a reference to the ticket, Issue, or Pull Request (Merge Request in Gitalb)
* the type of change this fragment describes
* optionally, a quick description of the change

`<issue>.<type>[.<anything-you-want>].<markup extension>`.

Eg:

* `23.feat.md`
* `42.docs.fix-typo-in-README.rst`
* `JIRA-123.test.fix-flakey-test.md`

The contents of the fragment file
should use the same markup language as the news file,
and should describe the change.

## Customizing Your News

While `stentor` comes with its own style of news file,
your project may either have an existing style
or a different idea of how a news file should be laid out.
In these case,
`stentor` provides several ways to customize the output:

1. [Section configuration](#section-configuration)
1. [Header template](#header-template)
1. [Section template](#section-template)

### Section configuration

Sections are how `stentor` groups fragments together
to present a cohesive release.
`stentor` comes with a default set of sections
that map to the [@commitlint/config-conventional](https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional)
types.
Projects can change both the keyword used in the name of your fragment files,
and how the section title is rendered by the section template
via the `stentor.sections` option in the `stentor` config file:

```toml
[stentor]
sections = [
    { name = "New Features", short_name = "ft", show_always = true },
    { name = "Fixed Bugs", short_name = "bg", show_always = true },
]
```

### Header template

The header template is used to generate the content of the news file
that comes before the `<!-- stentor output starts -->` line.
If there is no configured header template,
then `stentor` will preserve the existing news file header.
To use a custom header template,
create the template in your project's fragment directory,
and set the `header_template` config option
to the name of the template file:

```toml
[stentor]
header_template = "my-header.tmpl"
```

### Section template

The section template is used
to render the contents of the fragment files.
The default section template
is based on the [Keep a Changelog] format.
To use a custom section template,
create the template in your project's fragment directory,
and set the `section_template` config option
to the name of the template file:

```toml
[stentor]
section_template = "my-section.tmpl"
```

### Template writing

`stentor` uses [go templates](https://pkg.go.dev/text/template)
to format the contents of the fragment files
and the news file header.
The template is passed the following [data structure](https://pkg.go.dev/github.com/wfscheper/stentor@v0.3.0/release#Release):

<!-- markdownlint-disable MD010 -->
```golang
type Release struct {
	// Date is the date of the release.
	Date time.Time
	// Header is the markup character used when writing the release header.
	Header string
	// PreviousVersion is the version before this release.
	PreviousVersion string
	// Repository is the URL of the project repository.
	Repository string
	// SectionHeader is the markup character used when writing a section header.
	SectionHeader string
	// Sections is the list of change types in this release.
	Sections []section.Section
	// Version is the version of this release.
	Version string
}
```
<!-- markdownlint-enable -->

In addition to the standard functions of go templates,
`stentor` template authors can leverage the following functions:

* `indent` adds *n* spaces immediately after every newline in *s*.
  For example,

  ```text
  - {{ input 2 "Line 1\nLine2" }}
  ```
  
  would render as
  
  ```text
  - Line1
    Line2
  ```

  The default `setntor` section template uses this
  to render the fragments as an unordered list.

* `repeat` is the same as [`strings.Repeat`](https://pkg.go.dev/strings#Repeat).
  For example,

  ```text
  Header
  {{ "=" | repeat 6 }}
  ```

  would render as

  ```text
  Header
  ======
  ```

  The default `stentor` section template uses this to render rST heading.

* `sum` returns the sum of its arguments.
  For example,

  ```text
  {{ sum (len "foo") (len "bar") }}
  ```

  would render as

  ```text
  6
  ```

  The default `stentor` section template uses this to render rST heading.

[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
