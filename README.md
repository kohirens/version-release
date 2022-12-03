# Version Release

[![CircleCI Build Status](https://circleci.com/gh/kohirens/version-release-orb.svg?style=shield "CircleCI Build Status")](https://circleci.com/gh/kohirens/version-release-orb) [![CircleCI Orb Version](https://badges.circleci.com/orbs/kohirens/version-release.svg)](https://circleci.com/orbs/registry/orb/kohirens/version-release) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/kohirens/version-release-orb/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

Provides the following features:

1. Auto update the changelog from commit messages, then merge it into the main trunk
2. Once the changelog is merged, publish a release on GitHub.

## Pre-requisites

You will need to integrate CircleCI and GitHub so that Circle CI can perform
actions in GitHub on your behave. This is the automation part.

* Make an SSH key/pair to allow Circle Ci to publish releases on your
  repository in GitHub. See [Generate An SSH Key for Circle CI]
* Make a GitHub Personal Access token so that Circle CI can make branches,
  PRs, and merge changes for the CHANGELOG.md. See [Setup A Personal Access Token on GitHub]

## Resources

See the [Version Release Orb] for examples.

Visit the [Docs] for development details.

---

[Generate An SSH Key for Circle CI]: /docs/setup-keys.md#generate-an-ssh-key-for-circle-ci
[Setup A Personal Access Token on GitHub]: /docs/setup-keys.md#setup-a-personal-access-token-on-github
[Version Release Orb]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-examples
[Docs]: /docs/index.md