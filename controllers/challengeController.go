package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ListChallenge -> List all the challenges
var ListChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var challenges []*models.Challenge
	opts := options.Find().SetSort(bson.D{primitive.E{Key: "created_at", Value: -1}})

	collection := client.Database("challenge").Collection("challenges")
	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	for cursor.Next(context.TODO()) {
		var challenge models.Challenge
		err := cursor.Decode(&challenge)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		challenges = append(challenges, &challenge)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessChallengeArrRespond(challenges, rw)
})

// GetChallengs -> Get challenges for specific user
var GetChallenges = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var challenges []*models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	cursor, err := collection.Find(context.TODO(), bson.D{primitive.E{Key: "coordinator", Value: params["username"]}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	for cursor.Next(context.TODO()) {
		var challenge models.Challenge
		err := cursor.Decode(&challenge)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		challenges = append(challenges, &challenge)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessChallengeArrRespond(challenges, rw)
})

// CreateChallenge -> Create a challenge
var CreateChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var challenge models.Challenge
	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	challenge.Identity = props["identity"].(string)

	now := time.Now().UTC()
	challenge.FundDeliveredFlag = false
	challenge.CreatedAt = now
	challenge.UpdatedAt = now
	collection := client.Database("challenge").Collection("challenges")
	result, err := collection.InsertOne(context.TODO(), challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`Inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// GetChallenge -> Get a challenge by id
var GetChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	challengeJoined := client.Database("challenge").Collection("challengeJoined")
	cursor, err := challengeJoined.Find(context.TODO(), bson.D{primitive.E{Key: "challenge_id", Value: params["id"]}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var participants []*models.JoinChallenge

	for cursor.Next(context.TODO()) {
		var participant models.JoinChallenge
		err := cursor.Decode(&participant)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		participants = append(participants, &participant)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	challenge.Participants = participants

	middlewares.SuccessRespond(challenge, rw)
})

var UpdateChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var challenge models.Challenge

	err := json.NewDecoder(r.Body).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	collection := client.Database("challenge").Collection("challenges")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: challenge}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("Challenge does not exist", rw)
		return
	}

	middlewares.SuccessResponse("Updated", rw)
})

var DeleteChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

})

var JoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var req models.JoinChallenge
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	ChallengeID, _ := primitive.ObjectIDFromHex(req.ChallengeID)

	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: ChallengeID}}).Decode(&challenge)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ErrorResponse("challenge does not exists", rw)
			return
		}
		middlewares.ErrorResponse(err.Error(), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)

	challengeBet, err := strconv.ParseFloat(challenge.AddBet, 64)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if challengeBet > req.Bet {
		middlewares.ErrorResponse(fmt.Sprintf("bet must be greater than %v", challengeBet), rw)
		return
	}

	req.Identity = identity
	req.CreatedAt = time.Now().UTC()

	challengeJoined := client.Database("challenge").Collection("challengeJoined")

	filter := bson.M{
		"challenge_id": req.ChallengeID,
		"identity":     req.Identity,
	}

	var getjoinedChallenge models.JoinChallenge

	err = challengeJoined.FindOne(context.TODO(), filter).Decode(&getjoinedChallenge)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if getjoinedChallenge.ChallengeID != "" || getjoinedChallenge.Identity != "" {
		if err != nil {
			middlewares.ErrorResponse("you already joined this challenge", rw)
			return
		}
	}

	result, err := challengeJoined.InsertOne(context.TODO(), req)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var bet models.Bet
	bet.Identity = identity
	bet.ChallengeID = req.ChallengeID
	bet.Amount = req.Bet

	betCollection := client.Database("challenge").Collection("bets")
	_, err = betCollection.InsertOne(context.TODO(), bet)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`Inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

var UnJoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)

	filter := bson.M{
		"challenge_id": params["id"],
		"identity":     identity,
	}

	challengeJoined := client.Database("challenge").Collection("challengeJoined")

	deleteResult, err := challengeJoined.DeleteOne(context.TODO(), filter)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if deleteResult.DeletedCount == 0 {
		middlewares.ErrorResponse("challenge does not exists", rw)
		return
	}

	middlewares.SuccessResponse("unjoin challenge successfully", rw)
})

