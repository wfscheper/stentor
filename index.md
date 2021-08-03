# stentor

*noun*
1. (in the Iliad) a Greek herald with a loud voice.
1. (lowercase) a person having a very loud or powerful voice.
1. (lowercase) a trumpet-shaped, ciliate protozoan of the genus Stentor.

Current release: [v0.2.3](https://github.com/wfscheper/stentor/releases/tag/v0.2.3)

## The Problem

Managing changelog files is annoying.
Many project maintainers discover they can either have
a richly descriptive changelog
that is a constant source of merge conflicts
and last minute nightmares of cherry-picking releases,
or they can have an easily generated changelog
that is a word salad of filtered commit messages.

`stentor` builds your changelog from individual fragment files,
committed alongside the code they reference.
This means you can write your changelog in comprehensible prose,
while also avoiding the maintenance pain
of multiple commits modifying the same file.
