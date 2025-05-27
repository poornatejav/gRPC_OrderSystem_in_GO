# gRPC Order Management System in Go

This repository contains a sample gRPC-based Order Management System written in Go. It demonstrates microservices architecture using gRPC, Protocol Buffers, and Go modules, featuring separate services for handling customers, inventory, and orders, along with an HTTP gateway and a sample client.

---
### 1. Clone the Repository

```bash
git clone https://github.com/poornatejav/gRPC_OrderSystem_in_GO.git
cd  gRPC_OrderSystem_in_GO
```

### 2. Install Proto Requirements 

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### 3. Generating gRPC code from .proto files

```bash
protoc --go_out=. --go-grpc_out=. proto/*.proto
```

### 4. Install Go Dependencies

```bash
go mod tidy
```

### 3. Run Go Services in diferent terminals

```bash
go run customer/main.go
go run inventory/main.go
go run order/main.go
go run gateway/main.go
```

