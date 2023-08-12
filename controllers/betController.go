package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var AddBetChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var bet models.Bet
	err := json.NewDecoder(r.Body).Decode(&bet)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	ChallengeID, _ := primitive.ObjectIDFromHex(bet.ChallengeID)

	var challenge models.Challenge
	collection := client.Database("challenge").Collection("challenges")
	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: ChallengeID}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(fmt.Sprintf("challenge %v not found", bet.ChallengeID), rw)
		return
	}

	if challenge.MinBetAmount > bet.Amount {
		middlewares.ErrorResponse(fmt.Sprintf("bet must be greater than %v", challenge.MinBetAmount), rw)
		return
	}

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)
	bet.Identity = identity
	bet.CreatedAt = time.Now().UTC()

	betCollection := client.Database("challenge").Collection("bets")
	result, err := betCollection.InsertOne(context.TODO(), bet)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	middlewares.SuccessRespond(result, rw)
})

var GetAllBetsForChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var bets []*models.Bet

	betCollection := client.Database("challenge").Collection("bets")

	cursor, err := betCollection.Find(context.TODO(), bson.D{primitive.E{Key: "challenge_id", Value: params["id"]}})
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
	middlewares.SuccessRespond(bets, rw)
})
