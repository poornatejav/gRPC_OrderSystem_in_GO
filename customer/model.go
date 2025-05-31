package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Customer struct {
	ID string `bson:"id"`
}

type CustomerModel struct {
	Collection *mongo.Collection
}

func (m *CustomerModel) Exists(ctx context.Context, id string) (bool, error) {
	var result Customer
	err := m.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (m *CustomerModel) Insert(ctx context.Context, id string) error {
	_, err := m.Collection.InsertOne(ctx, Customer{ID: id})
	return err
}
