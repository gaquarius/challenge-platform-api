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
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username,omitempty" bson:"username,omitempty"`
	Role       string             `json:"role,omitempty" bson:"role,omitempty"`
	Bio        string             `json:"bio,omitempty" bson:"bio,omitempty"`
	Avatar     string             `json:"avatar,omitempty" bson:"avatar,omitempty"`
	Identity   string             `json:"identity,omitempty" bson:"identity,omitempty"`
	Mnemonic   string             `json:"mnemonic,omitempty" bson:"mnemonic,omitempty"`
	PrivateKey string             `json:"private_key,omitempty" bson:"private_key,omitempty"`
	Password   string             `json:"password,omitempty" bson:"password,omitempty"`
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
	AddBet          string             `json:"add_bet,omitempty" bson:"add_bet,omitempty"`
	Category        []string           `json:"category,omitempty" bson:"category,omitempty"`
	Name            string             `json:"name,omitempty" bson:"name,omitempty"`
	Description     string             `json:"description,omitempty" bson:"description,omitempty"`
	// Mnemonic         string             `json:"mnemonic,omitempty" bson:"mnemonic,omitempty"`
	FundDeliveredFlag bool            `json:"fund_delivered_flag" bson:"fund_delivered_flag"`
	Content           string          `json:"content,omitempty" bson:"content,omitempty"`
	HeaderImage       string          `json:"header_image,omitempty" bson:"header_image,omitempty"`
	Coordinator       string          `json:"coordinator,omitempty" bson:"coordinator,omitempty"`
	Identity          string          `json:"identity,omitempty" bson:"identity,omitempty"`
	Visible           bool            `json:"visible,omitempty" bson:"visible,omitempty"`
	RecipientAddress  string          `json:"recipient_address,omitempty" bson:"recipient_address,omitempty"`
	MinBetAmount      float64         `json:"min_bet_amount,omitempty" bson:"min_bet_amount,omitempty"`
	Participants      []JoinChallenge `json:"participant,omitempty" bson:"participant,omitempty"`
	CreatedAt         time.Time       `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt         time.Time       `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
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

type Bet struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Identity    string             `json:"identity,omitempty" bson:"identity,omitempty"`
	Amount      float64            `json:"amount,omitempty" bson:"amount,omitempty"`
	ChallengeID string             `json:"challenge_id,omitempty" bson:"challenge_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type JoinChallenge struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Identity    string             `json:"identity,omitempty" bson:"identity,omitempty"`
	Bet         float64            `json:"bet,omitempty" bson:"bet,omitempty"`
	ChallengeID string             `json:"challenge_id,omitempty" bson:"challenge_id,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type Steps struct {
	ID                   primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Identity             string             `json:"identity,omitempty" bson:"identity,omitempty"`
	StepsCount           int64              `json:"steps_count,omitempty" bson:"steps_count,omitempty"`
	StepsDistance        string             `json:"steps_distance,omitempty" bson:"steps_distance,omitempty"`
	Distance             string             `json:"distance,omitempty" bson:"distance,omitempty"`
	MinimumStepsCount    int64              `json:"minimum_steps_count,omitempty" bson:"minimum_steps_count,omitempty"`
	MinimumStepsDistance string             `json:"minimum_steps_distance,omitempty" bson:"minimum_steps_distance,omitempty"`
	ChallengeID          string             `json:"challenge_id,omitempty" bson:"challenge_id,omitempty"`
	CreatedAt            time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type WinnerResponse struct {
	Identity    string  `json:"identity,omitempty" bson:"identity,omitempty"`
	ChallengeID string  `json:"challenge_id,omitempty" bson:"challenge_id,omitempty"`
	Amount      float64 `json:"amount,omitempty" bson:"amount,omitempty"`
}

type UpdateFlagRequest struct {
	FundDeliveredFlag bool `json:"fund_delivered_flag" bson:"fund_delivered_flag"`
}

type GetChallenges struct {
	ID                primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	StartDate         string             `json:"start_date,omitempty" bson:"start_date,omitempty"`
	EndDate           string             `json:"end_date,omitempty" bson:"end_date,omitempty"`
	FundDeliveredFlag bool               `json:"fund_delivered_flag" bson:"fund_delivered_flag"`
	Coordinator       string             `json:"coordinator,omitempty" bson:"coordinator,omitempty"`
	Identity          string             `json:"identity,omitempty" bson:"identity,omitempty"`
}
