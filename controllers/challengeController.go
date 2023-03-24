package controllers

import (
	"context"
	"encoding/json"
	"log"
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
)

// ListChallenge -> List all the challenges
var ListChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var challenges []*models.Challenge
	collection := client.Database("challenge").Collection("challenges")
	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	for cursor.Next(context.TODO()) {
		var challenge models.Challenge
		err := cursor.Decode(&challenge)
		if err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
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
	challenge.CreatedAt = time.Now()
	challenge.UpdatedAt = time.Now()
	collection := client.Database("challenge").Collection("challenges")
	var existingChallenge models.Challenge
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "mnemonic", Value: challenge.Mnemonic}}).Decode(&existingChallenge)
	if err == nil {
		middlewares.ErrorResponse("Mnemonic Invalid", rw)
		return
	}
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
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)

	if challenge.Status == "private" {
		if challenge.Coordinator == props["username"] || challenge.RecipientAddress == identity {

			challenge.Participants = append(challenge.Participants, identity)

			res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
				return
			}

			if res.MatchedCount == 0 {
				middlewares.ErrorResponse("challenge does not exist", rw)
				return
			}

			middlewares.SuccessRespond(params["id"], rw)
			return
		}
		middlewares.ForbiddenResponse("you have no access for this challenge", rw)
		return
	}
	challenge.Participants = append(challenge.Participants, identity)

	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("challenge does not exist", rw)
		return
	}

	middlewares.SuccessRespond(params["id"], rw)
})

var UnJoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	var challenge models.Challenge

	collection := client.Database("challenge").Collection("challenges")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)

	var check bool = false
	for i, v := range challenge.Participants {
		if v == identity {
			challenge.Participants = append(challenge.Participants[:i], challenge.Participants[i+1:]...)
			check = true
			break
		}
	}
	if !check {
		middlewares.ErrorResponse("you have already leave this challenge", rw)
		return
	}

	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "Participants", Value: challenge.Participants}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("challenge does not exist", rw)
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
			log.Fatal(err)
		}
		steps = append(steps, &step)
	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	var winnerRecord []*models.Steps

	if challenges.Goal == "distance" {
		floatGoalThershold, err := strconv.ParseFloat(challenges.GoalThreshold, 32)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
		}
		for _, v := range steps {
			distance := strings.Split(v.StepsDistance, " ")
			floatUserDistance, err := strconv.ParseFloat(distance[0], 32)
			if err != nil {
				middlewares.ServerErrResponse(err.Error(), rw)
			}
			if floatUserDistance > floatGoalThershold {
				winnerRecord = append(winnerRecord, v)
			}
		}
	} else if challenges.Goal == "count" {
		intGoalThershold, err := strconv.ParseInt(challenges.GoalThreshold, 10, 64)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
		}
		for _, v := range steps {
			if v.StepsCount > intGoalThershold {
				winnerRecord = append(winnerRecord, v)
			}
		}
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
			log.Fatal(err)
		}
		bets = append(bets, &bet)
	}

	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
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

	middlewares.SuccessRespond(winners, rw)
})
