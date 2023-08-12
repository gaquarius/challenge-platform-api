package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/models"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var AddStepsChallenge = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var steps models.Steps
	err := json.NewDecoder(r.Body).Decode(&steps)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	ChallengeID, _ := primitive.ObjectIDFromHex(steps.ChallengeID)

	var challenge models.Challenge
	collection := client.Database("challenge").Collection("challenges")
	err = collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: ChallengeID}}).Decode(&challenge)
	if err != nil {
		middlewares.ServerErrResponse(fmt.Sprintf("challenge %v not found", steps.ChallengeID), rw)
		return
	}

	stepsCollection := client.Database("challenge").Collection("stepsDetails")

	props, _ := r.Context().Value("props").(jwt.MapClaims)
	identity := props["identity"].(string)
	steps.Identity = identity
	steps.CreatedAt = time.Now().UTC()

	filter := bson.M{
		"challenge_id": steps.ChallengeID,
		"identity":     steps.Identity,
	}
	var existedSteps models.Steps
	err = stepsCollection.FindOne(context.TODO(), filter).Decode(&existedSteps)
	if err != nil && err != mongo.ErrNoDocuments {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	if err == mongo.ErrNoDocuments {
		res, err := stepsCollection.InsertOne(context.TODO(), steps)
		if err != nil {
			middlewares.ServerErrResponse(err.Error(), rw)
			return
		}
		middlewares.SuccessRespond(res, rw)
		return
	}
	if steps.StepsCount != 0 {
		existedSteps.StepsCount = steps.StepsCount
	}
	if len(steps.StepsDistance) > 0 {
		existedSteps.StepsDistance = steps.StepsDistance
	}
	if steps.MinimumStepsCount != 0 {
		existedSteps.MinimumStepsCount = steps.MinimumStepsCount
	}
	if len(steps.MinimumStepsDistance) > 0 {
		existedSteps.MinimumStepsDistance = steps.MinimumStepsDistance
	}
	if len(steps.Distance) > 0 {
		existedSteps.Distance = steps.Distance
	}

	_, err = stepsCollection.ReplaceOne(context.TODO(), bson.M{"_id": existedSteps.ID}, existedSteps)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	middlewares.SuccessRespond(fmt.Sprintf("steps updated at %v", existedSteps.ID.Hex()), rw)
})
