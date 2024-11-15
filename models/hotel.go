package models

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Hotel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name" validate:"required,min=1,max=50"`
	Address    string             `bson:"address" validate:"required"`
	District   string             `bson:"district" validate:"required"`
	Province   string             `bson:"province" validate:"required"`
	PostalCode string             `bson:"postalcode" validate:"required,len=5"`
	Tel        string             `bson:"tel,omitempty"`
	Region     string             `bson:"region" validate:"required"`
}

// PreDeleteHook performs cascading deletion of related appointments when a hotel is deleted.
func (hotel *Hotel) PreDeleteHook(ctx context.Context, db *mongo.Database) error {
	appointmentsCollection := db.Collection("appointments")
	_, err := appointmentsCollection.DeleteMany(ctx, bson.M{"hotel": hotel.ID})
	if err != nil {
		return fmt.Errorf("failed to delete appointments for hotel %s: %w", hotel.ID.Hex(), err)
	}
	fmt.Printf("Appointments removed for hotel %s\n", hotel.ID.Hex())
	return nil
}
