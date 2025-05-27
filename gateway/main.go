package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"grpc-order-system/proto/customerpb"
	"grpc-order-system/proto/inventorypb"
	"grpc-order-system/proto/orderpb"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := customerpb.RegisterCustomerServiceHandlerFromEndpoint(ctx, mux, "localhost:50051", opts)
	if err != nil {
		log.Fatalf("customer: %v", err)
	}
	err = inventorypb.RegisterInventoryServiceHandlerFromEndpoint(ctx, mux, "localhost:50052", opts)
	if err != nil {
		log.Fatalf("inventory: %v", err)
	}
	err = orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, "localhost:50053", opts)
	if err != nil {
		log.Fatalf("order: %v", err)
	}

	log.Println("HTTP Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
