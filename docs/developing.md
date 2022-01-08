## Developing

Notes that can help to understand developing this Orb.

## When to Update Integration Test

There is a Chicken-Egg/Catch-22 scenario when making updates to the integration test.

So when the "test-pack" pipeline runs "dev-orb-version" will be "dev:alpha" from any previous publish (or not exist in
the case of a new Orb). Only after a successful run of "test-pack" (or manual push) does a new version of your Orb
exist with the latest features.

This means you cannot use new features in the integration test pipeline until after a successful run of "test-pack".

Updating integration test requires the following process:
1. Push out new features/breaking changes without updating the integration pipeline.
2. Make sure the pipeline succeeds.
3. Update your integration pipeline to test the new features/breaking changes.

Again when the "integration-test_deploy" pipeline runs it will use the latest dev image set to the latest hash.
However, the whole CI config is compiled before building your latest dev orb (with its latest features). So any calls
to those features will throw errors and the pipeline will fail. So the process above MUST be followed.

## How to Publish

In order to have a new version of this orb published, perform the following:

1. Create and push a branch with your new features.
2. Develop your features.
3. Push your branch to GitHub and make a PR to the main trunk of this repo.
4. Once merged there will be a new release published within 10 minute (assuming normal functionality of GitHub and CircleCI).
* When ready to publish a new production version, create a Pull Request from _feature branch_ to `main`.
* The first line of the commit message, which some have dubbed the __commit subject__, must contain a special semver
  tag: `[semver:<segment>]` where `<segment>` is replaced by one of the following values.

Add inay of the keywords in the Descriptions below, an in any commit message since the last release, to increment the
version number accordingly.

| Increment | Description                                                                       |
|-----------|-----------------------------------------------------------------------------------|
| major     | Adding the words `BREAKING CHANGE`                                                |
| minor     | Use the `add: ` tag at the beginning of any line.                                 |
| patch     | Use the `chg: ` or `rmv: ` or `fix` or `dep: ` tags at the beginning of any line. |
| skip      | Just make regular commit messages with no tags or the words `BREAKING CHANGE`    |

### Incrementing the Major

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
NOTE: the `add: `, `rmv: `, `chg: ` tags will be ignored in favor of the `BREAKING CHANGE` keywords.

### Incrementing the Minor

Example commit message for incrementing the minor number:

```text
add: Parameter to turn on/off checkout, as a first step, in the update-changelog
job.
```

### Incrementing the Patch

Example commit message for incrementing the patch number:

```text
chg: Optimized the update-changelog job to be 2 line instead of 6 in your CI
config.
```

NOTE: When the major number is incremented then the minor and patch numbers will reset to zero. Likewise, when the
minor number is incremented, then only the patch number will reset to zero.
