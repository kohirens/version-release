# Version Release

Automatically update your CHANGELOG and release a tag using
[Conventional Commits].

FYI: Changelog updates are performed using [git-cliff] tool.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/version-release/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/version-release/tree/main) [![CircleCI Orb Version](https://badges.circleci.com/orbs/kohirens/version-release.svg)](https://circleci.com/orbs/registry/orb/kohirens/version-release) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/kohirens/version-release/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

Provides the following features:

1. Auto update the change log from Git commit messages based on conventional
   commits.
2. Auto updated your CHANGELOG.md and merge it into the main trunk.
3. Auto publish a release tag on GitHub.

## Prerequisites

You will need to integrate CircleCI and GitHub so that Circle CI workflows can
write to your repository on your behalf. See the following sections to help you
along.

* See [Generate An SSH Key for Circle CI] to grant CircleCI write access to
  your repository.
* [Setup A Personal Access Token on GitHub] so that GitHub API request can
  be made from the Circle CI.

## Resources

* See the [Version Release Orb] documentation.
* For contributing, please visit the [Docs] for development details.

---

[Generate An SSH Key for Circle CI]: /docs/setup-keyss.md#generate-an-ssh-key-for-circle-ci
[Setup A Personal Access Token on GitHub]: /docs/setup-keys.md#setup-a-personal-access-token-on-github
[Version Release Orb]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-examples
[Docs]: /docs/index.md
[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[git-cliff]: https://git-cliff.org/docs/
[Setup Deploy Keys]: /docs/setup-keys.md
