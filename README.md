# Version Release

Use in CircleCI pipeline to automatically update your CHANGELOG.md file based
on [Conventional Commits] using a [git-cliff] configuration, then auto-publish
a release.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/version-release-orb/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/version-release-orb/tree/main) [![CircleCI Orb Version](https://badges.circleci.com/orbs/kohirens/version-release.svg)](https://circleci.com/orbs/registry/orb/kohirens/version-release) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/kohirens/version-release-orb/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

Provides the following features:

1. Auto update the change log from Git commit messages, then merge it into the
   main trunk.
2. Once the change log is merged, auto publish a release.

## Prerequisites

You will need to integrate CircleCI and GitHub so that Circle CI workflows can
write to your repository on your behalf. To allow GitHub and CircleCI
integration you'll need to make an SSH key/pair to allow CircleCI to write
branches, make pull request, and publish releases to your repository.
* See [Generate An SSH Key for Circle CI] to grant CircleCI write access to
  your repository.
* [Setup A Personal Access Token on GitHub] so that GitHub API request can
  be made from the Circle CI.

## Resources

* See the [Version Release Orb] for examples.
* Visit the [Docs] for development details.

---

[Generate An SSH Key for Circle CI]: https://github.com/kohirens/version-release-orb/blob/main/docs/setup-keys.md#generate-an-ssh-key-for-circle-ci
[Setup A Personal Access Token on GitHub]: https://github.com/kohirens/version-release-orb/blob/main/docs/setup-keys.md#setup-a-personal-access-token-on-github
[Version Release Orb]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-examples
[Docs]: https://github.com/kohirens/version-release-orb/blob/main/docs/index.md
[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[git-cliff]: https://git-cliff.org/docs/
