package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

// MockSamplingHandler implements client.SamplingHandler for demonstration.
// In a real implementation, this would integrate with an actual LLM API.
type MockSamplingHandler struct{}

func (h *MockSamplingHandler) CreateMessage(ctx context.Context, request mcp.CreateMessageRequest) (*mcp.CreateMessageResult, error) {
	// Extract the user's message
	if len(request.Messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	// Get the last user message
	lastMessage := request.Messages[len(request.Messages)-1]
	userText := ""
	if textContent, ok := lastMessage.Content.(mcp.TextContent); ok {
		userText = textContent.Text
	}

	// Generate a mock response
	responseText := fmt.Sprintf("Mock LLM response to: '%s'", userText)

	log.Printf("Mock LLM generating response: %s", responseText)

	result := &mcp.CreateMessageResult{
		SamplingMessage: mcp.SamplingMessage{
			Role: mcp.RoleAssistant,
			Content: mcp.TextContent{
				Type: "text",
				Text: responseText,
			},
		},
		Model:      "mock-model-v1",
		StopReason: "endTurn",
	}

	return result, nil
}

func main() {
	// Create sampling handler
	samplingHandler := &MockSamplingHandler{}

	// Create HTTP transport directly
	httpTransport, err := transport.NewStreamableHTTP(
		"http://localhost:8080/mcp", // Replace with your MCP server URL
		// You can add HTTP-specific options here like headers, OAuth, etc.
	)
	if err != nil {
		log.Fatalf("Failed to create HTTP transport: %v", err)
	}
	defer httpTransport.Close()

	// Create client with sampling support
	c := client.NewClient(
		httpTransport,
		client.WithSamplingHandler(samplingHandler),
	)

	// Start the client
	ctx := context.Background()
	err = c.Start(ctx)
	if err != nil {
		log.Fatalf("Failed to start client: %v", err)
	}

	// Initialize the MCP session
	initRequest := mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			Capabilities:    mcp.ClientCapabilities{
				// Sampling capability will be automatically added by the client
			},
			ClientInfo: mcp.Implementation{
				Name:    "mcp-http-client",
				Version: "1.0.0",
			},
		},
	}

	_, err = c.Initialize(ctx, initRequest)
	if err != nil {
		log.Fatalf("Failed to initialize MCP session: %v", err)
	}

	fmt.Println("Performing health check...")
	if err := c.Ping(ctx); err != nil {
		log.Fatalf("Health check failed: %v", err)
	}
	fmt.Println("Server is alive and responding")

	toolsRequest := mcp.ListToolsRequest{}
	tools, err := c.ListTools(ctx, toolsRequest)
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
		return
	}

	fmt.Printf("Available tools: %d\n", len(tools.Tools))
	for _, tool := range tools.Tools {
		log.Printf("- %s: %s\n", tool.Name, tool.Description)
	}

	callTool(ctx, c)

	// In a real application, you would keep the client running to handle sampling requests
	// For this example, we'll just demonstrate that it's working

	// Keep the client running (in a real app, you'd have your main application logic here)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Println("Client context cancelled")
	case <-sigChan:
		log.Println("Received shutdown signal")
	}
}

// calling MCP server tool request
func callTool(ctx context.Context, c *client.Client) {
	args := make(map[string]interface{})

	/////////////
	//execute_query examples

	//par := make(map[string]any)
	// par["tables"] = "customers"
	//par["2"] = "L%"
	// par := `{"tables": "customers"}`
	// args["database"] = "sql"
	args["query"] = "SELECT * FROM pg_settings"
	// args["query"] = "SELECT id, CustomerName, ContactName FROM Customers WHERE Country = $1"
	//  SELECT CustomerName, ContactName FROM customers WHERE Country = 'UK';
	// args["query"] = "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"
	// args["query"] = "SELECT * FROM customers"
	//args["query"] = "select column_name, data_type, character_maximum_length from INFORMATION_SCHEMA.COLUMNS where table_name =$1;"
	// args["query"] = "CREATE SCHEMA movies;"
	// args["query"] = `CREATE TABLE IF NOT EXISTS movies.comedy (
	// 	id SERIAL PRIMARY KEY,
	// 	title VARCHAR(200),
	// 	genre VARCHAR(250),
	// 	year int,
	// 	created TIMESTAMP
	// )`
	// args["query"] = `INSERT INTO movies.comedy (title, genre, year)
	// 		VALUES
	// 		('fMagellan', 'comedy', 2022),
	// 		('movie 34', 'comedy ', 2012),
	// 		('Mrxxx', 'romantic comedy ', 2011);`
	// args["query"] = `INSERT INTO Users (CustomerName, ContactName, Address, City, PostalCode, Country)
	// 		VALUES
	// 		('mr x', 'Tmrxxxx', 'Warsaw strasse', 'Berlin', '40555506', 'Germany'),
	// 		('Pablo', 'Pablo picasso', 'Costa del sol', 'Malaga', 'mal22', 'Spaim'),
	// 		('Mr H', 'h', 'Most posh place in Barcelona', 'Barcelona', 'bc 0AA', 'Spain');`
	// args["parameters"] = par
	args["format"] = "csv"

	/////////////
	//execute_prepared

	// args["statement_name"] = "SELECT CustomerName, Address FROM customers WHERE id =$1;"
	// args["statement_name"] = "SELECT CustomerName, Address FROM customers WHERE Country =$1 AND City LIKE $2"
	// args["format"] = "json"
	// par := []any{"UK", "L%"}
	// args["parameters"] = par

	/////////////
	// get_schema examples

	// var tables []string
	// t := append(tables, "customers")
	// t := append(tables, "comedy")
	// args["database"] = "sql"
	// args["tables"] = t
	// args["detailed"] = true

	/////////////
	// get_connection_status

	args["database"] = "mcp-query-db"
	args["connected"] = true
	args["pool_stats"] = true

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			// Name: "execute_prepared",
			// Name: "execute_query",
			// Name:  "get_schema",
			Name:      "get_connection_status",
			Arguments: args,
		},
	}
	result, err := c.CallTool(ctx, request)
	if err != nil {
		log.Printf("tool call failed: %v", err)
	}
	log.Printf("response from MCP %+v \n", result)
}