// ChallengeWinner -> Get all the winners of challenge
var ChallengeWinner = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var challenges *models.Challenge

	id, _ := primitive.ObjectIDFromHex(params["id"])

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenges)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ServerErrResponse("challenge not found", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var steps []*models.Steps

	stepsCollection := client.Database("challenge").Collection("stepsDetails")
	cursor, err := stepsCollection.Find(context.TODO(), bson.D{primitive.E{Key: "challenge_id", Value: params["id"]}})
	if err != nil {
		if err == mongo.ErrNoDocuments {
			middlewares.ServerErrResponse("users steps record not found for this challenge", rw)
			return
		}
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var step models.Steps
		err := cursor.Decode(&step)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		steps = append(steps, &step)
	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var bets []*models.Bet

	betCollection := client.Database("challenge").Collection("bets")

	cursor, err = betCollection.Find(context.TODO(), bson.D{primitive.E{Key: "challenge_id", Value: params["id"]}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	for cursor.Next(context.TODO()) {
		var bet models.Bet
		err := cursor.Decode(&bet)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		bets = append(bets, &bet)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	checkChallengeWins := false

	goalThreshold, err := strconv.ParseFloat(challenges.GoalThreshold, 64)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	for _, value := range steps {
		stepsDistance, err := strconv.ParseFloat(value.StepsDistance, 64)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		if goalThreshold < stepsDistance {
			checkChallengeWins = true
			break
		}

	}

	var winnerRecord []*models.Steps

	if checkChallengeWins {

		if challenges.Goal == "distance" {
			floatGoalThershold, err := strconv.ParseFloat(challenges.GoalThreshold, 32)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			for _, v := range steps {
				floatUserDistance, err := strconv.ParseFloat(v.StepsDistance, 64)
				if err != nil {
					middlewares.ServerErrResponse(err.Error(), rw)
					return
				}
				if floatUserDistance > floatGoalThershold {
					winnerRecord = append(winnerRecord, v)
				}
			}
		} else if challenges.Goal == "count" {
			intGoalThershold, err := strconv.ParseInt(challenges.GoalThreshold, 10, 64)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}
			for _, v := range steps {
				if v.StepsCount > intGoalThershold {
					winnerRecord = append(winnerRecord, v)
				}
			}
		}

		var totalAmount float64

		for _, bet := range bets {
			totalAmount = totalAmount + bet.Amount
		}
		var count int64
		var winners []models.WinnerResponse
		for _, winner := range winnerRecord {
			for _, user := range bets {
				var win models.WinnerResponse
				if user.Identity == winner.Identity {
					win.ChallengeID = user.ChallengeID
					win.Identity = user.Identity
					// win.Amount = averageAmount
					count++
					winners = append(winners, win)
				}
			}
		}
		averageAmount := totalAmount / float64(count)
		for i := range winners {
			winners[i].Amount = averageAmount
		}
		middlewares.SuccessRespondWithCustomMessage(winners, "challenge winners", rw)
		return
	}
	middlewares.SuccessRespondWithCustomMessage(bets, "no winner", rw)
})

var UpdateFlag = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var challenge models.UpdateFlagRequest

	id, _ := primitive.ObjectIDFromHex(params["id"])
	challenge.FundDeliveredFlag = true

	collection := client.Database("challenge").Collection("challenges")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: challenge}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("Challenge does not exist", rw)
		return
	}

	middlewares.SuccessResponse("Flag updated", rw)
})

var FinishedChallenges = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var challenges []*models.GetChallenges

	collection := client.Database("challenge").Collection("challenges")
	ctx := context.TODO()
	currentDateTime := time.Now().Format("2006-01-02")

	filter := bson.M{
		"end_date":            bson.M{"$gte": currentDateTime},
		"fund_delivered_flag": false,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var challenge models.GetChallenges
		err := cursor.Decode(&challenge)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}

		challenges = append(challenges, &challenge)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	middlewares.SuccessRespond(challenges, rw)
})
