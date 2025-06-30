# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- OData v4 support with automatic version detection
- Query parameter translation ($inlinecount to $count for v4)
- Automatic versioning based on git tags and commit history
- GitHub Actions workflows for automated releases
- WSL-specific build targets
- Comprehensive test suite for v4 functionality

### Changed
- Improved response parsing for both v2 and v4 formats
- Enhanced error handling with detailed OData error messages
- Makefile now uses dynamic versioning instead of hardcoded version

### Fixed
- Multiple main function declarations in test files
- Type assertion panics in response parser
- Count value parsing for v2 string responses

## [0.1.0] - 2024-06-30

### Added
- Initial Go implementation of OData MCP Bridge
- Support for OData v2 services
- Dynamic tool generation based on metadata
- Basic auth and cookie authentication
- SAP OData extensions with CSRF token support
- Comprehensive CRUD operations
- Advanced query support with OData query options
- Function import support
- Cross-platform builds for Linux, Windows, and macOS

### Notes
- This is a Go port of the Python OData-MCP bridge
- Maintains CLI compatibility with the original implementation

[Unreleased]: https://github.com/odata-mcp/go/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/odata-mcp/go/releases/tag/v0.1.0