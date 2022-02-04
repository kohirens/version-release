# How to Release

Explanation of how to incorporate the auto-release workflow into your
development life cycle, which the Version Release Orb provides.

## Prerequisites
1. Add a [git-chglog] config to your repo.
2. Add the [auto-release] workflow to your CI config.

## Get Started

In order to have a new version of your software published, use the following as
an example Gitflow:

1. Make a new branch based off your main trunk.
2. Develop your features, making commits according to the
   [Commit Message Guide](#commit-message-guide).
3. Push your branch to GitHub and make a PR back to the main trunk.
4. Once merged, and the CI pipeline reaches the `auto-release` workflow, there
   will be a new release published (assuming normal functionality of GitHub and
   CircleCI).

## Commit Message Guide

When the auto-release workflow is reached it scours the git logs for messages
skipping any without the tag prefixes mentioned below. It only goes back to the
last tag, or the beginning if there are no tags.

Prefix any of the keywords in the Descriptions below, an in any commit message,
to increment the version number accordingly.

| Increment | Description                                                                       |
|-----------|-----------------------------------------------------------------------------------|
| major     | Adding the words `BREAKING CHANGE` on a line by itself.                           |
| minor     | Use the `add: ` tag at the beginning of any line.                                 |
| patch     | Use the `chg: ` or `rmv: ` or `fix` or `dep: ` tags at the beginning of any line. |
| skip      | Just make regular commit messages with no tags or the words `BREAKING CHANGE`     |

## Incrementing the Major

Example commit message of incrementing the major number:

```text
Removed Update Check

rmv: The command to perform an update check on the changelog since the
update-changelog job handles that too.

add: Ability to turn on/off checkout, as a first step, in the update-changelog
job.

chg: Optimized the ready-for-tagging workflow.

Updated all corresponding documentation.

BREAKING CHANGE
```
NOTE: the `add: `, `rmv: `, `chg: ` tags will be ignored in favor of the
`BREAKING CHANGE` keywords.

## Incrementing the Minor

Example commit message for incrementing the minor number:

```text
add: Parameter to turn on/off checkout, as a first step, in the update-changelog
job.
```

## Incrementing the Patch

Example commit message for incrementing the patch number:

```text
chg: Optimized the update-changelog job to be 2 line instead of 6 in your CI
config.
```
NOTE: The `fix:` and `rmv:` and `deps:` will also increment the patch number.

Know that when the major number is incremented then the minor and patch numbers
will reset to zero. Likewise, when the minor number is incremented, then only
the patch number will reset to zero.

---

[auto-release]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-auto-release
[git-chglog]: https://github.com/git-chglog/git-chglog#table-of-contents