# SCRIPTS.md

This document provides details on the available scripts for building, updating, and synchronizing branches in the repository. Each script automates specific tasks to streamline the development and release workflow.

---

## Table of Contents

- [Script: `release.sh`](#script-releasesh)
  - [Prerequisites](#prerequisites)
  - [What the `release.sh` Script Does](#what-the-releasesh-script-does)
  - [How to Use the Script](#how-to-use-the-script)
  - [Example Commit and Tag Output](#example-commit-and-tag-output)
- [Script: `updateMain.sh` (YMMV)](#script-updatemainsh-ymmv)
  - [Prerequisites](#prerequisites-1)
  - [What the `updateMain.sh` Script Does](#what-the-updatemainsh-script-does)
  - [How to Use the Script](#how-to-use-the-script-1)

---

## Script: `release.sh`

The `release.sh` script automates the process of building, committing, tagging, and pushing changes. It offers options to purge existing GitHub releases and tags, and synchronizes both the current branch and `main`.

### Prerequisites

- **GoReleaser** installed for project building.
- **GitHub CLI** with a valid `GITHUB_TAGS_TOKEN` for purging releases, if using the `--purge-releases` option.
- **jq** and **curl** for purging releases via the GitHub API.

### What the `release.sh` Script Does

- Builds the project using **GoReleaser** without publishing a GitHub release.
- Adds, commits, and tags changes on the current branch and `main`.
- **Optional**: Purges all releases on GitHub (`--purge-releases`) or tags (`--purge-tags`).
  - **Note**: The `GITHUB_TAGS_TOKEN` environment variable must be set for purging releases.
- Generates a timestamp-based tag in the format `vYYYYMMDDHHMMSS` (e.g., `v20241023120000`).
- Updates and pushes the local `main` branch to the remote repository.

### How to Use the Script

1. **Standard Run (without purging releases or tags)**

   ```bash
   ./release.sh
   ```

2. **With Purging Releases**

   ```bash
   ./release.sh --purge-releases=true
   ```

   - **Note**: Requires `GITHUB_TAGS_TOKEN` to be set.

3. **With Purging Tags**

   ```bash
   ./release.sh --purge-tags=true
   ```

4. **With Both Purging Releases and Tags**

   ```bash
   ./release.sh --purge-releases=true --purge-tags=true
   ```

### Example Commit and Tag Output

After running `release.sh`, the following actions occur:

1. **Build**: The project is built using GoReleaser with a snapshot release.
2. **Commit**: All changes are committed with a message like:

   ```
   "Build and release version v20241023120000"
   ```

3. **Tag**: A timestamp-based tag is created, e.g., `v20241023120000`.
4. **Push**: The script pushes:
   - Changes to the current branch.
   - Synchronizes `main` with the current branch.
   - Pushes the new tag to the remote repository.

**Expected Console Output**:

```plaintext
Pushed all changes to current branch 'feature-branch', updated 'main', and tagged the commit as 'v20241023120000'
```

**Note**: If `GITHUB_TAGS_TOKEN` is missing and `--purge-releases=true` is set, the script will output an error and exit.

---

## Script: `updateMain.sh` (YMMV)

The `updateMain.sh` script updates both the current branch and `main`, syncing local changes without creating a release or tag.

### Prerequisites
- **GitHub Access** for pushing updates.
- Ensure you are on the intended development branch with staged changes.

### What the `updateMain.sh` Script Does
- Commits all changes with a timestamped message.
- Synchronizes the local `main` branch with the current branch.
- Pushes both the current branch and updated `main` branch to the remote repository.
- **Does Not** create a tag or trigger a release.

### How to Use the Script

1. **Run from the development branch**:

   ```bash
   ./updateMain.sh
   ```

2. The script will:
   - Commit all staged changes.
   - Update `main` locally to match the current branch.
   - Push changes to the remote repository.
