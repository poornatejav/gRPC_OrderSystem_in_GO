// package main
//
// import (
//
//	"context"
//	"fmt"
//	"log"
//	"net"
//	"time"
//
//	orderpb "grpc-order-system/proto/orderpb"
//
//	"google.golang.org/grpc"
//
// )
//
//	type server struct {
//		orderpb.UnimplementedOrderServiceServer
//	}
//
//	func (s *server) PlaceOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
//		customerID := req.GetCustomerId()
//		itemID := req.GetItemId()
//		qty := req.GetQuantity()
//
//		log.Printf("Received Order - CustomerID: %s, ItemID: %s, Quantity: %d", customerID, itemID, qty)
//
//		// TODO: Add calls to customer and inventory services for real validation
//
//		orderID := fmt.Sprintf("order-%d", time.Now().Unix())
//
//		return &orderpb.OrderResponse{
//			Success: true,
//			OrderId: orderID,
//			Message: "Order placed successfully",
//		}, nil
//	}
//
//	func main() {
//		lis, err := net.Listen("tcp", ":50052")
//		if err != nil {
//			log.Fatalf("Failed to listen: %v", err)
//		}
//
//		s := grpc.NewServer()
//		orderpb.RegisterOrderServiceServer(s, &server{})
//
//		log.Println("Order service listening on :50052")
//		if err := s.Serve(lis); err != nil {
//			log.Fatalf("Failed to serve: %v", err)
//		}
//	}
package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"grpc-order-system/proto/customerpb"
	"grpc-order-system/proto/inventorypb"
	"grpc-order-system/proto/orderpb"
)

type OrderServer struct {
	orderpb.UnimplementedOrderServiceServer
	customerClient  customerpb.CustomerServiceClient
	inventoryClient inventorypb.InventoryServiceClient
}

func (s *OrderServer) PlaceOrder(ctx context.Context, req *orderpb.OrderRequest) (*orderpb.OrderResponse, error) {
	log.Printf("Order request: customer=%s, item=%s, qty=%d", req.CustomerId, req.ItemId, req.Quantity)

	// Validate customer
	custResp, err := s.customerClient.GetCustomer(ctx, &customerpb.GetCustomerRequest{Id: req.CustomerId})
	if err != nil || !custResp.Exists {
		return &orderpb.OrderResponse{Status: "Invalid customer"}, nil
	}

	// Validate inventory
	invResp, err := s.inventoryClient.CheckItemAvailable(ctx, &inventorypb.CheckItemRequest{ItemId: req.ItemId, Quantity: req.Quantity})
	if err != nil || !invResp.Available {
		return &orderpb.OrderResponse{Status: "Insufficient inventory"}, nil
	}

	// Deduct inventory
	deductResp, err := s.inventoryClient.DeductItem(ctx, &inventorypb.DeductItemRequest{ItemId: req.ItemId, Quantity: req.Quantity})
	if err != nil || !deductResp.Success {
		return &orderpb.OrderResponse{Status: "Failed to deduct inventory"}, nil
	}

	return &orderpb.OrderResponse{Status: "Order placed successfully"}, nil
}

func main() {
	customerConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Customer connection failed: %v", err)
	}
	defer customerConn.Close()
	customerClient := customerpb.NewCustomerServiceClient(customerConn)

	inventoryConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Inventory connection failed: %v", err)
	}
	defer inventoryConn.Close()
	inventoryClient := inventorypb.NewInventoryServiceClient(inventoryConn)

	server := &OrderServer{
		customerClient:  customerClient,
		inventoryClient: inventoryClient,
	}

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(grpcServer, server)

	log.Println("Order Service running on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
