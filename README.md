# mcp-client
GO MCP client using HTTP transport that supports sampling requests from the server.

## Overview

This client:
- Connects to an MCP server via HTTP/HTTPS transport
- Declares sampling capability during initialization
- Handles incoming sampling requests from the server
- Uses a mock LLM to generate responses (replace with real LLM integration)

## Usage

1. Start an MCP server that supports sampling (e.g., using the `sampling_server` example)

2. Update the server URL in `main.go`:
   ```go
   httpClient, err := client.NewStreamableHttpClient(
       "http://your-server:port", // Replace with your server URL
   )
   ```
