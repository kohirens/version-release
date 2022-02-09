# Setup SSH Keys

SSH keys are needed to give CircleCI jobs the ability to read/write to the repo
for the follow features.

* push to the repo
* make and merge pull request
* publish a release

## Generate An SSH Key for Circle CI

This key will be used to give CircleCI the access it needs to your repository.

These keys are on a per repository basis, a requirement of GitHub. So you'll
have to do this for each repo you want to use these features.

1. Generate the SSH key by running this in a terminal:
   ```
   ssh-keygen -t ed25519 -C "youremail@example.com"
   ```
2. Login to GitHub and go to the repositories settings, then go to "Deploy Keys"
3. Click "add key" and Copy the *.pub value and paste it in
   1. Git it any name you like,
   2. and check the "Write" box
   3. then save.
4. Login to CircleCI then go to the repositories settings.
   1. Go to "SSH Keys"
   2. Click the "Add Key" button,
   3. give it the name "github.com" (it is important to be named after the host it's used for)
   4. paste in the private key
   5. then save
   6. Copy the fingerprint, you will paste this in your Circle CI config (in the near future).
5. Copy the Fingerprint and paste it in your CI config (the future is) now

# Setup A personal access token on GitHub

Unlike the SSH keys which are per repository. You will only need to make 1
token for CirlceCI to use across all the projects the token gives access to
for an Org. It can be risky to generate multiple keys that give access to the
Org. So one should be enough in this case.

This is used to push branches, for updating the CHANGELOG, and then merging it
to your main branch by making a pull request. It is also used to publish a
release.

1. Login into GitHub
2. Go to your profile settings
3. click "Developer Settings" then "personal access tokens"
4. click the "Generate new token"
5. The token should have the following checked:
   ```
     repo
        repo:status
        repo_deployment
        public_repo
        repo:invite
        security_events
    
    write:packages
        read:packages
    
    admin:org
        read:org
   ```

NOTE: Without these certain parts of the jobs may fail.