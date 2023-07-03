package main

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// form: "title" to map the JSON field name to the struct
// binding: "required" to enforce the value is required

type BookCreateUpdateRequest struct {
	Title  string `form:"title" binding:"required"`
	Author string `form:"author"`
}

// json: "id" to map the struct Name to its Json field name

type BookResponse struct {
	Id        primitive.ObjectID `json:"id"`
	Title     string             `json:"title"`
	Author    string             `json:"author"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
