package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

func main() {
	// Create Dapr service
	server := daprd.NewService(":8080")

	// Health endpoints
	server.AddServiceInvocationHandler("/healthz", healthHandler)
	server.AddServiceInvocationHandler("/readyz", healthHandler)

	// Agent API endpoints
	server.AddServiceInvocationHandler("/api/v1/agents", listAgentsHandler)
	server.AddServiceInvocationHandler("/api/v1/agents/run", runAgentHandler)

	log.Println("Agent service starting on :8080")
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting service: %v", err)
	}
}

func healthHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	return &common.Content{
		Data:        []byte("ok"),
		ContentType: "text/plain",
	}, nil
}

func listAgentsHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	agents := []map[string]interface{}{
		{
			"id":          "agt7890b",
			"name":        "AssistantAgent",
			"description": "General purpose assistant",
			"status":      "active",
		},
		{
			"id":          "agt7890b",
			"name":        "ResearchAgent",
			"description": "Research and analysis agent",
			"status":      "active",
		},
	}

	data, _ := json.Marshal(agents)
	return &common.Content{
		Data:        data,
		ContentType: "application/json",
	}, nil
}

type RunRequest struct {
	Input   string                 `json:"input"`
	Context map[string]interface{} `json:"context"`
}

type RunResponse struct {
	Output    string                 `json:"output"`
	Steps     int                    `json:"steps"`
	Completed bool                   `json:"completed"`
	Context   map[string]interface{} `json:"context"`
}

func runAgentHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	var req RunRequest
	if err := json.Unmarshal(in.Data, &req); err != nil {
		return nil, err
	}

	log.Printf("Running agent with input: %s", req.Input)

	// TODO: Implement actual agent logic
	// 1. Load agent state from memory store
	// 2. Call LLM via Dapr conversation component
	// 3. Execute tools as needed
	// 4. Save state back to memory store
	// 5. Return result

	response := RunResponse{
		Output:    "Agent response to: " + req.Input,
		Steps:     1,
		Completed: true,
		Context:   req.Context,
	}

	data, _ := json.Marshal(response)
	return &common.Content{
		Data:        data,
		ContentType: "application/json",
	}, nil
}
