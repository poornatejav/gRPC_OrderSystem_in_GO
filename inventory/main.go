//// package main
////
//// import (
////
////	"context"
////	"log"
////	"net"
////
////	"google.golang.org/grpc"
////	pb "grpc-order-system/proto/inventorypb"
////
//// )
////
////	type server struct {
////		pb.UnimplementedInventoryServiceServer
////	}
////
////	func (s *server) CheckStock(ctx context.Context, req *pb.StockRequest) (*pb.StockResponse, error) {
////		log.Println("Checking stock for:", req.ItemId)
////		inStock := req.ItemId == "item123"
////		return &pb.StockResponse{InStock: inStock}, nil
////	}
////
////	func main() {
////		lis, _ := net.Listen("tcp", ":50053")
////		s := grpc.NewServer()
////		pb.RegisterInventoryServiceServer(s, &server{})
////		log.Println("Inventory Service running on :50053")
////		s.Serve(lis)
////	}
////
//// ------------------------------------------
//// package main
////
//// import (
////
////	"context"
////	"google.golang.org/grpc"
////	"grpc-order-system/proto/inventorypb"
////	"log"
////	"net"
////	"sync"
////
//// )
////
////	type InventoryServer struct {
////		inventorypb.UnimplementedInventoryServiceServer
////		items map[string]int32
////		mu    sync.Mutex
////	}
////
////	func (s *InventoryServer) CheckItemAvailable(ctx context.Context, req *inventorypb.CheckItemRequest) (*inventorypb.CheckItemResponse, error) {
////		s.mu.Lock()
////		defer s.mu.Unlock()
////
////		qty, ok := s.items[req.ItemId]
////		if !ok || qty < req.Quantity {
////			return &inventorypb.CheckItemResponse{Available: false}, nil
////		}
////		return &inventorypb.CheckItemResponse{Available: true}, nil
////	}
////
////	func (s *InventoryServer) DeductItem(ctx context.Context, req *inventorypb.DeductItemRequest) (*inventorypb.DeductItemResponse, error) {
////		s.mu.Lock()
////		defer s.mu.Unlock()
////
////		qty, ok := s.items[req.ItemId]
////		if !ok || qty < req.Quantity {
////			return &inventorypb.DeductItemResponse{Success: false}, nil
////		}
////
////		s.items[req.ItemId] -= req.Quantity
////		return &inventorypb.DeductItemResponse{Success: true}, nil
////	}
////
////	func main() {
////		// Pre-fill some inventory for testing
////		items := map[string]int32{
////			"11": 5,
////			"22": 10,
////		}
////		server := &InventoryServer{items: items}
////
////		lis, err := net.Listen("tcp", ":50052")
////		if err != nil {
////			log.Fatalf("Failed to listen: %v", err)
////		}
////		grpcServer := grpc.NewServer()
////		inventorypb.RegisterInventoryServiceServer(grpcServer, server)
////
////		log.Println("Inventory Service running on :50052")
////		if err := grpcServer.Serve(lis); err != nil {
////			log.Fatalf("Failed to serve: %v", err)
////		}
////	}
////
//// --------------------------------------------
//package main
//
//import (
//	"bufio"
//	"context"
//	"fmt"
//	"log"
//	"net"
//	"os"
//	"strconv"
//	"strings"
//	"sync"
//
//	"google.golang.org/grpc"
//	"grpc-order-system/proto/inventorypb"
//)
//
//type InventoryServer struct {
//	inventorypb.UnimplementedInventoryServiceServer
//	items map[string]int32
//	mu    sync.RWMutex
//}
//
//func (s *InventoryServer) CheckItemAvailable(ctx context.Context, req *inventorypb.CheckItemRequest) (*inventorypb.CheckItemResponse, error) {
//	s.mu.RLock()
//	defer s.mu.RUnlock()
//
//	qty, ok := s.items[req.ItemId]
//	available := ok && qty >= req.Quantity
//	return &inventorypb.CheckItemResponse{Available: available}, nil
//}
//
//func (s *InventoryServer) DeductItem(ctx context.Context, req *inventorypb.DeductItemRequest) (*inventorypb.DeductItemResponse, error) {
//	s.mu.Lock()
//	defer s.mu.Unlock()
//
//	qty, ok := s.items[req.ItemId]
//	if !ok || qty < req.Quantity {
//		return &inventorypb.DeductItemResponse{Success: false}, nil
//	}
//
//	s.items[req.ItemId] -= req.Quantity
//	return &inventorypb.DeductItemResponse{Success: true}, nil
//}
//
//func main() {
//	server := &InventoryServer{
//		items: make(map[string]int32),
//	}
//
//	lis, err := net.Listen("tcp", ":50052")
//	if err != nil {
//		log.Fatalf("Failed to listen: %v", err)
//	}
//
//	grpcServer := grpc.NewServer()
//	inventorypb.RegisterInventoryServiceServer(grpcServer, server)
//
//	// Input goroutine
//	go func() {
//		scanner := bufio.NewScanner(os.Stdin)
//		fmt.Println("Enter inventory items (one per line) in format: <item_id> <quantity>")
//		fmt.Println("Example: item123 10")
//		fmt.Println("Separate by space or comma. Enter empty line to finish or continue adding:")
//
//		for {
//			fmt.Print("Item and quantity: ")
//			if !scanner.Scan() {
//				fmt.Println("\nInput closed.")
//				return
//			}
//			line := strings.TrimSpace(scanner.Text())
//			if line == "" {
//				fmt.Println("Empty input detected, you can continue or Ctrl+C to quit.")
//				continue
//			}
//
//			line = strings.ReplaceAll(line, ",", " ")
//			parts := strings.Fields(line)
//			if len(parts) != 2 {
//				fmt.Println("Invalid input. Please enter exactly 2 values: <item_id> <quantity>")
//				continue
//			}
//
//			itemID := parts[0]
//			qty, err := strconv.Atoi(parts[1])
//			if err != nil || qty < 0 {
//				fmt.Println("Quantity must be a positive integer.")
//				continue
//			}
//
//			server.mu.Lock()
//			server.items[itemID] = int32(qty)
//			server.mu.Unlock()
//
//			fmt.Printf("Added item: %s with quantity %d\n", itemID, qty)
//		}
//	}()
//
//	log.Println("Inventory Service running on :50052")
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("Failed to serve: %v", err)
//	}
//}
//-----------------------------------------------------

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"grpc-order-system/proto/inventorypb"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	collection := client.Database("orderdb").Collection("inventory")
	model := &InventoryModel{Collection: collection}

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter inventory items (format: <item_id> <quantity>):")
		for {
			fmt.Print("Item: ")
			if !scanner.Scan() {
				return
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) != 2 {
				fmt.Println("Invalid format. Use: item123 10")
				continue
			}
			qty, err := strconv.Atoi(parts[1])
			if err != nil || qty < 0 {
				fmt.Println("Invalid quantity.")
				continue
			}
			if err := model.Insert(context.Background(), parts[0], int32(qty)); err != nil {
				fmt.Println("Insert failed:", err)
			} else {
				fmt.Printf("Added item: %s with qty %d\n", parts[0], qty)
			}
		}
	}()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal("listen:", err)
	}

	s := grpc.NewServer()
	inventorypb.RegisterInventoryServiceServer(s, &InventoryService{Model: model})

	log.Println("Inventory Service running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatal("serve:", err)
	}
}
