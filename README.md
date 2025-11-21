# mcp-client
GO MCP client using HTTP transport that supports sampling requests from the server.

## Overview

This client:
- Connects to an MCP server via HTTP/HTTPS transport
- Declares sampling capability during initialization
- Handles incoming sampling requests from the server
- Uses a mock LLM to generate responses (replace with real LLM integration)

## Usage

1. Start an MCP server that supports sampling 
2. To utilise examples below you can use [database-query-server](https://github.com/PawelK2012/database-query-server)
3. Update the server URL in `main.go`:
```go
httpClient, err := client.NewStreamableHttpClient(
    "http://your-server:port", // Replace with your server URL
)
 ```
4. Start client 
```
go run .
````

# Examples

Copy and paste below examples into `callTool()` in `main.go` file 

## execute_query examples

### execute_query examples - SELECT * FROM table_name

```
# Execute a simple SELECT query
echo '{
  "method": "tools/call",
  "params": {
    "name": "execute_query",
    "arguments": {
      "database": "primary",
      "query": "SELECT * FROM movies",
      "parameters": {},
      "format": "json",
      "limit": 100
    }
  }
}' | mcp-database-query-server

```
This maps to [mcp.CallToolRequest](https://pkg.go.dev/github.com/mark3labs/mcp-go@v0.41.1/mcp#CallToolRequest)

```go
args := make(map[string]interface{})
args["database"] = "postgres"
args["query"] = "SELECT * FROM movies"
args["format"] = "json"
args["limit"] = 100

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "execute_query",
			Arguments: args,
		},
}

```

### execute_query examples - SELECT * FROM table_name

```go

args := make(map[string]interface{})
par["database"] = "postgres"
args["query"] = `CREATE TABLE IF NOT EXISTS movies (
		id SERIAL PRIMARY KEY,
		title VARCHAR(200),
		genre VARCHAR(250),
		year int,
		created TIMESTAMP
)`

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "execute_query",
			Arguments: args,
		},
}

```

### execute_query examples - SELECT * FROM table_name

```go

args := make(map[string]interface{})
par["database"] = "postgres"
args["query"] = `INSERT INTO movies.comedy (title, genre, year)
			VALUES
			('Magellan', 'comedy', 2022),
			('movie 34', 'comedy ', 2012),
			('Mrxxx', 'romantic comedy ', 2011);`

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "execute_query",
			Arguments: args,
		},
}

```

## execute_prepared examples

```go

args := make(map[string]interface{})
var tables []string
args["statement_name"] = "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"
par := []any{"UK", "L%"}
args["parameters"] = par
args["format"] = "json"

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "execute_prepared",
			Arguments: args,
		},
}

```

## get_schema examples

```
# Execute a simple SELECT query
echo '{
  "method": "tools/call",
  "params": {
    "name": "get_schema",
    "arguments": {
      "database": "primary",
      "tables": ["customers"],
      "detailed": true,
    }
  }
}' | mcp-database-query-server

```
This maps to [mcp.CallToolRequest](https://pkg.go.dev/github.com/mark3labs/mcp-go@v0.41.1/mcp#CallToolRequest)


```go
args := make(map[string]interface{})
var tables []string
t := append(tables, "customers")
t := append(tables, "comedy")
args["database"] = "sql"
args["tables"] = t
args["detailed"] = true

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_schema",
			Arguments: args,
		},
}

```

## get_connection_status examples

```
# Execute a simple SELECT query
echo '{
  "method": "tools/call",
  "params": {
    "name": "get_connection_status",
    "arguments": {
      "database": "the name of the DB you want to check the status for",
      "tables": ["customers"],
      "detailed": true,
    }
  }
}' | mcp-database-query-server

```
This maps to [mcp.CallToolRequest](https://pkg.go.dev/github.com/mark3labs/mcp-go@v0.41.1/mcp#CallToolRequest)

```go
args := make(map[string]interface{})
args["database"] = "mcp-query-db"

request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "get_connection_status",
			Arguments: args,
		},
}

```