package controllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chattertechno/challenge-platform-api/db"
	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/chattertechno/challenge-platform-api/models"
	"github.com/chattertechno/challenge-platform-api/validators"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var client = db.Dbconnect()

// RegisterUser -> Register User with Menmonic and username
var RegisterUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	collection := client.Database("challenge").Collection("users")
	var existingUser models.User
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("Username is already taken.", rw)
		return
	}
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "identity", Value: user.Identity}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("Identity is already in use.", rw)
		return
	}
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "mnemonic", Value: user.Mnemonic}}).Decode(&existingUser)
	if err == nil {
		middlewares.ErrorResponse("Mnemonic Invalid", rw)
		return
	}
	passwordHash, err := middlewares.HashPassword(user.Password)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	user.Password = passwordHash
	result, err := collection.InsertOne(r.Context(), user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`Inserted at `+strings.Replace(string(res), `"`, ``, 2), rw)
})

// LoginUser -> Let the user login with identity and password
var LoginUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}
	collection := client.Database("challenge").Collection("users")
	var existingUser models.User
	err = collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: user.Username}}).Decode(&existingUser)

	if err != nil {
		middlewares.ErrorResponse("User doesn't exist", rw)
		return
	}
	isPasswordMatch := middlewares.CheckPasswordHash(user.Password, existingUser.Password)
	if !isPasswordMatch {
		middlewares.ErrorResponse("Password doesn't match", rw)
		return
	}
	token, err := middlewares.GenerateJWT(user.Username)
	if err != nil {
		middlewares.ErrorResponse("Failed to generate JWT", rw)
		return
	}
	middlewares.SuccessResponse(string(token), rw)
})

// GetMe -> Get user details from Authorization token
var GetMe = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var user models.User
	collection := client.Database("challenge").Collection("users")
	err := collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: props["username"]}}).Decode(&user)
	if err != nil {
		middlewares.AuthorizationResponse("Malformed token", rw)
		return
	}

	user.Password = ""
	middlewares.SuccessRespond(user, rw)
})

// GetUser -> Get user details from username
var GetUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user models.User

	collection := client.Database("challenge").Collection("users")
	err := collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: params["username"]}}).Decode(&user)
	if err != nil {
		middlewares.ErrorResponse("User doesn't exist", rw)
		return
	}

	user.Password = ""
	middlewares.SuccessRespond(user, rw)
})

// UpdateUser -> Update user details from username
var UpdateUser = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
	props, _ := r.Context().Value("props").(jwt.MapClaims)
	var user models.User

	collection := client.Database("challenge").Collection("users")
	err := collection.FindOne(r.Context(), bson.D{primitive.E{Key: "username", Value: props["username"]}}).Decode(&user)
	if err != nil {
		middlewares.AuthorizationResponse("Malformed token", rw)
		return
	}

	var newUser models.User
	err = json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), rw)
		return
	}

	res, err := collection.UpdateOne(r.Context(), bson.D{primitive.E{Key: "_id", Value: user.ID}}, bson.D{
		primitive.E{
			Key: "$set",
			Value: bson.D{
				primitive.E{Key: "username", Value: newUser.Username},
				primitive.E{Key: "bio", Value: newUser.Bio},
			},
		},
	})

	if err != nil {
		middlewares.ErrorResponse("Username is already taken.", rw)
		return
	}
	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("User doesn't exist", rw)
		return
	}

	token, err := middlewares.GenerateJWT(newUser.Username)
	if err != nil {
		middlewares.ErrorResponse("Failed to generate JWT", rw)
		return
	}
	middlewares.SuccessResponse(string(token), rw)
})

// CreatePersonEndpoint -> create person
var CreatePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	var person models.Person
	err := json.NewDecoder(request.Body).Decode(&person)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	if ok, errors := validators.ValidateInputs(person); !ok {
		middlewares.ValidationResponse(errors, response)
		return
	}
	collection := client.Database("golang").Collection("people")
	result, err := collection.InsertOne(context.TODO(), person)
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	res, _ := json.Marshal(result.InsertedID)
	middlewares.SuccessResponse(`Inserted at `+strings.Replace(string(res), `"`, ``, 2), response)
})

// GetPeopleEndpoint -> get people
var GetPeopleEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	var people []*models.Person

	collection := client.Database("golang").Collection("people")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	for cursor.Next(context.TODO()) {
		var person models.Person
		err := cursor.Decode(&person)
		if err != nil {
			log.Fatal(err)
		}

		people = append(people, &person)
	}
	if err := cursor.Err(); err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	middlewares.SuccessArrRespond(people, response)
})

// GetPersonEndpoint -> get person by id
var GetPersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person models.Person

	collection := client.Database("golang").Collection("people")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&person)
	if err != nil {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	middlewares.SuccessRespond(person, response)
})

// DeletePersonEndpoint -> delete person by id
var DeletePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var person models.Person

	collection := client.Database("golang").Collection("people")
	err := collection.FindOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&person)
	if err != nil {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	_, derr := collection.DeleteOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}})
	if derr != nil {
		middlewares.ServerErrResponse(derr.Error(), response)
		return
	}
	middlewares.SuccessResponse("Deleted", response)
})

// UpdatePersonEndpoint -> update person by id
var UpdatePersonEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	type fname struct {
		Firstname string `json:"firstname"`
	}
	var fir fname
	json.NewDecoder(request.Body).Decode(&fir)
	collection := client.Database("golang").Collection("people")
	res, err := collection.UpdateOne(context.TODO(), bson.D{primitive.E{Key: "_id", Value: id}}, bson.D{primitive.E{Key: "$set", Value: bson.D{primitive.E{Key: "firstname", Value: fir.Firstname}}}})
	if err != nil {
		middlewares.ServerErrResponse(err.Error(), response)
		return
	}
	if res.MatchedCount == 0 {
		middlewares.ErrorResponse("Person does not exist", response)
		return
	}
	middlewares.SuccessResponse("Updated", response)
})

// UploadFileEndpoint -> upload file
var UploadFileEndpoint = http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
	file, handler, err := request.FormFile("file")
	// fileName := request.FormValue("file_name")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f, err := os.OpenFile("uploaded/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = io.Copy(f, file)

	middlewares.SuccessResponse("Uploaded Successfully", response)
})
