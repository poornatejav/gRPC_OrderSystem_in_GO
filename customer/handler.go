package main

import (
	"context"
	"strings"

	"grpc-order-system/proto/customerpb"
)

type CustomerService struct {
	customerpb.UnimplementedCustomerServiceServer
	Model *CustomerModel
}

func (s *CustomerService) GetCustomer(ctx context.Context, req *customerpb.GetCustomerRequest) (*customerpb.GetCustomerResponse, error) {
	id := strings.TrimSpace(req.Id)
	exists, err := s.Model.Exists(ctx, id)
	if err != nil {
		return nil, err
	}
	return &customerpb.GetCustomerResponse{Exists: exists}, nil
}
