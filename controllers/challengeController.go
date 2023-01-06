package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// TO-DO USER TO JOIN THE CHALLENGE
var JoinChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	// id, _ := primitive.ObjectIDFromHex(params["id"])
	// var challenge models.Challenge

	// err := json.NewDecoder(r.Body).Decode(&challenge)
	// if err != nil {
	// 	middlewares.ServerErrResponse(err.Error(), rw)
	// 	return
	// }
	// collection := client.Database("challenge").Collection("challenges")
	// res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: challenge}})

})
