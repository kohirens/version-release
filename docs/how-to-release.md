# How to Release

Explanation of how to incorporate the auto-release workflow into your
development life cycle, which the Version Release Orb provides.

## Prerequisites

1. Depending on which system you plan to use read on:
   1. [CircleCI Setup]
   2. [GitHub Actions]
2. Add an optional [git-cliff] config to your repo. If you do not, then one will
   be generated and added with the changelog the first time and auto-release
   occurs.
3. Add the [auto-release] workflow to your CI config.
4. When using a private CircleCI Server ensure you set the parameters
   `circleci_api_host` and `circleci_app_host` to the correct values. This will
   allow the workflow-selector workflow to pipeline to trigger workflow on the
   CI server. For example:
   ```yaml
   circleci_api_host="https://api.example-circleci.com"
   circleci_app_host="https://app.example-circleci.com"
   ```
   Otherwise, the jobs will fail since it will look at public Circle CI server.

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

In order to have the changelog updated automatically, please format commits
based on https://www.conventionalcommits.org.

The auto-release workflow depends on git-cliff tool to detect changes in the
repository that need to be release. Please make sure you use a git-cliff
configuration that suites your release needs.

You should know that it only goes back to the last tag, or the beginning if
there are no tags.

For example, prefixing any commit with these commit tags (out-of-the-box) will
cause the release tag to increment the version number accordingly.

| Increment | Description                                                                        |
|-----------|------------------------------------------------------------------------------------|
| major     | Adding the words `BREAKING CHANGE` on a line by itself.                            |
| minor     | Use the `feat: ` tag at the beginning of any line.                                 |
| patch     | Use the `chg: ` or `rmv: ` or `fix:` or `dep: ` tags at the beginning of any line. |
| skip      | See [git-cliff] documentation.                                                     |

## Incrementing the Major

Example commit message of incrementing the major number:

```text
Removed Update Check

rmv: The API function Get().

BREAKING CHANGE: Removed functionality
```

NOTE: All tags will be ignored in favor of the `BREAKING CHANGE` keywords.

## Incrementing the Minor

By default, git-cliff only allows `feat` to increment the minor number
out-of-the-box. Example commit message for incrementing the minor number:

```text
feat: Parameter to turn on/off checkout, as a first step, in the update-changelog
job.
```

## Incrementing the Patch

By default, any tag beside `feat` will increment the patch number.
Example commit message for incrementing the patch number:

```text
chg: Optimized the update-changelog job to be 2 line instead of 6 in your CI
config.
```

Know that when the major number is incremented then the minor and patch numbers
will reset to zero. Likewise, when the minor number is incremented, then only
the patch number will reset to zero.

---

[auto-release]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-auto-release
[git-cliff]: https://git-cliff.org/docs/configuration
[CircleCI Setup]: /docs/setup-keys.md#circleci-setup
[GitHub Actions]: /docs/setup-keys.md#github-actions
