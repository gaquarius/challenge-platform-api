package routes

import (
	"net/http"

	"github.com/chattertechno/challenge-platform-api/controllers"
	middlewares "github.com/chattertechno/challenge-platform-api/handlers"
	"github.com/gorilla/mux"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// User API routes

	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	user.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	user.HandleFunc("/me", middlewares.IsAuthorized(controllers.GetMe)).Methods("GET")
	user.HandleFunc("/me", middlewares.IsAuthorized(controllers.UpdateUser)).Methods("PUT")
	user.HandleFunc("/{username}", controllers.GetUser).Methods("GET")

	// Challenge API routes

	challenge := api.PathPrefix("/challenge").Subrouter()
	challenge.HandleFunc("/", middlewares.IsAuthorized(controllers.ListChallenge)).Methods("GET")
	challenge.HandleFunc("/", middlewares.IsAuthorized(controllers.CreateChallenge)).Methods("POST")
	challenge.HandleFunc("/user/{username}", controllers.GetChallenges).Methods("GET")
	challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.GetChallenge)).Methods("GET")
	challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.UpdateChallenge)).Methods("PUT")
	challenge.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.DeleteChallenge)).Methods("DELETE")
	challenge.HandleFunc("/join/", middlewares.IsAuthorized(controllers.JoinChallenge)).Methods("POST")
	challenge.HandleFunc("/{id}/unjoin/", middlewares.IsAuthorized(controllers.UnJoinChallenge)).Methods("POST")
	challenge.HandleFunc("/{id}/winner/", middlewares.IsAuthorized(controllers.ChallengeWinner)).Methods("GET")

	challenge.HandleFunc("/finished/", controllers.FinishedChallenges).Methods("GET")
	challenge.HandleFunc("/update/flag/{id}", controllers.UpdateFlag).Methods("PUT")

	// Challenge bet routes
	bet := challenge.PathPrefix("/bet").Subrouter()
	bet.HandleFunc("/add/", middlewares.IsAuthorized(controllers.AddBetChallenge)).Methods("POST")
	bet.HandleFunc("/{id}", middlewares.IsAuthorized(controllers.GetAllBetsForChallenge)).Methods("GET")

	// User steps routes
	steps := challenge.PathPrefix("/user/steps").Subrouter()
	steps.HandleFunc("/add/", middlewares.IsAuthorized(controllers.AddStepsChallenge)).Methods("PUT")

	api.HandleFunc("/person", controllers.CreatePersonEndpoint).Methods("POST")
	api.HandleFunc("/people", middlewares.IsAuthorized(controllers.GetPeopleEndpoint)).Methods("GET")
	api.HandleFunc("/person/{id}", controllers.GetPersonEndpoint).Methods("GET")
	api.HandleFunc("/person/{id}", controllers.DeletePersonEndpoint).Methods("DELETE")
	api.HandleFunc("/person/{id}", controllers.UpdatePersonEndpoint).Methods("PUT")

	router.HandleFunc("/upload", controllers.UploadFileEndpoint).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./uploaded/"))))
	return router
}
