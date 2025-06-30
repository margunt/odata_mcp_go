# Release Process

This document describes how to create releases for the OData MCP Bridge.

## Automated Release Process

### Method 1: Push a Tag (Recommended)

The simplest way to create a release is to push a git tag:

```bash
# Create and push a tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# Or use the Makefile
make release TAG=v1.0.0
```

This will trigger the GitHub Actions workflow that:
1. Builds binaries for all platforms
2. Creates release archives
3. Generates checksums
4. Creates a GitHub release with all artifacts

### Method 2: Manual Release via GitHub UI

1. Go to the repository's Actions tab
2. Select "Manual Release" workflow
3. Click "Run workflow"
4. Enter the version (e.g., `v1.0.0`)
5. Choose whether to create as draft
6. Run the workflow

### Method 3: GitHub CLI

If you have the GitHub CLI installed:

```bash
# Create a release with the CLI
gh release create v1.0.0 --generate-notes
```

## Release Workflow Details

### Automatic Builds

When you push a tag starting with `v`, the release workflow:

1. **Builds on multiple platforms**:
   - Linux (amd64)
   - Windows (amd64)
   - macOS (Intel and Apple Silicon)

2. **Creates archives**:
   - `.tar.gz` for Linux and macOS
   - `.zip` for Windows

3. **Generates checksums**:
   - SHA256 checksums for all files

4. **Creates GitHub release**:
   - Includes all binaries
   - Adds installation instructions
   - Publishes as draft for review

### Version Numbering

The project uses semantic versioning (MAJOR.MINOR.PATCH):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes

Development builds use: `0.1.<commit-count>`

## Local Release Process

To create release archives locally:

```bash
# Build all platforms and create archives
make release-local

# Files will be in dist/
ls -la dist/
```

## Pre-release Checklist

Before creating a release:

1. [ ] Update CHANGELOG.md with release notes
2. [ ] Run all tests: `make test`
3. [ ] Build all platforms: `make build-all`
4. [ ] Test binaries on target platforms
5. [ ] Commit all changes
6. [ ] Push to main branch

## Post-release Steps

After the release is created:

1. [ ] Verify all artifacts are uploaded
2. [ ] Test download links
3. [ ] Update documentation if needed
4. [ ] Announce the release

## Troubleshooting

### Release workflow fails

1. Check GitHub Actions logs
2. Ensure Go version is compatible
3. Verify all tests pass locally

### Missing artifacts

1. Check build logs for errors
2. Ensure all platforms built successfully
3. Verify upload steps completed

### Version conflicts

1. Ensure tag doesn't already exist
2. Pull latest changes before tagging
3. Use unique version numbers