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

Again when the "integration-test" pipeline runs it will use the latest dev image set to the latest hash.
However, the whole CI config is compiled before building your latest dev orb (with its latest features). So any calls
to those features will throw errors and the pipeline will fail. So the process above MUST be followed.

## How to Publish Locally

In order to validate the config you'll need the latest changes sometimes. This
is very true when things go wrong. Which is expected in development.
To get to a working state fix the src files. Then publish so that you can
validate the orb, CI config, and get the CI pipeline working as expected.

```shell
circleci orb pack .\src\ > orb.yml
circleci.exe orb publish .\orb.yml  kohirens/version-release@dev:alpha
circleci config validate
```

NOTE: If you're using CircleCI server then set the flag
      `--host https://circleci.example.com` for example, at the end of each
      command.

```shell
circleci orb pack src > orb.yml
circleci orb validate orb.yml
circleci orb publish orb.yml kohirens/version-release@dev:alpha
circleci config validate
```

In case you have not registered the namespace (only needed once):

```shell
circleci namespace create kohirens github kohirens
```

```shell
circleci orb create kohirens/version-release
```