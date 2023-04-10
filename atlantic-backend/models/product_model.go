package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id          primitive.ObjectID `json:"id,omitempty"`
	Productid   string             `json:"productid,omitempty"`
	Title       string             `json:"title,omitempty"`
	Price       uint               `json:"price,omitempty"`
	Size        string             `json:"size,omitempty"`
	Description string             `json:"description,omitempty"`
	Color       string             `json:"color,omitempty"`
	Createdat   string             `json:"createdat,omitempty"`
	Updatedat   string             `json:"updatedat,omitempty"`
	Image       string             `json:"image,omitempty"`
}
