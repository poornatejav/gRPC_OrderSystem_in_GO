// package main
//
// import (
//
//	"bufio"
//	"context"
//	"fmt"
//	"log"
//	"os"
//	"strconv"
//	"strings"
//	"time"
//
//	orderpb "grpc-order-system/proto/orderpb"
//
//	"google.golang.org/grpc"
//
// )
//
//	func main() {
//		conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
//		if err != nil {
//			log.Fatalf("Failed to connect to order service: %v", err)
//		}
//		defer conn.Close()
//
//		client := orderpb.NewOrderServiceClient(conn)
//
//		reader := bufio.NewReader(os.Stdin)
//
//		fmt.Println("Place your order")
//
//		fmt.Print("Customer ID: ")
//		custID, _ := reader.ReadString('\n')
//		custID = strings.TrimSpace(custID)
//
//		fmt.Print("Item ID: ")
//		itemID, _ := reader.ReadString('\n')
//		itemID = strings.TrimSpace(itemID)
//
//		fmt.Print("Quantity: ")
//		qtyStr, _ := reader.ReadString('\n')
//		qtyStr = strings.TrimSpace(qtyStr)
//		qty, err := strconv.Atoi(qtyStr)
//		if err != nil {
//			log.Fatalf("Invalid quantity: %v", err)
//		}
//
//		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//		defer cancel()
//
//		req := &orderpb.OrderRequest{
//			CustomerId: custID,
//			ItemId:     itemID,
//			Quantity:   int32(qty),
//		}
//
//		resp, err := client.PlaceOrder(ctx, req)
//		if err != nil {
//			log.Fatalf("PlaceOrder failed: %v", err)
//		}
//
//		if resp.Success {
//			fmt.Printf("Order placed! Order ID: %s\n", resp.OrderId)
//		} else {
//			fmt.Printf("Order failed: %s\n", resp.Message)
//		}
//	}
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"grpc-order-system/proto/orderpb"
)

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer conn.Close()

	client := orderpb.NewOrderServiceClient(conn)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Place your order")

		fmt.Print("Customer ID: ")
		custID, _ := reader.ReadString('\n')
		custID = strings.TrimSpace(custID)

		fmt.Print("Item ID: ")
		itemID, _ := reader.ReadString('\n')
		itemID = strings.TrimSpace(itemID)

		fmt.Print("Quantity: ")
		qtyStr, _ := reader.ReadString('\n')
		qtyStr = strings.TrimSpace(qtyStr)
		qty, err := strconv.Atoi(qtyStr)
		if err != nil {
			fmt.Println("Invalid quantity, try again")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		resp, err := client.PlaceOrder(ctx, &orderpb.OrderRequest{
			CustomerId: custID,
			ItemId:     itemID,
			Quantity:   int32(qty),
		})
		if err != nil {
			log.Printf("Error placing order: %v", err)
			continue
		}

		fmt.Printf("Order response: %s\n\n", resp.Status)
	}
}
