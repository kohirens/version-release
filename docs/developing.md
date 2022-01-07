## Developing

Notes that can help to understand some of the quirks with developing this Orb.

## When to Update Integration Test

There s a Chicken-Egg/Catch-22 scenario when making updates to the integration test.

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