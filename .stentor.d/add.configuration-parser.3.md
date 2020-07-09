Added a toml configuration parser for `stentor`.
Supported options are:

- *repository*

  The username/repo of the github or gitlab repository.

- *fragments*

  The directory news fragments will be found in.
  Defaults to `.stentor.d`.

- *hosting*

  The SCM host,
  either *github* or *gitlab*

- *markup*

  The style of markup to use.
  Either *markdown* or *rst*
