package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" binding:"required" bson:"name,omitempty"`
	Price       int                `json:"price" binding:"required" bson:"price,omitempty"`
	Description string             `json:"description" binding:"required" bson:"description,omitempty"`
}
