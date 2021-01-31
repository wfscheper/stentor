# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Changes for the next release can be found in the [".stentor.d" directory](./.stentor.d).

<!-- stentor output starts -->

## [v0.2.2] - 2020-11-30

### Fixed

- Pass `indent` and `repeat` template functions to custom templates.
  [#26](https://github.com/wfscheper/stentor/issues/26)
- Fixed an issue where `stentor` could fail to find the start comment.

  The original way `stentor` scanned for the start comment
  worked so long as the newline after the comment
  didn't align just after the end of the internal read buffer.
  To fix that,
  stentor now scans the end of the read buffer
  looking for partial matches with the start comment.
  [#28](https://github.com/wfscheper/stentor/issues/28)


[v0.2.2]: https://github.com/wfscheper/stentor/compare/v0.2.1...v0.2.2


----


## [v0.2.1] - 2020-11-12

### Fixed

- `stentor` now produces a release,
  even if there are no fragment files.
  The built-in templates
  will produce a release that says "No significant changes."
  [#21](https://github.com/wfscheper/stentor/issues/21)


[v0.2.1]: https://github.com/wfscheper/stentor/compare/v0.2.0...v0.2.1


----


## [v0.2.0] - 2020-10-30

### Changed

- The repository for a project must be a http or https URL.
  This breaking change is required
  to support privately hosted repositories
  with the built-in templates.

  API changes:
  - `release.New` now returns a `*release.Release` and an `error`
  - Config structs moved
    from the `main` package in `cmd/stentor`
    to a new `config` package,
    so that they are importable.
  - `SectionConfig` is now `config.Section` to comply with go naming rules.

  Behavior changes:
  - `config.ValidateConfig` returns an error
    if `repository` is not parseable by `url.Parse`
    or is not a http or https URL.
  [#14](https://github.com/wfscheper/stentor/issues/14)
- The call signature of `newsfile.WriteFragments`
  was changed to take a bool `keepHeader`.
  This breaking change is required
  to fix the duplication of the newsfile header.

  API Changes:
  - `newsfile.WriteFragments` now takes a new boolean argument,
    indicating whether to keep the existing newsfile header or not.

  Behavior changes:
  - `stentor` no longer provides a default header template.
    Instead,
    the existing newsfile header will be preserved,
    unless configured with a `header_template`.
  [#18](https://github.com/wfscheper/stentor/issues/18)


### Added

- Added a `SetSections` method to `release.Release`.

  `SetSectiosn` takes a `[]config.Section` and `[]fragment.Fragment`,
  and populates the `Release`'s `Section` member
  and their `Fragments`.
  [#14](https://github.com/wfscheper/stentor/issues/14)
- Added built-in markdown and rst templates
  for gitlab repositories.
  [#15](https://github.com/wfscheper/stentor/issues/15)


[v0.2.0]: https://github.com/wfscheper/stentor/compare/v0.1.0...v0.2.0


----


## [v0.1.0] - 2020-09-20

Initial release

[v0.1.0]: https://github.com/wfscheper/stentor/compare/2e808ef...v0.1.0
