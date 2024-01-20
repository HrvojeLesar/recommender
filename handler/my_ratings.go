package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/HrvojeLesar/recommender/db/models"
	"github.com/gorilla/mux"
)

type MyRatings struct {
	User         *models.User
	SimilarUsers []models.UserSimilarity
}

func (h *Handler) MyRatings(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	userRatings, err := h.globalInstances.Mongo.MyRatings(userId)
	if err != nil {
		log.Panicln(err)
	}

	similarUsers := make([]models.UserSimilarity, 0, 0)
	similarUsersQueryParam := r.URL.Query().Get("similarUsers")
	if len(similarUsersQueryParam) > 0 {
		users, err := h.globalInstances.Mongo.SimilarUsers(userId)
		if err != nil {
			log.Panicln(err)
		}
		similarUsers = users.SimilarUsers()
	}

	err = h.myRatingsTemplate.Execute(w, MyRatings{User: userRatings, SimilarUsers: similarUsers})
	if err != nil {
		log.Panicln(err)
	}
}

func (h *Handler) MyRatingsUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rating, err := strconv.ParseInt(r.FormValue("rating"), 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	userId, err := strconv.ParseInt(vars["userId"], 10, 64)
	if err != nil {
		log.Panicln(err)
	}
	bookId, err := strconv.ParseInt(r.FormValue("bookId"), 10, 64)
	if err != nil {
		log.Panicln(err)
	}

	insertion := r.FormValue("new")

	if len(insertion) != 0 {
		book, err := h.globalInstances.Mongo.Book(bookId)
		if err != nil {
			log.Panicln(err)
		}
		err = h.globalInstances.Mongo.InsertRating(rating, userId, *book)
		if err != nil {
			log.Panicln(err)
		}
	} else if rating > 0 && rating <= 5 {
		err := h.globalInstances.Mongo.UpdateRating(rating, userId, bookId)
		if err != nil {
			log.Panicln(err)
		}
	} else if rating == 0 {
		h.globalInstances.Mongo.RemoveRating(rating, userId, bookId)
		if err != nil {
			log.Panicln(err)
		}
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)
}
