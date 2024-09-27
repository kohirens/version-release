# GitHub Auto Version Release

GitHub Auto Version Release or AVR.

1. Workflow Select Workflow.
   1. Merge in a pull-request.
   2. On Merge to main:
      1. Check if the changelog is up-to-date:
         1. If NO, then go to "Publish Changelog Workflow."
      2. Check if there are commits to tag:
         1. If YES, go to the "Tag Changes Workflow."
      3. DONE.
2. Publish Changelog Workflow.
   1. Run Commands to prepend a new section to the changelog.
   2. Commit the changelog.
   3. Push the changelog.
   4. Make a new pull-request and merge it to main.
   5. DONE.
3. Tag Changes Workflow.
   1. Check if there are commits to tag.
   2. Calculate the next tag:
      1. If NO tag, then STOP.
      2. If YES, there is a tag:
         1. Check if GitHub already has a release with such a tag.
            1. If NO, then proceed to the NEXT STEP.
            2. If YES, then report an error and STOP.
   3. Tag the commits by making a Release on GitHub, DONE.
