# Developing

Please read this document to help understand developing this Orb.

## Local Development

At this time only the container environment is supported out-of-the-box for
developing locally. The tests require reaching out to __*circleci.com__ and
__*github.com__. The container is configured with these host names so request do
not reach the internet and mock services respond with canned responses. Review
the `.docker/compose.yml` and code in the `avr/mock-server` directories to see
how that is accomplished. To spin up the environment run:

`docker compose up`

There will be 2 container running, but you will only run the test in the
__version-release-web-1__ container. This is the container with the pseudo
__circleci.com__ and __github.com__ mock services.

Login to the mock-server container:

`docker exec -it version-release-web-1 sh`

Then run the test:
```shell
$ go test
```
## When to Update Integration Test

There is a Chicken-Egg/Catch-22 scenario with running integration test with new
Orb features. Because CircleCI must compile the config.yml before it is run, we
don't know the version of dev:<git-hash> to use. For this reason, the
integration test are triggered **after** the pipeline publishes a new dev:alpha
Orb.

This means you cannot run integration test for new features until
after a successful run of "quality-checks" job.

Also dev:* version of Orbs are deleted after 90 days automatically. If test
are not run for 89+, then the test will fail since it will look for dev:aplha.
For this reason, it will point to production version, but the test will generate
a new dev:<git-hash> orb and the integration test will use that. So test should
NOT break after 90 days.

## How to Publish Locally

This is handy to know how to do when there is a problem with dev:* Orbs.

__You are required to have circleci CLI installed locally  and configured
with an CircleCI API key__

```shell
circleci orb pack .\src\ > orb.yml
circleci orb publish .\orb.yml kohirens/version-release@dev:alpha
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

In case you have not registered the namespace (which is only needed once):

```shell
circleci namespace create <orb-namespace> github <github-org>
```

NOTE: If you are having trouble making your namespace with a GitHub org name
then use the option to make it from the Org ID. You can get the org ID from
your Circle CI organizations context page. For example:

```shell
circleci namespace create <your-namespace> --org-id <circleci-org-id>
```
