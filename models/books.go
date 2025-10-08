package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	//ID     int    `json:"id"`
	ID     primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title  string             `json:"title"`
	Author string             `json:"author"`
	Year   int                `json:"year"`
	Status string             `json:"status"`
}

type StatusRequest struct {
	Action string `json:"action"`
}
