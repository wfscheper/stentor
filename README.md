# stentor

noun

1. (in the Iliad) a Greek herald with a loud voice.
1. (lowercase) a person having a very loud or powerful voice.
1. (lowercase) a trumpet-shaped, ciliate protozoan of the genus Stentor.

Stentor is a CLI for generating a change log or release notes from a set of fragment files and templates.
It was inspired by [towncrier](https://github.com/twisted/towncrier) and [git-chlog](https://github.com/git-chglog/git-chglog).

## Badges

![Build](https://github.com/wfscheper/stentor/workflows/Build/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/wfscheper/stentor/badge.svg?branch=master)](https://coveralls.io/github/wfscheper/stentor?branch=master)
[![License](https://img.shields.io/github/license/wfscheper/stentor)](/LICENSE)

## Installation

### Binary

Download a pre-built binary for your OS and Architecture from the [releases](./releases) page.

### Go

Run

```golang
GO111MODULES=on go get github.com/wfscheper/stentor@v1.0.0
```

## Usage

### Setup

This example assumes that there is already a v0.1.0 tag.

1. Create a `.stentor.d` directory in your git repositry.

   ```bash
   mkdir .stentor.d
   ```

   This is where your fragments, configuration, and templates will go.

1. Create a minimal stentor config file.

   ```bash
   $ cat >.stentor.d/stentor.yaml << EOF
   repository: myname/myrepo
   EOF
   ```

1. Create some fragment files.

   ```bash
   $ cat >.stentor.d/1.feature.md << EOF
   Added the foo feature.

   The foo feature is full of foos, and is awesome.
   EOF
   $ cat >.stentor.d/2.fix.md << EOF
   Fixed parsing foos that contain special characters

   `fooer` no longer chokes when parsing a foo with the special characters `!@#$%`.
   EOF
   ```

1. Run stentor to see the output it would add to CHANGELOG.md.

   ```bash
   $ stentor
   # Changelog

   All notable changes to this project will be documented in this file.

   The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
   and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

   <!-- stentor output starts -->


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
   $ git add .stentory.d/
   $ git commit -m "Setup stentor to generate CHANGELOG.md"
   ```

`stentor` will attempt to determine the next version from the types of fragements that it finds.
You can override this by giving `stentor` an explicit version:

```bash
stentor -version v1.0.0
```

1. Use the `-release` flag to consume the news fragments.

   ```bash
   $ stentor -release
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
