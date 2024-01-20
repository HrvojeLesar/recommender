package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type SearchData struct {
	SearchTerm string
	Books      []bookRecommendation
}

func (h *Handler) Search(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIdStr := vars["userId"]
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	searchTerm := r.URL.Query().Get("query")

	usersRatings, err := h.globalInstances.Mongo.MyRatings(userId)
	if err != nil {
		log.Panicln(err)
	}

	books, err := h.globalInstances.Mongo.Search(searchTerm)
	if err != nil {
		log.Panicln(err)
	}

	bookRecommendations, _ := integrateUserRatings(usersRatings, books, IndexData{})

	data := SearchData{
		Books:      bookRecommendations,
		SearchTerm: searchTerm,
	}

	err = h.searchTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
