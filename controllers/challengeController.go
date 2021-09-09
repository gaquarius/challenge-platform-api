package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	middlewares "github.com/gaquarius/challenge-platform-api/handlers"
	"github.com/gaquarius/challenge-platform-api/models"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	layoutISO = "2006-01-02"
)

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
