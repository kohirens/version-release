# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [4.0.6] - 2024-06-14

### Changed

- Continue To Publish Checkpoint
- Workflow Selector
- Improve User Feedback
- Selecting Publish Changelog Workflow

### Fixed

- No Changelog Config Bug
- Publish Changelog Commit Message
- Change Publish Commit Message

### Regular Maintenance

- Added Debug Messaging
- Code Cleanup
- Updated Development Container

### Testing

- Setting Global SemVer Flag
- Updated

## [4.0.5] - 2024-06-01

### Changed

- Load Variable For Publish Docker Image

## [4.0.4] - 2024-05-28

### Miscellaneous Tasks

- Cleaned Up Stdout/Stderr Messages

## [4.0.3] - 2024-05-16

### Added

- Publish Multiple Images With A Single Job

## [4.0.2] - 2024-05-09

### Fixed

- Error Detecting No Changes To Tag

## [4.0.1] - 2024-05-05

### Fixed

- Checking Changelog Entries

## [4.0.0] - 2024-05-04

### Added

- Help Flag For Publish And Release Command
- Unreleased Changes Detection

### Changed

- Renamed Repository To version-release
- Detecting A Commit Is Tagged
- Checking How The Changelog Is Up-To-Date
- Attach Workspace To Jobs
- Checking How The Changelog Is Up-To-Date
- Semantic Version Tagging
- Detecting Semantic Tag For a Commit
- Distinguish Between A New/Existing Changelog
- Dev Tag To Latest
- Update Git Cliff Configuration
- Calls To Build A Changelog
- Updated Development Environment
- Calls To Build A Changelog
- Updated Development Environment
- Phase Out Obsolete Tools
- Configuring SSH Keys
- Updated Formatting
- Updated GitHub Integration
- Updated the SSH Fingerprint
- Switch To Git-Cliff
- Default Orb Tag
- Version Release Tool Base Image

### Documentation

- Updated How It Works
- Improve User Output

### Fixed

- Executing Publish Release Tag Command
- Publish Changlog Update Loop
- Updating Changelog
- Image Builds
- Error Message

### Miscellaneous Tasks

- Add Clarity For Changelog IsUpToDate Message
- Upgrade Version Release Orb
- Message Output
- Various Formatting Things
- Fixed CI Checkout Image
- Fixed CI Checkout Image
- Fixed Building Dev Image in CI
- Fixed Lint Errors
- Code Refactoring
- Function Checking Unreleased Changes
- Building Changelog
- Configure Changlog Section Types
- Removed Unused CI Configuration
- Updated Test Fixture
- Updated Git Ignore for .ash_history File
- Cleanup CI Configuration

### Removed

- Method of Detecting A Repository Needs A Tag
- Hardcoded Value For Orb Image
- Tag File References
- Obsolete References
- Removed Obsolete Configuration

### Testing

- Fixed Test Fixture
- When Changes are Taggable

<a name="3.1.0"></a>
## [3.1.0] - 2023-08-17
### Added
- Container Image Notes

### Changed
- Update Auto Release Example
- Correct Project Documentation


<a name="3.0.0"></a>
## [3.0.0] - 2023-08-09
### Added
- workflow-selector Job
- Publish Docker Hub Job
- workflow-selector Job
- Version Release Orb Container Image
- Version Release CLI Application

### Changed
- Fixed CI Pipeline
- Over To Use Version Release Orb CLI

### Fixed
- Auto Release


<a name="2.4.1"></a>
## [2.4.1] - 2023-05-26
### Fixed
- Orb Publish On Tag Release


<a name="2.4.0"></a>
## [2.4.0] - 2023-05-25
### Added
- Deley After Publishing Changelog

### Changed
- Fixed Typo In Output

### Removed
- Debug Output


<a name="2.3.5"></a>
## [2.3.5] - 2022-12-05
### Fixed
- Setting CircleCI API Token from Parameter


<a name="2.3.4"></a>
## [2.3.4] - 2022-12-05
### Fixed
- Environment Variable Name In Publish Changelog


<a name="2.3.3"></a>
## [2.3.3] - 2022-12-05
### Changed
- Improve Output Messages

### Fixed
- Job Parameters


<a name="2.3.2"></a>
## [2.3.2] - 2022-12-05
### Changed
- Append gh_token_var Parameter To The publish-changelog Job


<a name="2.3.1"></a>
## [2.3.1] - 2022-12-04
### Changed
- GH_TOKEN Environment Variable Name Check

### Fixed
- Environment Name Parameters For publish-changelog


<a name="2.3.0"></a>
## [2.3.0] - 2022-12-04
### Added
- Parameter CIRCLE_TOKEN Environment Name

### Changed
- publish-changelog Unit Test
- Updated auto-release Example

### Fixed
- YAML Lint Error
- Orb Description and Links


<a name="2.2.3"></a>
## [2.2.3] - 2022-12-03
### Fixed
- Output Grammer


<a name="2.2.2"></a>
## [2.2.2] - 2022-10-01
### Changed
- Make Shell if statements Consistent

### Fixed
- Triggering Github Release


<a name="2.2.1"></a>
## [2.2.1] - 2022-07-18
### Changed
- Upgrade to git-tool-belt version 2.1.0.


<a name="2.2.0"></a>
## [2.2.0] - 2022-07-18
### Added
- Command to check last commit for updated changelog to tag-and-release job.
- Check commit for updated changelog.

### Fixed
- Undo a breaking change.


<a name="2.1.0"></a>
## [2.1.0] - 2022-07-17
### Added
- Gernerate workflows.
- Repository path param to git-chglog-update command.

### Changed
- Refactored CI Dynamic configurations.

