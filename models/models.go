package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Person Model
type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty" validate:"required,alpha"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty" validate:"required,alpha"`
}

// User Model
type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	Role     string             `json:"role" bson:"role"`
	Bio      string             `json:"bio" bson:"bio"`
	Avatar   string             `json:"avatar" bson:"avatar"`
	Identity string             `json:"identity" bson:"identity"`
	Password string             `json:"password" bson:"password"`
}
