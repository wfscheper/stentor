# How To Contribute

Thank you for considering contributing to `stentor`.
Some things we would like you to know:

- There's no such thing as a contribution that is "too small".
  Grammar and spelling corrections are just as important as fixing bugs or adding features.
- Each pull request should focus on one change.
  This makes them easier to review and merge.
- Tests are required for all changes.
  If you fix a bug,
  add a test to ensure that bug stays fixed.
  If you add a feature,
  add a test to show how that feature works.
- New features require new documentation.
- `stentor`'s API is still considered in flux,
  but API breaking changes need to clear a higher bar than new features.
- Pull requests that do not pass our [CI] checks will not receive feedback.
  If you need help passing the CI checks,
  add the `CI Triage` label to your pull request.
- `stentor` adheres to the [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/>)
  standard for commit messages.
  You also may want to read [How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/).
  The conventional commits standard
  takes precendence when *How to Write a Git Commit Message* disagrees.


## Tools

`stentor` uses [mage](https://github.com/magefile/mage)
to manage local build tasks.

In order to isolate build dependnecies like `mage`
from `stentor`'s runtime dependnecies,
the `magefile.go` is kept in a separate module.
A simple Makefile is provided to make running mage correctly easier.
Run `make help` to see a basic list of mage targets,
and what they do.


## Tests

Use the `make go:test` command to run the test suite locally.
This ensures that your tests are run in the same way as our CI.

`stentor` uses the [testify] package for most tests,
but some tests are written using the [rapid] property testing framework.
When writing tests, some guidelines that we suggest you follow are:

- `rapid` and `testify` do not play well with each other,
  so tests written using [rapid] should use the standard library `testing`
  package directly.
- When testing multiple scenarios,
  prefer [table-driven tests] over multiple test functions.
- Use the variables `want` and `got` to hold data for assertions:

  ```golang
  got := meaningofLife()
  want := 42
  assert.Equal(t, want, got)
  ```

- If the function under test returns an error,
  use an inline style to handle checking the error
  before doing any other assertions:

  ```golang
  if got, err := meaningOfLife(); assert.NoError(t, err) {
      want := 42
      assert.Equal(t, want, got)
  }
  ```

- Use the `require` package to handle test setup errors:

  ```golang
  require.NoError(t, ioutil.WriteFile("some_file.txt", ...))
  ```


## Documentation

- `stentor` uses [semantic newlines] in all documentation files and comments:

  ```text
  This is a sentence.
  This is another sentance,
  and it has a clause.
  ```


### Changelog

All changes should include a changelog entry.
Add a single file to the `.stentor.d` directory as part of your pull request
named `<issue #>.(breaking|build|chore|deprecation|feature|fix|test).md`.

Changelog entries should follow these rules:

- Use [semantic newlines],
  just like other documentation.
- Wrap the names of things in backticks,
  `like this`.
- Wrap arguments with asterisks:
  *these* or *attributes*.
- Names of functions or other callables should be followed by parentheses,
  `my_cool_function()`.
- Use the active voice
  and either present tense
  or simple past tense.

  - Added `my_cool_function()` to do cool things.
  - Creating `Foo` objects
    with the *many* argument
    no longer raises a `RuntimeError`.

- For the rare change that addresses multiple pull requests,
  create multiple fragments with the same contents.

To see what `stentor` will add to the `CHANGELOG.md`, run `make changelog`.


## Development

### VS Code Dev Container

The `stentor` project provides a devcontainer setup for VS Code,
and a set of recommended extensions.

To use the devcontainer with VS Code,
first install the [Remote - Container] extension.

You can either follow the [Local Development](#local-development) instructions
and mount your local clone into the devcontainer,
or clone the repository into a docker volume.

Read the official [documentation](https://code.visualstudio.com/docs/remote/containers)
for details.


### Local Development

First,
make sure you have the latest version of [go 1.17](https://golang.org/dl/) installed.
`stentor` supports the two most recent releases,
so development should be done with older stable release.

Next,
make a fork of the `stentor` repository
by going to <https://github.com/wfscheper/stentor>
and clicking on the **Fork** button near the top of the page.

Then clone your fork of the `stentor` repository:

```bash
git clone git@github.com:<username>/stentor.git
```

Installing the [pre-commit] hooks is recommend
to ensure your commit will pass our CI checks:

```bash
pre-commit install -t pre-commit -t commit-msg
pre-commit run --all-files
```

[Remote - Container]: https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers
[ci]: https://github.com/wfscheper/stentor/actions?query=workflow%3ACI
[pre-commit]: https://pre-commit.com/
[semantic newlines]: https://rhodesmill.org/brandon/2012/one-sentence-per-line/
[rapid]: https://github.com/flyingmutant/rapid
[table-driven tests]: https://github.com/golang/go/wiki/TableDrivenTests
[testify]: https://github.com/stretchr/testify
