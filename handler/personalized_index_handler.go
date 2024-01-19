package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/HrvojeLesar/recommender/db/models"
	"github.com/gorilla/mux"
)

type PersonalizedIndexData struct {
	BookRecommendations []models.Book
}

func (h *Handler) PersonalizedIndex(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	books, err := h.globalInstances.Mongo.NearestNeighbour(userId)
	if err != nil {
		log.Panicln(err)
	}

	data := PersonalizedIndexData{
		BookRecommendations: books,
	}

	err = h.personalizedIndexTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
