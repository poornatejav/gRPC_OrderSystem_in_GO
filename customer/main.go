// package main
//
// import (
//
//	"context"
//	"log"
//	"net"
//
//	pb "grpc-order-system/proto/customerpb"
//	"google.golang.org/grpc"
//
// )
//
//	type server struct {
//	   pb.UnimplementedCustomerServiceServer
//	}
//
//	func (s *server) GetCustomer(ctx context.Context, req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
//	   log.Println("Checking customer:", req.CustomerId)
//	   exists := req.CustomerId == "cust123"
//	   return &pb.CustomerResponse{Exists: exists}, nil
//	}
//
//	func main() {
//	   lis, _ := net.Listen("tcp", ":50051")
//	   s := grpc.NewServer()
//	   pb.RegisterCustomerServiceServer(s, &server{})
//	   log.Println("Customer Service running on :50051")
//	   s.Serve(lis)
//	}

// ---------------------------------------------
// package main
//
// import (
//
//	"context"
//	"log"
//	"net"
//	"strings"
//
//	"google.golang.org/grpc"
//	"grpc-order-system/proto/customerpb"
//
// )
//
//	type CustomerServer struct {
//		customerpb.UnimplementedCustomerServiceServer
//		customers map[string]bool
//	}
//
//	func (s *CustomerServer) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
//		_, exists := s.customers[strings.TrimSpace(req.Id)]
//		return &customerpb.GetCustomerResponse{Exists: exists}, nil
//	}
//
//	func main() {
//		customers := map[string]bool{
//			"123": true,
//			"456": true,
//		}
//
//		server := &CustomerServer{customers: customers}
//
//		lis, err := net.Listen("tcp", ":50051")
//		if err != nil {
//			log.Fatalf("Failed to listen: %v", err)
//		}
//
//		grpcServer := grpc.NewServer()
//		customerpb.RegisterCustomerServiceServer(grpcServer, server)
//
//		log.Println("Customer Service running on :50051")
//		if err := grpcServer.Serve(lis); err != nil {
//			log.Fatalf("Failed to serve: %v", err)
//		}
//	}
//
// ----------------------------------------------
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"
	"grpc-order-system/proto/customerpb"
)

type CustomerServer struct {
	customerpb.UnimplementedCustomerServiceServer
	customers map[string]bool
	mu        sync.RWMutex
}

func (s *CustomerServer) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	id := strings.TrimSpace(req.Id)
	s.mu.RLock()
	_, exists := s.customers[id]
	s.mu.RUnlock()
	return &customerpb.GetCustomerResponse{Exists: exists}, nil
}

func main() {
	server := &CustomerServer{
		customers: make(map[string]bool),
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	customerpb.RegisterCustomerServiceServer(grpcServer, server)

	// Input goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Enter customer IDs (one per line). Enter empty line to finish or continue adding:")

		for {
			fmt.Print("Customer ID: ")
			if !scanner.Scan() {
				fmt.Println("\nInput closed.")
				return
			}
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				fmt.Println("Empty input detected, you can continue or Ctrl+C to quit.")
				continue
			}

			server.mu.Lock()
			server.customers[line] = true
			server.mu.Unlock()
			fmt.Printf("Added customer: %s\n", line)
		}
	}()

	log.Println("Customer Service running on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
