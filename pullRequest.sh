#!/bin/bash

# Temporary exit to prevent the script from running; remove this line after testing
exit 0

# Check if the working directory is clean
if [[ -n $(git status --porcelain) ]]; then
  echo "Error: A clean working tree is needed. Please commit or stash your changes."
  exit 1
fi

# Commit and push with timestamped message
COMMIT_MESSAGE="Build version $(date +'%Y%m%d%H%M%S')"
git add .
git commit -m "$COMMIT_MESSAGE"
git push

# Get the current branch name
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Create a new branch name based on the last commit message, lowercased and slugified
LAST_COMMIT_MSG=$(git log -1 --pretty=%B)
NEW_BRANCH=$(echo "$LAST_COMMIT_MSG" | tr '[:upper:]' '[:lower:]' | sed -e 's/[^a-z0-9]/-/g' | sed -e 's/--*/-/g')

# Checkout the new branch based on the current branch's latest commit
git checkout -b "$NEW_BRANCH"

# Push the new branch
git push -u origin "$NEW_BRANCH"

# Create a pull request to main using GitHub CLI with autofill
gh pr create --base main --head "$NEW_BRANCH" --fill

# Switch back to the original branch
git checkout "$CURRENT_BRANCH"

# Delete the pull request branch locally and remotely
git branch -D "$NEW_BRANCH"
git push origin --delete "$NEW_BRANCH"

echo "Pull request created from $NEW_BRANCH to main, switched back to $CURRENT_BRANCH, and deleted $NEW_BRANCH."
