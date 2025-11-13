# Auto Version Release

Use A CI/CD pipeline to automate updating the CHANGELOG and releasing on GitHub
using [Conventional Commits].

FYI: Changelog updates are performed using the [git-cliff] tool.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/version-release/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/version-release/tree/main) [![CircleCI Orb Version](https://badges.circleci.com/orbs/kohirens/version-release.svg)](https://circleci.com/orbs/registry/orb/kohirens/version-release) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/kohirens/version-release/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

## Features

* Integrate with existing CircleCI and GitHub Actions workflows.
* Automated change log publishing by aggregating Git commit messages that
  follow conventional commits standard, then generating a new CHANGELOG.md with
  those updates; Then automatically commit, branch and then merge the
  CHANGELOG.md into a trunk (named or ref) of your choosing.
* Automated GitHub release publishing and updating the semantic version number.

## Definitions

This application defines "Release" and "Unreleased" as such.

**Release** - An annotated tag that points to Git commit, and any
commits that came before it, in reverse chronological order by committer date.

All commits are considered to be part of that release, until you hit another
annotated tag or the first commit.

**Unreleased** - Are commits that do not fall under an annotated tag.

## What Determines a Release

Conventional commit rely on semantic versioning rules. The commit
messages themselves will contain keywords/markers/labels that control which
number of the semantic version will be incremented based on "Unreleased"
commits.

## Get Started

Make use of this Auto Version Release tool in your new or existing CI/CD
pipelines.

### CircleCI

1. You will need to grant your CircleCI app read/write permissions to your GitHub
repository, you can follow the [CircleCI Setup].
2. Then review the [Orb example] to integrate the necessary changes in your
   `.circleci/config.yml`
3. Then you can try it out by following the [How to Release] guide.

See the [Version Release Orb] documentation for an example.

### GitHub Actions

1. You will need to grant your repo GitHub Actions read/write permissions to
   your GitHub repository, which you can follow the [GitHub Actions] to set that
   up. This also gives an example of how you could integrate the necessary
   changes in your GitHub Actions.
2. Once that is done, then you can try it out by following the [How to Release]
   guide.

## Resources

* See the [Version Release Orb] documentation.
* For contributing, please visit the [Docs] for development details.

---

[Set up a personal access token on GitHub]: /docs/setup-keys.md#setup-a-personal-access-token-on-github
[Version Release Orb]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-examples
[Docs]: /docs/index.md
[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[git-cliff]: https://git-cliff.org/docs/
[Setup Deploy Keys]: /docs/setup-keys.md
[.circleci/config.yml]: /.circleci/config.yml
[CircleCI Setup]: /docs/setup-keys.md#circleci-setup
[GitHub Actions]: /docs/setup-keys.md#github-actions
[How to Release]: /docs/how-to-release.md
[Orb example]: /src/examples/auto-release.yml
