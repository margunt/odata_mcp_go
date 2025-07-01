#!/bin/bash

echo "=== Testing OData MCP with SSE Transport ==="
echo

# Start the server with SSE transport in the background
echo "Starting OData MCP server with SSE transport on port 8080..."
./odata-mcp --transport http --http-addr :8080 https://services.odata.org/V2/Northwind/Northwind.svc/ &
SERVER_PID=$!

# Give the server time to start
sleep 2

echo
echo "Server started with PID: $SERVER_PID"
echo "Server is running at http://localhost:8080"
echo
echo "Available endpoints:"
echo "  - SSE endpoint: http://localhost:8080/sse"
echo "  - RPC endpoint: http://localhost:8080/rpc"
echo "  - Health check: http://localhost:8080/health"
echo

# Test health endpoint
echo "Testing health endpoint..."
curl -s http://localhost:8080/health | jq .
echo

echo "To test SSE:"
echo "  1. Open examples/sse_client.html in a web browser"
echo "  2. Or use curl: curl -N -H 'Accept: text/event-stream' http://localhost:8080/sse"
echo

echo "Press Ctrl+C to stop the server..."

# Wait for user to stop
trap "kill $SERVER_PID 2>/dev/null; echo 'Server stopped'; exit" INT
wait $SERVER_PID