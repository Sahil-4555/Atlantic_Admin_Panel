package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Email     string             `json:"email,omitempty" gorm:"unique"`
	Password  string             `json:"password,omitempty"`
	Createdat string             `json:"createdat,omitempty"`
}