### Fixed
- Missing CI parameter from publish-changelog.
- Publishing CHANGELOG.
- Integration test to use a fixture.


<a name="2.0.0-rc1"></a>
## [2.0.0-rc1] - 2022-07-12
### Added
- Parameters to set CircleCI URL for publish trigger.

### Changed
- Updated examples.
- CI config clean up.
- Make note of using CIRCLECI_CLI_TOKEN environment variable.
- Refactored CI trigger workflow code.
- Upgraded CircleCI CLI Orb to version 0.1.9.
- Upgrade Shellcheck Orb to version 3.1.1.
- Updated development doc.
- Updated how to release doc.
- GH unattended login.
- Make github server configurable.
- Updated setup-ssh-keys doc.
- Replaced PARAM_OUTPUTFILE with PARAM_CHANGELOG_FILE.
- Upgraded to git-tool-belt version 1.2.9.
- Upgraded to git-tool-belt 1.2.3 in script test.
- Only trigger a release publish when there are taggable changes.

### Fixed
- Setting dynamic parameter in main dynamic CI config.
- Second GitHub CLI login when publishing a changelog update.
- GitHub CLI login when publishing a changelog update.
- CI location for deploy.yml.
- Broken unit test.
- adding missing git-chglog config.
- Setting GitHub server for tagging a release.
- Used variable before set.
- Looking for trigger.txt when it was not expected to be found.

### Removed
- Obsolete parameters from trigger-tag-and-release command.


<a name="1.1.0"></a>
## [1.1.0] - 2022-05-04
### Added
- Proper tag-and-release example.
- Send API parameter 	riggered-by-bot from trigger-tag-and-release job.
- Example for using trigger-tag-and-release job in a workflow.
- Ability to set the working directory in each job.

### Changed
- Updated the auto-release example.
- Updated publish-changelog example.
- Renamed DoCurl to TriggerPipeline.
- Make checkout and attach_workspace conditional on jobs that require them.
- Make quoting environment variables in jobs consistent.

### Fixed
- Errors in examples.
- Type in trigger-and-tag-release.
- trigger-tag-and-release example.
- Tag and Release job description.


<a name="1.0.6"></a>
## [1.0.6] - 2022-05-01
### Changed
- Updated error message in trigger-tag-and-release job.

### Fixed
- Typo in tag-and-release job.
- Passing VCS type to trigger-tag-and-release job.

### Removed
- Checking CHANGLOG in the trigger-tag-and-release job.


<a name="1.0.5"></a>
## [1.0.5] - 2022-04-29
### Fixed
- Error on persist when no taggable changes.
- Commiting a new Git-ChgLog config.


<a name="1.0.4"></a>
## [1.0.4] - 2022-04-22
### Fixed
- Commit newly generated Git-ChgLog configurations.


<a name="1.0.3"></a>
## [1.0.3] - 2022-04-22
### Changed
- Use GITHUB_TOKEN environment variable in to login to GitHub to make a PR.


<a name="1.0.2"></a>
## [1.0.2] - 2022-04-22

<a name="1.0.1"></a>
## [1.0.1] - 2022-04-21

<a name="1.0.0"></a>
## [1.0.0] - 2022-04-18

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


[Unreleased]: https://github.com/kohirens/version-release-orb/compare/3.1.0...HEAD
[3.1.0]: https://github.com/kohirens/version-release-orb/compare/3.0.0...3.1.0
[3.0.0]: https://github.com/kohirens/version-release-orb/compare/2.4.1...3.0.0
[2.4.1]: https://github.com/kohirens/version-release-orb/compare/2.4.0...2.4.1
[2.4.0]: https://github.com/kohirens/version-release-orb/compare/2.3.5...2.4.0
[2.3.5]: https://github.com/kohirens/version-release-orb/compare/2.3.4...2.3.5
[2.3.4]: https://github.com/kohirens/version-release-orb/compare/2.3.3...2.3.4
[2.3.3]: https://github.com/kohirens/version-release-orb/compare/2.3.2...2.3.3
[2.3.2]: https://github.com/kohirens/version-release-orb/compare/2.3.1...2.3.2
[2.3.1]: https://github.com/kohirens/version-release-orb/compare/2.3.0...2.3.1
[2.3.0]: https://github.com/kohirens/version-release-orb/compare/2.2.3...2.3.0
[2.2.3]: https://github.com/kohirens/version-release-orb/compare/2.2.2...2.2.3
[2.2.2]: https://github.com/kohirens/version-release-orb/compare/2.2.1...2.2.2
[2.2.1]: https://github.com/kohirens/version-release-orb/compare/2.2.0...2.2.1
[2.2.0]: https://github.com/kohirens/version-release-orb/compare/2.1.0...2.2.0
[2.1.0]: https://github.com/kohirens/version-release-orb/compare/2.0.0-rc1...2.1.0
[2.0.0-rc1]: https://github.com/kohirens/version-release-orb/compare/1.1.0...2.0.0-rc1
[1.1.0]: https://github.com/kohirens/version-release-orb/compare/1.0.6...1.1.0
[1.0.6]: https://github.com/kohirens/version-release-orb/compare/1.0.5...1.0.6
[1.0.5]: https://github.com/kohirens/version-release-orb/compare/1.0.4...1.0.5
[1.0.4]: https://github.com/kohirens/version-release-orb/compare/1.0.3...1.0.4
[1.0.3]: https://github.com/kohirens/version-release-orb/compare/1.0.2...1.0.3
[1.0.2]: https://github.com/kohirens/version-release-orb/compare/1.0.1...1.0.2
[1.0.1]: https://github.com/kohirens/version-release-orb/compare/1.0.0...1.0.1
[1.0.0]: https://github.com/kohirens/version-release-orb/compare/0.8.1...1.0.0
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
