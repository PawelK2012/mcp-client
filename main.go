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
	args["databse"] = "sql"
	args["query"] = "SELECT * FROM masterbranch"
	args["format"] = "json"

	request := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "execute_query",
			Arguments: args,
		},
	}
	result, err := c.CallTool(ctx, request)
	if err != nil {
		log.Printf("tool call failed: %w", err)
	}
	log.Printf("response from MCP %v", result)
}
