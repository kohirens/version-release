# How To Contribute

Follow these basic steps.

1. Create a fork or request **write** access to this repository.
2. Make a new branch and add your changes there.
3. **Required** Write a git commit message using the [Conventional Commits]
   format (without scope) when you need to update the [CHANGELOG.md]
4. Push your branch.
5. Make a PR when you are ready for a final review.
6. Request a review from the owning team.

NOTE: It is **VERY IMPORTANT** that you follow the [Conventional Commits]
format to add to the [CHANGELOG.md] since it is automatically updated.
This repository uses tools to not only auto-update the [CHANGELOG.md] but also
produce semantic versioned releases. Forgetting the add a proper commit
message or just junk commit messages will fail to produce a release.

__Blank__ format will be used to generate updates to the `CHANGELOG.md` file.

## Versioning

The docker tags are created using the Circle CI Orb [Versions Release] tool.

Just make a PR and add Git commit message using the [Conventional Commits].

Once the PR is merged a new release will be made based on your commit messages.

---

[Conventional Commits]: https://www.conventionalcommits.org/en/v1.0.0/
[CHANGELOG.md]: CHANGELOG.md
[Versions Release]: https://circleci.com/developer/orbs/orb/kohirens/version-release
