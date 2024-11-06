#!/bin/bash

# Get the current branch name
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)

# Get the current repository in "owner/repo" format
REPO=$(git remote get-url origin | sed -e 's/.*:\/\/[^\/]*\///' -e 's/\.git$//')

# Set variables
TARGET_BRANCH="main"  # The branch to also push to (e.g., main)
TAG="v$(date +'%Y%m%d%H%M%S')"  # Generate a timestamp-based tag
COMMIT_MESSAGE="Build and release version $TAG"

# Parse the --purge-releases and --purge-tags arguments
PURGE_RELEASES=false
PURGE_TAGS=false
for arg in "$@"
do
    case $arg in
        --purge-releases=true)
        PURGE_RELEASES=true
        shift
        ;;
        --purge-tags=true)
        PURGE_TAGS=true
        shift
        ;;
    esac
done

# Function to purge existing GitHub releases
purge_releases() {
  echo "Purging existing releases..."
  # List all releases using the GitHub API
  releases=$(curl -s -H "Authorization: token $GITHUB_TAGS_TOKEN" \
    "https://api.github.com/repos/$REPO/releases" | jq -r '.[].id')

  # Delete each release
  for release_id in $releases; do
    echo "Deleting release ID: $release_id"
    curl -s -X DELETE -H "Authorization: token $GITHUB_TAGS_TOKEN" \
      "https://api.github.com/repos/$REPO/releases/$release_id"
  done

  echo "All releases have been purged."
}

# Function to purge all local and remote tags
purge_tags() {
  echo "Purging all local and remote tags..."

  # Delete all local tags
  git tag -l | xargs git tag -d

  # Fetch remote tags
  git fetch --tags

  # Delete all remote tags
  git tag -l | xargs -n 1 git push --delete origin

  echo "All tags have been purged."
}

# Check if the --purge-releases argument is set to true
if [ "$PURGE_RELEASES" = true ]; then
  # Ensure GITHUB_TAGS_TOKEN is set
  if [ -z "$GITHUB_TAGS_TOKEN" ]; then
    echo "Error: GITHUB_TAGS_TOKEN environment variable must be set to purge releases."
    exit 1
  fi
  purge_releases
fi

# Check if the --purge-tags argument is set to true
if [ "$PURGE_TAGS" = true ]; then
  purge_tags
fi

# Build the binaries and create the release using GoReleaser
goreleaser release --snapshot --clean --skip archive

# Add all files (not just dist/)
git add .

# Commit the changes
git commit -m "$COMMIT_MESSAGE"

# Tag the commit
git tag -a "$TAG" -m "Release version $TAG"

# Update the local main branch to match the current branch without switching branches
git fetch origin $TARGET_BRANCH
git branch -f $TARGET_BRANCH $CURRENT_BRANCH

# Push to the current branch
git push origin "$CURRENT_BRANCH"

# Push the updated local main branch to the remote main branch
git push origin "$TARGET_BRANCH"

# Push the tag to the remote repository
git push origin "$TAG"

echo "Pushed all changes to current branch '$CURRENT_BRANCH', updated '$TARGET_BRANCH', and tagged the commit as '$TAG'"
