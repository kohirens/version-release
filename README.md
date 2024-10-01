# Auto Version Release

Use A pipeline to automate updating the CHANGELOG and releasing on GitHub using
[Conventional Commits].

FYI: Changelog updates are performed using the [git-cliff] tool.

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/kohirens/version-release/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/kohirens/version-release/tree/main) [![CircleCI Orb Version](https://badges.circleci.com/orbs/kohirens/version-release.svg)](https://circleci.com/orbs/registry/orb/kohirens/version-release) [![GitHub License](https://img.shields.io/badge/license-MIT-lightgrey.svg)](https://raw.githubusercontent.com/kohirens/version-release/master/LICENSE) [![CircleCI Community](https://img.shields.io/badge/community-CircleCI%20Discuss-343434.svg)](https://discuss.circleci.com/c/ecosystem/orbs)

## Features

1. Supports CircleCI pipelines.
2. Auto update the change log from Git commit messages based on conventional
   commits then merge into the main trunk.
3. Auto publish a release tag on GitHub.

## How It Works

**CircleCI Integration**

You will need to integrate CircleCI and GitHub so that Circle CI workflows can
write to your repository on your behalf.

1. [Set up a personal access token on GitHub] to allow CircleCI to make GitHub
 API request on your behalf.
2. [Generate an SSH key for your repository] to grant CircleCI write access to
  your repository.
3. Add or edit `.circleci/config.yml` to:
   1. Update the change log by extracting changes from (conventional) commits
      messages.
   2. Make a pull request and with the updated changes then auto merge them.
   3. Perform a release if the changes warrant one.

See this repository's [.circleci/config.yml] for example.

## Resources

* See the [Version Release Orb] documentation.
* For contributing, please visit the [Docs] for development details.

---

[Generate an SSH key for your repository]: /docs/setup-keyss.md#generate-an-ssh-key-for-circle-ci
[Set up a personal access token on GitHub]: /docs/setup-keys.md#setup-a-personal-access-token-on-github
[Version Release Orb]: https://circleci.com/developer/orbs/orb/kohirens/version-release#usage-examples
[Docs]: /docs/index.md
[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[git-cliff]: https://git-cliff.org/docs/
[Setup Deploy Keys]: /docs/setup-keys.md
[.circleci/config.yml]: /.circleci/config.yml
