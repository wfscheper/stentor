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
  Conventional commits takes precendence when it *How to Write a Git Commit Message* disagree.


## Tools

`stentor` uses [mage](https://github.com/magefile/mage)
to manage local build tasks.

The `mage lint` command will run our [pre-commit] hooks.
Please run these before creating your pull request to save yourself wasted time.

## Tests

`stentor` uses the standard `testing` package plus the [rapid] property testing framework.

When writing tests, some guidelines that we suggest you follow are:

- Use the variables `want` and `got` to hold data for assertions:

  ```golang
  x := f()
  if got, want := x, 42; got != want {
      t.Errorf("f() returned %v, want %v", got, want)
  }
  ```

- `stentor` relies on the [rapid] proprety testing framework for most tests.
  These style of tests are generally preferred to [table-driven tests],
  though the latter have their place.
- Use the `mage test` command to run the test suite locally.
  This ensures that your tests are run in the same way as our CI.

## Documentation

- `stentor` uses [semantic newlines] in all documentation files and comments:

  ```text
  This is a sentence.
  This is another sentance,
  with a clause.
  ```

### Changelog

All changes should include a changelog entry.
Add a single file to the `.stentor.d` directory as part of your pull request
named `<pull request #>.(breaking|build|chore|deprecation|feature|fix|test).md`.

Changelog entries should follow these rules:

- Use [semantic newlines], just like other documentation.
- Wrap the names of things in backticks, `like this`.
- Wrap arguments with asterisks: _these_ or _attributes_.
- Names of functions or other callables should be followed by parentheses,
  `my_cool_function()`.
- Use the active voice and either present tense or simple past tense.

  - Added `my_cool_function()` to do cool things.
  - Creating `Foo` objects with the _many_ argument no longer raises a `RuntimeError`.

- For changes that address multiple pull requests,
  create multiple fragments with the same contents.

To see what `stentor` will add to the `CHANGELOG.md`, run `mage changelog`.

## Development

First, make sure you have the latest version of [go] installed.
While `stentor` supports the two most recent releases,
development should be done with the most recent version.

Next, make a fork of the `stentor` repository by going to <https://github.com/wfscheper/stentor>
and clicking on the **Fork** button near the top of the page.

Then clone your fork of the `stentor` repository:

```bash
git clone git@github.com:<username>/pymaven.git
```

Installing the [pre-commit] hooks is recommend to ensure your commit will pass our CI checks:

```bash
pre-commit install
pre-commit install -t commit-msg
pre-commit run --all-files
```

[ci]: https://github.com/wfscheper/stentor/actions?query=workflow%3ACI
[pre-commit]: https://pre-commit.com/
[semantic newlines]: https://rhodesmill.org/brandon/2012/one-sentence-per-line/
[rapid]: https://github.com/flyingmutant/rapid
[table-driven tests]: https://github.com/golang/go/wiki/TableDrivenTests
[go]: https://golang.org/dl/
