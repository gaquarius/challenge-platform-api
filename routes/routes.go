package routes

import (
	"net/http"

	"github.com/gaquarius/challenge-platform-api/controllers"
	middlewares "github.com/gaquarius/challenge-platform-api/handlers"
	"github.com/gorilla/mux"
)

// Routes -> define endpoints
func Routes() *mux.Router {
	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	user := api.PathPrefix("/user").Subrouter()
	user.HandleFunc("/register", controllers.RegisterUser).Methods("POST")
	user.HandleFunc("/login", controllers.LoginUser).Methods("POST")
	api.HandleFunc("/person", controllers.CreatePersonEndpoint).Methods("POST")
	api.HandleFunc("/people", middlewares.IsAuthorized(controllers.GetPeopleEndpoint)).Methods("GET")
	api.HandleFunc("/person/{id}", controllers.GetPersonEndpoint).Methods("GET")
	api.HandleFunc("/person/{id}", controllers.DeletePersonEndpoint).Methods("DELETE")
	api.HandleFunc("/person/{id}", controllers.UpdatePersonEndpoint).Methods("PUT")
	router.HandleFunc("/upload", controllers.UploadFileEndpoint).Methods("POST")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./uploaded/"))))
	return router
}
