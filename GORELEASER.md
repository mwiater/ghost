# GORELEASER.md

This file provides a detailed overview of the GoReleaser setup for this project. GoReleaser automates the process of building and packaging binaries for multiple platforms, allowing for efficient releases.

## Table of Contents

1. [Overview](#overview)
2. [Configuration Breakdown](#configuration-breakdown)
   - [Version and Schema](#version-and-schema)
   - [Pre-Build Hooks](#pre-build-hooks)
   - [Builds](#builds)
3. [Building the Project](#building-the-project)

---

## Overview

The `.goreleaser.yml` configuration file defines the build and release process for this project. It is set up to produce binaries for Windows, Linux (amd64), and Linux (arm64) platforms. GoReleaser is configured to create binaries without CGO (C Go) dependencies, making them portable across different systems.

## Configuration Breakdown

### Version and Schema

```yaml
version: 2
# yaml-language-server: $schema=https://goreleaser.com/static/schema.json
```

- **Version**: This specifies that the GoReleaser configuration file uses version 2 of the configuration format, which includes several enhanced features.
- **Schema**: This modeline allows editors to validate the YAML file against the GoReleaser schema for configuration errors. It helps ensure that the configuration adheres to GoReleaser’s standards.

### Pre-Build Hooks

```yaml
before:
  hooks:
    - go mod tidy
```

- **Hooks**: The `before` hook runs commands before the build process begins. In this setup, `go mod tidy` is executed, which cleans up the `go.mod` and `go.sum` files by removing any unused dependencies. This ensures that the project dependencies are optimized before building.

### Builds

The `builds` section defines specific build configurations for each target platform.

#### Windows 64-bit Build (`win64`)

```yaml
  - id: win64
    env:
      - CGO_ENABLED=0
    goos:
      - windows
    goarch:
      - amd64
    binary: win64/ghost
    no_unique_dist_dir: true
```

- **ID**: `win64` uniquely identifies this build configuration.
- **Environment Variables**: Setting `CGO_ENABLED=0` disables CGO, creating a statically-linked binary. This is important for portability across different Windows environments.
- **Target Platform**: `goos` specifies the OS as `windows`, and `goarch` sets the architecture to `amd64`.
- **Binary Name**: The `binary` field specifies the output binary's name and path. Here, it will be `dist/win64/ghost.exe`.
- **No Unique Distribution Directory**: `no_unique_dist_dir: true` places all build artifacts in the same output directory without additional versioning subfolders. This is useful for consistent binary paths across builds.

#### Linux 64-bit Build (`linux64`)

```yaml
  - id: linux64
    env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - amd64
    binary: linux64/ghost
    no_unique_dist_dir: true
```

- **ID**: `linux64` identifies this build configuration.
- **Environment Variables**: `CGO_ENABLED=0` is also set here, ensuring a statically-linked binary compatible with various Linux distributions.
- **Target Platform**: `goos` is set to `linux`, with `goarch` as `amd64` for 64-bit Linux systems.
- **Binary Name**: The `binary` field directs GoReleaser to output `dist/linux64/ghost`.
- **No Unique Distribution Directory**: This field simplifies access to binaries in the `dist/` directory, keeping organization consistent across releases.

#### Linux ARM 64-bit Build (`linuxarm64`)

```yaml
  - id: linuxarm64
    env:
      - CGO_ENABLED=0
    goos:
      - linux
    goarch:
      - arm64
    binary: linuxarm64/ghost
    no_unique_dist_dir: true
```

- **ID**: `linuxarm64` identifies the ARM build configuration for 64-bit Linux systems.
- **Environment Variables**: As with the other builds, `CGO_ENABLED=0` creates a fully static binary.
- **Target Platform**: `goos` is set to `linux`, and `goarch` is specified as `arm64`, making the binary compatible with ARM architecture.
- **Binary Name**: The output path for this binary is `dist/linuxarm64/ghost`.
- **No Unique Distribution Directory**: Again, this field keeps all binaries in a consistent output path, simplifying access and automation workflows.

---

#### UPX: Enable binary compression

```yaml
upx:
  - enabled: true
    compress: best
    lzma: true
```

You must install UPX on your host platform for this to work. See: https://github.com/alegrey91/go-upx

It's worth installing, there can be significant reductions in binary size:

```
• upx
    • packed                                         before=7.434MB after=2.227MB ratio=29% binary=dist/win64/ghost.exe
    • packed                                         before=7.078MB after=1.926MB ratio=27% binary=dist/linuxarm64/ghost
    • packed                                         before=7.234MB after=2.21MB ratio=30% binary=dist/linux64/ghost
    • took: 21s
```

---



## Building the Project

To build the project with GoReleaser, use the following command:

```bash
goreleaser release --snapshot --clean --skip archive
```

- **`--snapshot`**: Creates a snapshot release, suitable for local testing.
- **`--clean`**: Cleans the build environment before starting.
- **`--skip-archive`**: Skips the archive creation step, producing only the uncompressed binaries in the `dist/` directory.

```
  • skipping announce, archive, publish and validate...
  • cleaning distribution directory
  • loading environment variables
  • getting and validating git state
    • git state                                      commit=d99dca23c32c4b987e44382198b8771ae5a7aa02 branch=development current_tag=v20241111171704 previous_tag=v20241107101326 dirty=true
    • pipe skipped                                   reason=disabled during snapshot mode
  • parsing tag
  • setting defaults
  • snapshotting
    • building snapshot...                           version=20241111171704-SNAPSHOT-d99dca2
  • running before hooks
    • running                                        hook=go mod tidy
  • ensuring distribution directory
  • setting up metadata
  • writing release metadata
  • loading go mod information
  • build prerequisites
  • building binaries
    • building                                       binary=dist\linuxarm64\ghost
    • building                                       binary=dist\linux64\ghost
    • building                                       binary=dist\win64\ghost.exe
  • upx
    • packed                                         before=7.434MB after=2.227MB ratio=29% binary=dist/win64/ghost.exe
    • packed                                         before=7.078MB after=1.926MB ratio=27% binary=dist/linuxarm64/ghost
    • packed                                         before=7.234MB after=2.21MB ratio=30% binary=dist/linux64/ghost
    • took: 21s
  • calculating checksums
  • writing artifacts metadata
  • release succeeded after 38s
  • thanks for using goreleaser!
```

### Expected Output

Running this command will create the following binaries in the `dist/` directory:

```
├── dist
│   ├── artifacts.json
│   ├── config.yaml
│   ├── linux64
│   │   └── ghost
│   ├── linuxarm64
│   │   └── ghost
│   ├── metadata.json
│   └── win64
│       └── ghost.exe
```

Each binary is configured for portability and is compatible with various environments due to the disabled CGO.