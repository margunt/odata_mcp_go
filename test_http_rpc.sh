#!/bin/bash

echo "=== Testing OData MCP HTTP/RPC Transport ==="
echo

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Function to send JSON-RPC request
send_request() {
    local method=$1
    local params=$2
    local id=$3
    
    echo -e "${GREEN}Sending $method request...${NC}"
    
    local request="{
        \"jsonrpc\": \"2.0\",
        \"id\": $id,
        \"method\": \"$method\",
        \"params\": $params
    }"
    
    echo "Request:"
    echo "$request" | jq .
    echo
    
    echo "Response:"
    curl -s -X POST http://localhost:8080/rpc \
        -H "Content-Type: application/json" \
        -d "$request" | jq .
    
    echo
    echo "---"
    echo
}

# Check if server is already running
if curl -s http://localhost:8080/health > /dev/null 2>&1; then
    echo "Server is already running on port 8080"
else
    echo "Starting OData MCP server with HTTP transport..."
    ./odata-mcp --transport http --http-addr :8080 https://services.odata.org/V2/Northwind/Northwind.svc/ &
    SERVER_PID=$!
    sleep 2
    
    if ! curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo -e "${RED}Failed to start server${NC}"
        exit 1
    fi
    
    echo "Server started with PID: $SERVER_PID"
fi

echo

# Test initialize
send_request "initialize" '{
    "protocolVersion": "0.1.0",
    "capabilities": {},
    "clientInfo": {
        "name": "HTTP RPC Test Client",
        "version": "1.0.0"
    }
}' 1

# Send initialized notification (no ID for notifications)
echo -e "${GREEN}Sending initialized notification...${NC}"
curl -s -X POST http://localhost:8080/rpc \
    -H "Content-Type: application/json" \
    -d '{
        "jsonrpc": "2.0",
        "method": "initialized",
        "params": {}
    }'
echo "Notification sent (no response expected)"
echo
echo "---"
echo

# Test tools/list
send_request "tools/list" '{}' 2

# Test ping
send_request "ping" '{}' 3

# Clean up if we started the server
if [ ! -z "$SERVER_PID" ]; then
    echo -e "${GREEN}Stopping server...${NC}"
    kill $SERVER_PID 2>/dev/null
fi

echo -e "${GREEN}Test completed!${NC}"