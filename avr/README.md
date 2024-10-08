# Auto Version Release

This directory serves as the home directory for the `avr` tool, currently
written in Go.

## Options

* `-branch` - Git repository main ref to evaluate, such as commit messages to add to a changelog and or tag (default:
              main).
* `-cicd` - set the CI/CD platform, options are `circleci|github` (default: circleci).
* `-gh-server` - Git hub server (default: github.com).
* `-semver` - Semantic version, required for most commands.
* `-wd` - Working directory (default: ./).

Note: when setting `cicd` to either platform there are required environment
variables that need to be set. For CircleCI you need to set `CIRCLE_TOKEN` and
a few more options.

## Argument

The command has no arguments.

## publish-release-tag Command

`avr -[options] publish-release-tag <semver> <repository-url>`

**Example:**

```shell
avr \
  -branch "main" \
  -gh-server "github.com" \
  -wd "./" \
  -cicd "github"
  publish-release-tag "0.1.0" "git@github.com:kohirens/version-release.git"
```

```shell
avr \
  -branch "main" \
  -gh-server "github.com" \
  -wd "./" \
  -cicd "circleci" \
  publish-release-tag "0.1.0" "${CIRCLE_REPOSITORY_URL}"
```

## References

* [Default environment variables]

---

[Default environment variables]: https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables#default-environment-variables
