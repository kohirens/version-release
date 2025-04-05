# Setup Deploy Keys

## CircleCI Setup

A GitHub fine-grained access token is needed to give CircleCI jobs
publish-changelog and tag-and-release the ability to write to the repo
for the follow features:

* Push change for the CHANGELOG.md to the repo, make and merge pull request.
* Publish a release tag.

### Setup A GitHub Personal Access Token

This is used to push branches, for updating the CHANGELOG, and then merging it
to your main branch by making a pull request. It is also used to publish a
release.

1. Login into GitHub
2. Go to your profile settings
3. click "Developer Settings" > "Personal access tokens" > "Fine-grained tokens"
4. click the "Generate new token"
5. Set the "Resource owner" to the organization of the repositories that you
   need write access to.
6. Set "Expiration" from 7 days to "No expiration" (that is the longest you can
   set it at this time of writing).
7. Set "Repository access" to the option that contains the repos you need write
   access to.
8. Select the following "Permissions" options:
   ```
   Permissions
     Commit statuses Access: Read-only
     Contents Access: Read and Write (needed to push files, make commits)
     Metadata Access: Read-only (You cannot turn this off)
     Pull requests Access: Read and Write (needed to open pull-request)
   ```
9. Go to CircleCI App, then select your organization.
10. Select "Organization Settings" > "Context".
11. Select the appropriate context (it should be one that the repos has access
    to).
12. Click the "Add Environment variable" button and paste the fine-grained token
    as the variable `GH_WRITE_TOKEN`.

That finises up this section.

NOTE: Without this setup the jobs will fail in CircleCI.

### Setup A CircleCI API Token

You will also need a CircleCI Token for accessing the API to trigger workflows.

1. Go to your "User Settings" and then select "Personal API Tokens".
2. Click the "Create New Token" button and save it as "CircleCI Automated
   Releases" (or give it any name you like).
3. Save the token as `CIRCLE_TOKEN` in the same context as the GH_WRITE_TOKEN.

NOTE: `GH_WRITE_TOKEN` and `CIRCLE_TOKEN` cannot be changed as they are 
      hardcoded into the AVR application.

## GitHub Actions

### Fined-Grained Access Token

**WARNING:** If you make a new repository after creating a fine-grained token,
you may get 404 when making API requesting such as `repos/{owner}/{repo}/pulls`.
The fix was generating a new token in my case.

1. Go to your settings page and Developer settings > Personal Access Tokens:

   ![Path to fine-grained access tokens](/docs/assets/path-to-fine-grained-tokens.png)

2. Make a fine-grained personal access token with the following permissions:

   ![img.png](/docs/assets/fine-grained-access-permissions.png)

   NOTE: Be sure to select which organization you want to use the token with:

   ![Fine-grained organization selection](/docs/assets/fine-grained-organization-selection.png)

3. If this is for an organization you'll need to have it approved before it can
   access any of those repositories. For details see [Setting a personal access token policy for your organization]
4. Add the token as the `GH_WRITE_TOKEN` secret, either in the repository settings or
   the organization settings.

   NOTE: You need to have a paid plan to use organization secrets in private
   repositories.
5. Pass the secret to the selector job in your workflow.
   ```yaml
   jobs:
     selector: # trigger this on push to main and only when a PR merge.
       uses: kohirens/version-release/.github/workflows/selector.yml
       name: workflow-selector
       secrets:
         github_write_token: ${{ secrets.GH_WRITE_TOKEN }}
   ```

   ![Passing secrets to job step](/docs/assets/passing-secrets.png)
6. Example curl test:
   ```shell
   curl -i -H "Authorization: Bearer ${GH_WRITE_TOKEN}" https://api.github.com/repos/blast-zone/pulls
   ```

NOTE: Even though were running GitHub Actions which will automatically generate
a temporary GITHUB_TOKEN for use in steps; it is prevented from triggering any
additional actions/workflows we may want to run. I assume this is a security
measure and hand-holding for those new to GitHub Actions. However, when we
publish the change log then merge it in, we'll need the workflow-selector to
run once more, automatically, to decided if a release should be published.

---

[this orbs docs]: https://circle`ci.com/developer/orbs/orb/kohirens/version-release
[Setting a personal access token policy for your organization]: https://docs.github.com/en/organizations/managing-programmatic-access-to-your-organization/setting-a-personal-access-token-policy-for-your-organization
