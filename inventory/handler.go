package main

import (
	"context"

	"grpc-order-system/proto/inventorypb"
)

type InventoryService struct {
	inventorypb.UnimplementedInventoryServiceServer
	Model *InventoryModel
}

func (s *InventoryService) CheckItemAvailable(ctx context.Context, req *inventorypb.CheckItemRequest) (*inventorypb.CheckItemResponse, error) {
	available, err := s.Model.IsAvailable(ctx, req.ItemId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &inventorypb.CheckItemResponse{Available: available}, nil
}

func (s *InventoryService) DeductItem(ctx context.Context, req *inventorypb.DeductItemRequest) (*inventorypb.DeductItemResponse, error) {
	success, err := s.Model.Deduct(ctx, req.ItemId, req.Quantity)
	if err != nil {
		return nil, err
	}
	return &inventorypb.DeductItemResponse{Success: success}, nil
}
