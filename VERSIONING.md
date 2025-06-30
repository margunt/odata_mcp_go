# Versioning Guide

This project uses automatic versioning based on git tags and commit history.

## How It Works

The Makefile automatically determines the version using the following strategy:

1. **Tagged Releases**: If the current commit has a git tag, that tag is used as the version
2. **Development Builds**: If no tag exists, the version is generated as `0.1.<commit-count>`
3. **Dirty Builds**: If there are uncommitted changes, `-dirty` is appended to the version

## Examples

```bash
# Check current version
make version

# Show detailed build information
make info

# Build with current version
make build
```

## Creating a Release

To create a new release version:

```bash
# Create a tag for the release
git tag -a v1.1.0 -m "Release version 1.1.0"

# Push the tag to remote
git push origin v1.1.0

# Build will now use v1.1.0 as the version
make build
```

## Version Format

- **Release versions**: `v1.0.0`, `v1.1.0`, etc. (following semantic versioning)
- **Development versions**: `0.1.25`, `0.1.26`, etc. (based on commit count)
- **Uncommitted changes**: Append `-dirty` to indicate uncommitted changes

## Checking Version in Binary

The built binary includes version information:

```bash
./odata-mcp --version
```

This will display the version, commit hash, and build time.