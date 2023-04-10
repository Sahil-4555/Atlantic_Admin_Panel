package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AppUsers struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Uid       string             `json:"uid,omitempty"`
	Email     string             `json:"email,omitempty"`
	Photourl  string             `json:"photourl,omitempty"`
	Name      string             `json:"name,omitempty"`
	Createdat string             `json:"createdat,omitempty"`
	Updatedat string             `json:"updatedat,omitempty"`
}
