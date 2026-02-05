package main

import (
	"context"
	"log"
	"net/http"

	"github.com/dapr/go-sdk/service/common"
	daprd "github.com/dapr/go-sdk/service/http"
)

func main() {
	// Create Dapr service
	server := daprd.NewService(":8080")

	// Health endpoints
	server.AddServiceInvocationHandler("/healthz", func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		return &common.Content{
			Data:        []byte("ok"),
			ContentType: "text/plain",
		}, nil
	})
	server.AddServiceInvocationHandler("/readyz", func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		return &common.Content{
			Data:        []byte("ok"),
			ContentType: "text/plain",
		}, nil
	})

	// Subscribe to user-events topic
	sub := &common.Subscription{
		PubsubName: "shared-pubsub",
		Topic:      "user-events",
		Route:      "/events/user",
	}
	if err := server.AddTopicEventHandler(sub, handleUserEvent); err != nil {
		log.Fatalf("error adding topic handler: %v", err)
	}

	log.Println("Worker service starting on :8080")
	if err := server.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("error starting service: %v", err)
	}
}

func handleUserEvent(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	log.Printf("Received user event: topic=%s, data=%v", e.Topic, string(e.RawData))

	// Process the event
	// - Validate data
	// - Call other services via Dapr service invocation
	// - Update state store
	// - Publish follow-up events

	return false, nil
}
