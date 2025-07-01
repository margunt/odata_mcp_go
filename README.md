# OData MCP Bridge (Go)

A Go implementation of the OData to Model Context Protocol (MCP) bridge, providing universal access to OData services through MCP tools.

This is a Go port of the Python OData-MCP bridge implementation, designed to be easier to run on different operating systems with better performance and simpler deployment. It supports both OData v2 and v4 services.

## Features

- **Universal OData Support**: Works with both OData v2 and v4 services
- **Dynamic Tool Generation**: Automatically creates MCP tools based on OData metadata
- **Multiple Authentication Methods**: Basic auth, cookie auth, and anonymous access
- **SAP OData Extensions**: Full support for SAP-specific OData features including CSRF tokens
- **Comprehensive CRUD Operations**: Generated tools for create, read, update, delete operations
- **Advanced Query Support**: OData query options ($filter, $select, $expand, $orderby, etc.)
- **Function Import Support**: Call OData function imports as MCP tools
- **Flexible Tool Naming**: Configurable tool naming with prefix/postfix options
- **Entity Filtering**: Selective tool generation with wildcard support
- **Cross-Platform**: Native Go binary for easy deployment on any OS

## Installation

### Download Binary

Download the appropriate binary for your platform from the [releases page](https://github.com/odata-mcp/go/releases).

Pre-built binaries are available for:
- Linux (amd64)
- Windows (amd64)
- macOS (Intel and Apple Silicon)

### Build from Source

#### Quick Build (Go required)
```bash
git clone <repository-url>
cd odata_mcp_go
go build -o odata-mcp cmd/odata-mcp/main.go
```

#### Using Makefile (Recommended)
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build for all platforms with WSL integration (copies to /mnt/c/bin)
make build-all-wsl

# Build and test
make dev

# Check current version
make version

# See all options
make help
```

#### Using Build Script
```bash
# Build for current platform
./build.sh

# Build for all platforms
./build.sh all

# See all options
./build.sh help
```

#### Cross-Compilation Examples
```bash
# Using Make
make build-linux     # Linux (amd64)
make build-windows   # Windows (amd64)
make build-macos     # macOS (Intel + Apple Silicon)

# WSL-specific builds (copies Windows binary to /mnt/c/bin)
make build-windows-wsl  # Build Windows + WSL integration
make build-all-wsl      # Build all platforms + WSL integration

# Using build script
./build.sh linux     # Linux (amd64)
./build.sh windows   # Windows (amd64)
./build.sh macos     # macOS (Intel + Apple Silicon)

# Manual Go build
GOOS=linux GOARCH=amd64 go build -o odata-mcp-linux cmd/odata-mcp/main.go
GOOS=windows GOARCH=amd64 go build -o odata-mcp.exe cmd/odata-mcp/main.go
```

#### Docker Build
```bash
# Build Docker image
make docker

# Or manually
docker build -t odata-mcp .

# Run in container
docker run --rm -it odata-mcp --help
```

#### Building in WSL (Windows Subsystem for Linux)

When building in WSL, you can use special targets that automatically copy the Windows binary to your Windows file system:

```bash
# Build all platforms and copy Windows binary to C:\bin
make build-all-wsl

# Build only Windows and copy to C:\bin
make build-windows-wsl
```

Note: These commands will check if `/mnt/c/bin` exists and skip the copy if not found, so they're safe to use on any system.

## Usage

### Claude Desktop Configuration

Claude Desktop uses the stdio transport by default. Here are example configurations:

#### Finding Your Configuration File

The Claude Desktop configuration file location varies by platform:

- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`
- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Linux**: `~/.config/Claude/claude_desktop_config.json`

#### Basic Configuration

```json
{
    "mcpServers": {
        "northwind-v2": {
            "command": "C:/bin/odata-mcp.exe",
            "args": [
                "--service",
                "https://services.odata.org/V2/Northwind/Northwind.svc/",
                "--tool-shrink"
            ]
        },
        "northwind-v4": {
            "command": "C:/bin/odata-mcp.exe",
            "args": [
                "--service",
                "https://services.odata.org/V4/Northwind/Northwind.svc/",
                "--tool-shrink"
            ]
        }
    }
}
```

#### With Authentication

```json
{
    "mcpServers": {
        "my-sap-service": {
            "command": "/usr/local/bin/odata-mcp",
            "args": [
                "--service",
                "https://my-sap-system.com/sap/opu/odata/sap/MY_SERVICE/",
                "--user",
                "myusername",
                "--password",
                "mypassword",
                "--tool-shrink",
                "--entities",
                "Products,Orders,Customers"
            ]
        }
    }
}
```

#### Using Environment Variables (More Secure)

```json
{
    "mcpServers": {
        "my-secure-service": {
            "command": "/usr/local/bin/odata-mcp",
            "args": [
                "--service",
                "https://my-service.com/odata/",
                "--tool-shrink"
            ],
            "env": {
                "ODATA_USERNAME": "myusername",
                "ODATA_PASSWORD": "mypassword"
            }
        }
    }
}
```

**Note:** Claude Desktop currently doesn't support reading environment variables from your system. The `env` field in the configuration sets environment variables specifically for that MCP server process.

#### Security Best Practices for Claude Desktop

1. **Use environment variables** in the `env` field rather than hardcoding credentials in `args`
2. **Limit entity access** using the `--entities` flag to only expose necessary data
3. **Use read-only accounts** when possible for OData services
4. **Store configuration file securely** with appropriate file permissions

**Note:** Claude Desktop does not currently support API key authentication for MCP servers. All MCP servers run locally with the same permissions as Claude Desktop itself.

### Transport Options

The OData MCP bridge supports two transport mechanisms:

1. **STDIO (default)** - Standard input/output communication, used by Claude Desktop
2. **HTTP/SSE** - HTTP server with Server-Sent Events for web-based clients

> ⚠️ **Security Warning**: The HTTP/SSE transport currently does not include authentication. It should only be used in secure, trusted environments such as:
> - Local development (localhost only)
> - Private networks with proper firewall rules
> - Behind a reverse proxy with authentication
> 
> **Do NOT expose the HTTP transport directly to the internet without additional security measures.**

#### Using HTTP/SSE Transport

```bash
# Start server with HTTP transport on default port 8080
./odata-mcp --transport http https://services.odata.org/V2/Northwind/Northwind.svc/

# Use custom port
./odata-mcp --transport http --http-addr :3000 https://services.odata.org/V2/Northwind/Northwind.svc/

# Bind to specific interface
./odata-mcp --transport http --http-addr 127.0.0.1:8080 https://services.odata.org/V2/Northwind/Northwind.svc/
```

When using HTTP transport, the following endpoints are available:

- `GET /health` - Health check endpoint
- `GET /sse` - Server-Sent Events endpoint for real-time communication
- `POST /rpc` - JSON-RPC endpoint for request/response communication

#### Testing HTTP/SSE Transport

1. **Using the provided HTML client:**
   ```bash
   # Start the server
   ./odata-mcp --transport http https://services.odata.org/V2/Northwind/Northwind.svc/
   
   # Open examples/sse_client.html in a web browser
   ```

2. **Using curl:**
   ```bash
   # Test SSE endpoint
   curl -N -H 'Accept: text/event-stream' http://localhost:8080/sse
   
   # Test RPC endpoint
   curl -X POST http://localhost:8080/rpc \
     -H "Content-Type: application/json" \
     -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}'
   ```

3. **Using the test scripts:**
   ```bash
   # Test SSE with interactive script
   ./test_sse.sh
   
   # Test HTTP/RPC communication
   ./test_http_rpc.sh
   ```


### Basic Usage

```bash
# Using positional argument
./odata-mcp https://services.odata.org/V2/Northwind/Northwind.svc/

# Using --service flag
./odata-mcp --service https://services.odata.org/V2/Northwind/Northwind.svc/

# Using environment variable
export ODATA_SERVICE_URL=https://services.odata.org/V2/Northwind/Northwind.svc/
./odata-mcp
```

### Authentication

```bash
# Basic authentication
./odata-mcp --user admin --password secret https://my-service.com/odata/

# Cookie file authentication
./odata-mcp --cookie-file cookies.txt https://my-service.com/odata/

# Cookie string authentication  
./odata-mcp --cookie-string "session=abc123; token=xyz789" https://my-service.com/odata/

# Environment variables
export ODATA_USERNAME=admin
export ODATA_PASSWORD=secret
./odata-mcp https://my-service.com/odata/
```

### Tool Naming Options

```bash
# Use custom prefix instead of postfix
./odata-mcp --no-postfix --tool-prefix "myservice" https://my-service.com/odata/

# Use custom postfix
./odata-mcp --tool-postfix "northwind" https://my-service.com/odata/

# Use shortened tool names
./odata-mcp --tool-shrink https://my-service.com/odata/
```

### Entity and Function Filtering

```bash
# Filter to specific entities (supports wildcards)
./odata-mcp --entities "Products,Categories,Order*" https://my-service.com/odata/

# Filter to specific functions (supports wildcards)  
./odata-mcp --functions "Get*,Create*" https://my-service.com/odata/
```

### Debugging and Inspection

```bash
# Enable verbose output
./odata-mcp --verbose https://my-service.com/odata/

# Trace mode - show all tools without starting server
./odata-mcp --trace https://my-service.com/odata/
```

## Configuration

### Command Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--service` | OData service URL | |
| `-u, --user` | Username for basic auth | |
| `-p, --password` | Password for basic auth | |
| `--cookie-file` | Path to cookie file (Netscape format) | |
| `--cookie-string` | Cookie string (key1=val1; key2=val2) | |
| `--tool-prefix` | Custom prefix for tool names | |
| `--tool-postfix` | Custom postfix for tool names | |
| `--no-postfix` | Use prefix instead of postfix | `false` |
| `--tool-shrink` | Use shortened tool names | `false` |
| `--entities` | Comma-separated entity filter (supports wildcards) | |
| `--functions` | Comma-separated function filter (supports wildcards) | |
| `--sort-tools` | Sort tools alphabetically | `true` |
| `-v, --verbose` | Enable verbose output | `false` |
| `--debug` | Alias for --verbose | `false` |
| `--trace` | Show tools and exit (debug mode) | `false` |

### Environment Variables

| Variable | Description |
|----------|-------------|
| `ODATA_SERVICE_URL` or `ODATA_URL` | OData service URL |
| `ODATA_USERNAME` or `ODATA_USER` | Username for basic auth |
| `ODATA_PASSWORD` or `ODATA_PASS` | Password for basic auth |
| `ODATA_COOKIE_FILE` | Path to cookie file |
| `ODATA_COOKIE_STRING` | Cookie string |

### .env File Support

Create a `.env` file in the working directory:

```env
ODATA_SERVICE_URL=https://my-service.com/odata/
ODATA_USERNAME=admin
ODATA_PASSWORD=secret
```

## Generated Tools

The bridge automatically generates MCP tools based on the OData service metadata:

### Entity Set Tools

For each entity set, the following tools are generated (if the entity set supports the operation):

- `filter_{EntitySet}` - List/filter entities with OData query options
- `count_{EntitySet}` - Get count of entities with optional filter
- `search_{EntitySet}` - Full-text search (if supported by the service)
- `get_{EntitySet}` - Get a single entity by key
- `create_{EntitySet}` - Create a new entity (if allowed)
- `update_{EntitySet}` - Update an existing entity (if allowed)  
- `delete_{EntitySet}` - Delete an entity (if allowed)

### Function Import Tools

Each function import is mapped to an individual tool with the function name.

### Service Information Tool

- `odata_service_info` - Get metadata and capabilities of the OData service

## Examples

### Northwind Service (v2)

```bash
# Connect to the public Northwind OData v2 service
./odata-mcp --trace https://services.odata.org/V2/Northwind/Northwind.svc/

# This will show generated tools like:
# - filter_Products_for_northwind
# - get_Products_for_northwind  
# - filter_Categories_for_northwind
# - get_Orders_for_northwind
# - etc.
```

### Northwind Service (v4)

```bash
# Connect to the public Northwind OData v4 service
./odata-mcp --trace https://services.odata.org/V4/Northwind/Northwind.svc/

# OData v4 is automatically detected and handled appropriately
# Supports v4 specific features like:
# - $count parameter instead of $inlinecount
# - contains() filter function
# - New data types (Edm.Date, Edm.TimeOfDay, etc.)
```

### SAP OData Service

```bash
# Connect to SAP service with CSRF token support
./odata-mcp --user admin --password secret \
  https://my-sap-system.com/sap/opu/odata/sap/SERVICE_NAME/
```

## Differences from Python Version

While maintaining the same CLI interface and functionality, this Go implementation offers:

- **Better Performance**: Native compiled binary with lower memory usage
- **Easier Deployment**: Single binary with no runtime dependencies
- **Cross-Platform**: Native binaries for Windows, macOS, and Linux
- **Type Safety**: Go's type system provides better reliability
- **Simpler Installation**: No need for Python runtime or package management

## Versioning

This project uses automatic versioning based on git tags and commit history:

- **Tagged releases**: Uses git tags (e.g., `v1.0.0`)
- **Development builds**: Uses `0.1.<commit-count>` format
- **Uncommitted changes**: Appends `-dirty` suffix

```bash
# Check current version
make version

# Create a release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

See [VERSIONING.md](VERSIONING.md) for detailed versioning guide.

## Releasing

This project uses automated GitHub Actions for releases. See [RELEASING.md](RELEASING.md) for the release process.

## Security

This project includes comprehensive security measures to prevent credential leaks. See [SECURITY.md](SECURITY.md) for details.

**Important**: Never commit `.zmcp.json` or any files containing real credentials.

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

## License

This project is licensed under the same terms as the original Python implementation.