package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type InventoryItem struct {
	ItemID   string `bson:"item_id"`
	Quantity int32  `bson:"quantity"`
}

type InventoryModel struct {
	Collection *mongo.Collection
}

func (m *InventoryModel) IsAvailable(ctx context.Context, itemID string, quantity int32) (bool, error) {
	var item InventoryItem
	err := m.Collection.FindOne(ctx, bson.M{"item_id": itemID}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return item.Quantity >= quantity, nil
}

func (m *InventoryModel) Deduct(ctx context.Context, itemID string, quantity int32) (bool, error) {
	res, err := m.Collection.UpdateOne(ctx,
		bson.M{"item_id": itemID, "quantity": bson.M{"$gte": quantity}},
		bson.M{"$inc": bson.M{"quantity": -quantity}},
	)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (m *InventoryModel) Insert(ctx context.Context, itemID string, quantity int32) error {
	_, err := m.Collection.InsertOne(ctx, InventoryItem{ItemID: itemID, Quantity: quantity})
	return err
}
