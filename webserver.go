package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/HrvojeLesar/recommender/config"
	"github.com/HrvojeLesar/recommender/handler"
	"github.com/gorilla/mux"
)

func StartWebserver(handler *handler.Handler) {
	r := mux.NewRouter()
	r.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {}).Methods("GET")
	r.HandleFunc("/", handler.Index).Methods("GET")
	r.HandleFunc("/{userId}", handler.PersonalizedIndex).Methods("GET")
	r.HandleFunc("/{userId}", handler.MyRatingsUpdate).Methods("POST")
	r.HandleFunc("/{userId}/my-ratings", handler.MyRatings).Methods("GET")
	r.HandleFunc("/{userId}/my-ratings", handler.MyRatingsUpdate).Methods("POST")

	http.Handle("/", r)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	port := config.LookupEnvVariableOrDefault("PORT", "8000")
	fmt.Printf("Listening on http://localhost:%s\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Panic(err)
	}
}
