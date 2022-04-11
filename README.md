<!-- [https://github.com/egonelbre/gophers/blob/master/vector/friends/heart-balloon.svg]() -->

# prt: The missing gh pr functionality

`gh pr` is a fantastic set of utilities, but it only accepts the PR number, URL, or branch name. This extension extends this capability by allowing you to enter part of the PR title instead.

## Installation

Install the gh CLI - see the [installation](https://github.com/cli/cli#installation)

*Installation requires a minimum version (2.0.0) of the the GitHub CLI that supports extensions.*

```
gh extension install jdahm/gh-prt
```

## Getting started

The basic functionality is

```
gh prt core_command "some string" <options...>
```

For example, the command below will attempt to checkout the latest PR with the title matching `Fix bug`

```
gh prt checkout "Fix bug"
```

The search is **case insensitive**.

## Details

All of the core commands are supported, with the exception of `create`.
`close` and `merge` are considered somewhat dangerous, and therefore require the `-s` or `--sudo` flags passed to `prt`.

There are a few `prt`-related options that can go after the `core_command` token in the command:

* `--sudo`: Take off the shackles!
* `--dry-run`: Print the sub-command that will be run to stdout, but do not execute it.
