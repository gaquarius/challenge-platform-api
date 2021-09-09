package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type ChallengeStatus string

const (
	Open           ChallengeStatus = "open"
	InvitationOnly ChallengeStatus = "invites only"
	Private        ChallengeStatus = "private"
	Draft          ChallengeStatus = "draft"
)

// Challenge Model
type Challenge struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartDate   string             `json:"start_date" bson:"start_date"`
	EndDate     string             `json:"end_date" bson:"end_date"`
	Status      ChallengeStatus    `json:"status" bson:"status"`
	Category    []string           `json:"category" bson:"category"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Content     string             `json:"content" bson:"content"`
	HeaderImage string             `json:"header_image" bson:"header_image"`
	Coordinator string             `json:"coordinator" bson:"coordinator"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}

type EscrowStatus string

const (
	Pending EscrowStatus = "pending"
	Paid    EscrowStatus = "paid"
)

type Escrow struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Amount    float32            `json:"amount" bson:"amount"`
	Challenge primitive.ObjectID `json:"challenge" bson:"challenge"`
	Status    EscrowStatus       `json:"status" bson:"status"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type ActivityType string

const (
	Joined ActivityType = "joined"
	Won    ActivityType = "won"
	Lost   ActivityType = "lost"
	Played ActivityType = "played"
)

type Activity struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Challenge   primitive.ObjectID `json:"challenge" bson:"challenge"`
	Participant string             `json:"participant" bson:"participant"`
	Type        ActivityType       `json:"type" bson:"type"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
}
