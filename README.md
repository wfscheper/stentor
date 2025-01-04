# stentor

noun

1. (in the Iliad) a Greek herald with a loud voice.
1. (lowercase) a person having a very loud or powerful voice.
1. (lowercase) a trumpet-shaped, ciliate protozoan of the genus Stentor.

`stentor` is a CLI
for generating a change log or release notes
from a set of fragment files and templates.
It was inspired by [towncrier](https://github.com/twisted/towncrier)
and [git-chlog](https://github.com/git-chglog/git-chglog).


## Badges

![Build](https://github.com/wfscheper/stentor/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/wfscheper/stentor/branch/master/graph/badge.svg)](https://codecov.io/gh/wfscheper/stentor)
[![License](https://img.shields.io/github/license/wfscheper/stentor)](/LICENSE)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)


## Installation

### Binary

Download a pre-built binary for your OS and Architecture from the [releases](./releases) page.


### Go

You can also build `stentor` directly using `go install`:

```bash
go install github.com/wfscheper/stentor/cmd/stentor@latest
```


## Usage

### Setup

1. Create a `.stentor.d` directory in your git repository.

   ```bash
   mkdir .stentor.d
   ```

   This is where your fragments, configuration, and templates will go.

1. Create a minimal stentor config file.

   ```bash
   $ cat >.stentor.d/stentor.toml << EOF
   [stentor]
   repository = "https://github.com/myname/myrepo"
   EOF
   ```

1. Create some fragment files.

   ```bash
   $ cat >.stentor.d/1.feature.md << EOF
   Added the foo feature.

   The foo feature is full of foos,
   and is awesome.
   EOF
   $ cat >.stentor.d/2.fix.md << EOF
   Fixed parsing foos that contain special characters

   `fooer` no longer chokes when parsing a foo with the special characters `!@#$%`.
   EOF
   ```

1. *(Optional)* Create initial `CHANGELOG.md`.
   This is optional,
   but lets you write an intro section.

   ```bash
   $ cat >CHANGELOG.md << EOF
   # Changelog

   All notable changes to this project will be documented in this file.

   The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
   and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

   Changes for the next release can be found in the [".stentor.d" directory](./.stentor.d).

   <!-- stentor output starts -->
   EOF
   ```

1. Commit the changes.

### First release

This assumes that you are making a first release,
ie. there is no "previous" version.

1. Run `git log`,
   go to the start of your commit history,
   and record the first commit hash.

1. Run `stentor` to see the output it would add to the CHANGELOG.md file.

   ```bash
   $ stentor v0.1.0 2e808ef3f3a64e8c5965bcc130d4006d6abb56a1
   ## [v0.1.0] - 2006-01-02

   ### Features

   - Added the foo feature

     The foo feature is full of foos, and is awesome.
     [#1](https://github.com/myname/myrepo/issues/1)


   ### Bug fixes

   - Fixed parsing foos that contain special characters

     `fooer` no longer chokes when parsing a foo with the special characters `!@#$%`.
     [#2](https://github.com/myname/myrepo/issues/2)

   [v0.1.0]: https://github.com/myname/myrepo/compare/2e808ef3f3a64e8c5965bcc130d4006d6abb56a1...v0.1.0

   ---

   ```

### General release

1. . Run `stentor` to see the output it would add to the CHANGELOG.md file.

   ```bash
   $ stentor v0.2.0 v0.1.0
   ## [v0.2.0] - 2006-01-02

   ### Features

   - Added the foo feature

     The foo feature is full of foos, and is awesome.
     [#1](https://github.com/myname/myrepo/issues/1)


   ### Bug fixes

   - Fixed parsing foos that contain special characters

     `fooer` no longer chokes when parsing a foo with the special characters `!@#$%`.
     [#2](https://github.com/myname/myrepo/issues/2)

   [v0.2.0]: https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0

   ---

   ```

1. Use the `-release` flag to consume the fragments
   and update the news file.

   **Note:** If a CHANGELOG.md does not exist already, one will be created.

   ```bash
   $ stentor -release v0.2.0 v0.1.0
   $ git status
   On branch master
   Changes not staged for commit:
     (use "git add <file>..." to update what will be committed)
     (use "git restore <file>..." to discard changes in working directory)
          modified:   CHANGELOG.md
          deleted:    .stentor.d/1.feature.md
          deleted:    .stentor.d/2.fix.md

   no changes added to commit (use "git add" and/or "git commit -a")
   ```

### Subsequent releases

1. Run stentor to see the output it would add to CHANGELOG.md.

   ```bash
   $ stentor v0.2.0 v0.1.0
   ## [v0.2.0] - 2006-01-02

   ### Features

   - Added the foo feature

     The foo feature is full of foos, and is awesome.
     [#1](https://github.com/myname/myrepo/issues/1)


   ### Bug fixes

   - Fixed parsing foos that contain special characters

     `fooer` no longer chokes when parsing a foo with the special characters `!@#$%`.
     [#2](https://github.com/myname/myrepo/issues/2)

   [v0.2.0]: https://github.com/myname/myrepo/compare/v0.1.0...v0.2.0


   ---

   $ git add .stentor.d/
   $ git commit -m "Setup stentor to generate CHANGELOG.md"
   ```

1. Use the `-release` flag to consume the fragments
   and update the news file.

   **Note:** If a CHANGELOG.md does not exist already, one will be created.

   ```bash
   $ stentor -release v0.2.0 v0.1.0
   $ git status
   On branch master
   Changes not staged for commit:
     (use "git add <file>..." to update what will be committed)
     (use "git restore <file>..." to discard changes in working directory)
          modified:   CHANGELOG.md
          deleted:    .stentor.d/1.feature.md
          deleted:    .stentor.d/2.fix.md

   no changes added to commit (use "git add" and/or "git commit -a")
   ```
