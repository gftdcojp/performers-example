package main

import (
	"log"
	"net"
	"net/http"

	daprd "github.com/dapr/go-sdk/service/grpc"
)

func main() {
	// Health server on HTTP (for K8s probes)
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("ok"))
		})
		log.Println("Health server starting on :8080")
		http.ListenAndServe(":8080", mux)
	}()

	// Dapr gRPC server
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server, err := daprd.NewService(lis)
	if err != nil {
		log.Fatalf("failed to create service: %v", err)
	}

	// Register actors
	server.RegisterActorImplFactoryContext(UserActorFactory)
	server.RegisterActorImplFactoryContext(SessionActorFactory)
	server.RegisterActorImplFactoryContext(CartActorFactory)

	log.Println("Actor service starting on :50051")
	if err := server.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

// Actor factories (simplified)
func UserActorFactory() any    { return &UserActor{} }
func SessionActorFactory() any { return &SessionActor{} }
func CartActorFactory() any    { return &CartActor{} }

// Actor types (simplified)
type UserActor struct{}
type SessionActor struct{}
type CartActor struct{}

func (a *UserActor) Type() string    { return "UserActor" }
func (a *SessionActor) Type() string { return "SessionActor" }
func (a *CartActor) Type() string    { return "CartActor" }
