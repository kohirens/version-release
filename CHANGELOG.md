<a name="unreleased"></a>
## [Unreleased]


<a name="0.8.1"></a>
## [0.8.1] - 2022-03-08

<a name="0.8.0"></a>
## [0.8.0] - 2022-02-22
### Added
- File to export the release tag.


<a name="0.7.6"></a>
## [0.7.6] - 2022-02-20
### Changed
- Upgraded git-tool-belt to version 0.9.0.


<a name="0.7.5"></a>
## [0.7.5] - 2022-02-20
### Fixed
- Allow tagging repos with no tags.


<a name="0.7.4"></a>
## [0.7.4] - 2022-02-18
### Fixed
- Use of git alias not defined.


<a name="0.7.3"></a>
## [0.7.3] - 2022-02-12
### Changed
- Removed obsolite code.

### Removed
- Square brackets from GitHub release tags.


<a name="0.7.2"></a>
## [0.7.2] - 2022-02-11
### Changed
- merge-changelog scritp to exit when gen branch exist remotely.


<a name="0.7.1"></a>
## [0.7.1] - 2022-02-09
### Changed
- tag-and-release.sh verbosity.
- Continue CHANGELOG merge when new and not in the repo.

### Fixed
- tag-and-release not resetting to correct commit.
- Processing empty commit txt file.


<a name="0.7.0"></a>
## [0.7.0] - 2022-02-03
### Added
- Check that merging changelog has completed.
- Auto release workflow example.

### Changed
- Refactored merge-changelog.sh.
- Add a newline in the changelog commit message.
- Removed failed attempt at newlines in commmit message.
- Updated merge changelog output for user feedback.
- Placed required checkout step in jobs.

### Fixed
- Error when missing commit-to-tag.txt file.
- tag-and-release job.
- typo in command.
- No auto tagging.


<a name="0.6.2"></a>
## [0.6.2] - 2022-01-31
### Changed
- Workflows into a single auto-release workflow.
- Renamed script that merged the changelog.


<a name="0.6.1"></a>
## [0.6.1] - 2022-01-19
### Fixed
- Setting revision range for tag-and-release job.


<a name="0.6.0"></a>
## [0.6.0] - 2022-01-19
### Changed
- Turn on verbosity when tagging.
- Default Executor image to 0.8.0.
- Upgraded executor image.
- More logging in Tag and Release checks.
- Exit normal when no change in the changelog.
- Moved tag-and-release to script for linting.
- Updated Executore image to 0.7.0.

### Fixed
- Fetch all refs before isTaggable check.
- Tag and release checks.


<a name="0.5.4"></a>
## [0.5.4] - 2022-01-08

<a name="0.5.3"></a>
## [0.5.3] - 2022-01-08

<a name="0.5.2"></a>
## [0.5.2] - 2022-01-08

<a name="0.5.1"></a>
## [0.5.1] - 2022-01-08

<a name="0.5.0"></a>
## [0.5.0] - 2022-01-07
### Added
- Initial docs.
- Job for triggering a release tag.

### Changed
- Split changelog update and tagging.

### Fixed
- Missing checkout step in tag-and-release job.
- Minor and patch version not resetting.
- Setting TOKEN to auth trigger.
- Example documentation.


<a name="0.4.7"></a>
## [0.4.7] - 2022-01-07

<a name="0.4.6"></a>
## [0.4.6] - 2022-01-07
### Added
- Flag to auto generate title and notes changelog merge.
- Tag auto release with title.

### Changed
- Fetch the latest main branch to tag after merging the changelog.
- Reduced time to merge changelog updates.
- Turn auto merge of changelog update back on.
- Updated changelog command to use git-tool-belt.
- Upgrade default executor image.
- Renamed job to update-and-merge-changelog.
- Only run tag-and-release-flow on main branch.
- Refactor tag-and-release job and workflow.

### Fixed
- Passing merge type to update and merge changelog command.
- commit-and-merge-changelog job.

### Removed
- Double quotes from the release tag.
- Skip CI tag in commit message for changelog update.


