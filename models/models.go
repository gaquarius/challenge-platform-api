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
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Role     string             `json:"role,omitempty" bson:"role,omitempty"`
	Bio      string             `json:"bio,omitempty" bson:"bio,omitempty"`
	Avatar   string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Identity string             `json:"identity,omitempty" bson:"identity,omitempty"`
	Mnemonic string             `json:"mnemonic,omitempty" bson:"mnemonic,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
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
	ID              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartDate       string             `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate         string             `json:"end_date,omitempty" bson:"end_date,omitempty"`
	Status          ChallengeStatus    `json:"status,omitempty" bson:"status,omitempty"`
	Goal            string             `json:"goal,omitempty" bson:"goal,omitempty"`
	GoalIncreaments string             `json:"goal_increaments,omitempty" bson:"goal_increaments,omitempty"`
	GoalThreshold   string             `json:"goal_threshold,omitempty" bson:"goal_threshold,omitempty"`
	Category        []string           `json:"category,omitempty" bson:"category,omitempty"`
	Participants    []string           `json:"participants,omitempty" bson:"participants,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Description     string             `json:"description,omitempty" bson:"description,omitempty"`
	Content         string             `json:"content,omitempty" bson:"content,omitempty"`
	Identity        string             `json:"identity,omitempty" bson:"identity,omitempty"`
	Mnemonic        string             `json:"mnemonic,omitempty" bson:"mnemonic,omitempty"`
	HeaderImage     string             `json:"header_image,omitempty" bson:"header_image,omitempty"`
	Coordinator     string             `json:"coordinator,omitempty" bson:"coordinator,omitempty"`
	CreatedAt       time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt       time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type EscrowStatus string

const (
	Pending EscrowStatus = "pending"
	Paid    EscrowStatus = "paid"
)

type Escrow struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Amount    float32            `json:"amount,omitempty" bson:"amount,omitempty"`
	Challenge primitive.ObjectID `json:"challenge,omitempty" bson:"challenge,omitempty"`
	Status    EscrowStatus       `json:"status,omitempty" bson:"status,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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
	Challenge   primitive.ObjectID `json:"challenge,omitempty" bson:"challenge,omitempty"`
	Participant string             `json:"participant,omitempty" bson:"participant,omitempty"`
	Type        ActivityType       `json:"type,omitempty" bson:"type,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
