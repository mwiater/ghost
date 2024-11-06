#!/bin/bash

# Get the current branch name
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Set variables
TARGET_BRANCH="main"  # The branch to also push to (e.g., main)
COMMIT_MESSAGE="Build version $(date +'%Y%m%d%H%M%S')"  # Message with a timestamped version

# Add all files (including new, modified, and deleted ones)
git add .

# Commit the changes
git commit -m "$COMMIT_MESSAGE"

# Update the local main branch to match the current branch without switching branches
git fetch origin $TARGET_BRANCH
git branch -f $TARGET_BRANCH $CURRENT_BRANCH

# Push to the current branch
git push origin "$CURRENT_BRANCH"

# Push the updated local main branch to the remote main branch
git push origin "$TARGET_BRANCH"

echo "Pushed all changes to current branch '$CURRENT_BRANCH' and updated '$TARGET_BRANCH'."