<a name="0.3.6"></a>
## [0.3.6] - 2022-01-02
### Added
- Updated CHANGELOG command example.


<a name="0.3.5"></a>
## [0.3.5] - 2022-01-02

<a name="0.3.4"></a>
## [0.3.4] - 2022-01-02

<a name="0.3.3"></a>
## [0.3.3] - 2022-01-02
### Fixed
- Dupe CI runs.
- CI auto Orb publish after tagging.


<a name="0.3.1"></a>
## [0.3.1] - 2022-01-02

<a name="0.3.2"></a>
## [0.3.2] - 2022-01-02
### Changed
- Allow tags to be run in CI deploy workflow.


<a name="0.3.0"></a>
## [0.3.0] - 2022-01-02
### Added
- Publish on making a new tag.
- Check for tags before updating CHANGELOG.md.
- SSH fingerprint to CI config.

### Fixed
- Type in default executor.


<a name="0.2.0"></a>
## 0.2.0 - 2021-12-26
### Added
- Docker environment.
- CHANGELOG update command.
- git-chglog check command.

### Changed
- Point to currect script for git-chglog-check command.
- Moved the default executer to an Alpine image.

### Fixed
- image access denied not found in CI.
- YAML lint issue.
- Using incorrect image with orb command.

### Removed
- Sample command and job.


[Unreleased]: https://github.com/kohirens/version-release-orb/compare/0.8.1...HEAD
[0.8.1]: https://github.com/kohirens/version-release-orb/compare/0.8.0...0.8.1
[0.8.0]: https://github.com/kohirens/version-release-orb/compare/0.7.6...0.8.0
[0.7.6]: https://github.com/kohirens/version-release-orb/compare/0.7.5...0.7.6
[0.7.5]: https://github.com/kohirens/version-release-orb/compare/0.7.4...0.7.5
[0.7.4]: https://github.com/kohirens/version-release-orb/compare/0.7.3...0.7.4
[0.7.3]: https://github.com/kohirens/version-release-orb/compare/0.7.2...0.7.3
[0.7.2]: https://github.com/kohirens/version-release-orb/compare/0.7.1...0.7.2
[0.7.1]: https://github.com/kohirens/version-release-orb/compare/0.7.0...0.7.1
[0.7.0]: https://github.com/kohirens/version-release-orb/compare/0.6.2...0.7.0
[0.6.2]: https://github.com/kohirens/version-release-orb/compare/0.6.1...0.6.2
[0.6.1]: https://github.com/kohirens/version-release-orb/compare/0.6.0...0.6.1
[0.6.0]: https://github.com/kohirens/version-release-orb/compare/0.5.4...0.6.0
[0.5.4]: https://github.com/kohirens/version-release-orb/compare/0.5.3...0.5.4
[0.5.3]: https://github.com/kohirens/version-release-orb/compare/0.5.2...0.5.3
[0.5.2]: https://github.com/kohirens/version-release-orb/compare/0.5.1...0.5.2
[0.5.1]: https://github.com/kohirens/version-release-orb/compare/0.5.0...0.5.1
[0.5.0]: https://github.com/kohirens/version-release-orb/compare/0.4.7...0.5.0
[0.4.7]: https://github.com/kohirens/version-release-orb/compare/0.4.6...0.4.7
[0.4.6]: https://github.com/kohirens/version-release-orb/compare/0.3.6...0.4.6
[0.3.6]: https://github.com/kohirens/version-release-orb/compare/0.3.5...0.3.6
[0.3.5]: https://github.com/kohirens/version-release-orb/compare/0.3.4...0.3.5
[0.3.4]: https://github.com/kohirens/version-release-orb/compare/0.3.3...0.3.4
[0.3.3]: https://github.com/kohirens/version-release-orb/compare/0.3.1...0.3.3
[0.3.1]: https://github.com/kohirens/version-release-orb/compare/0.3.2...0.3.1
[0.3.2]: https://github.com/kohirens/version-release-orb/compare/0.3.0...0.3.2
[0.3.0]: https://github.com/kohirens/version-release-orb/compare/0.2.0...0.3.0
