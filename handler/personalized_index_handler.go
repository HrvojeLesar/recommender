package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/HrvojeLesar/recommender/db/models"
	"github.com/gorilla/mux"
)

type PersonalizedIndexData struct {
	BookRecommendations []bookRecommendation
	IndexData           indexDataWithUserRatings
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

	usersRatings, err := h.globalInstances.Mongo.MyRatings(userId)
	if err != nil {
		log.Panicln(err)
	}

	bookRecommendations, indexData := integrateUserRatings(usersRatings, books, h.resolveIndexData(r))

	data := PersonalizedIndexData{
		BookRecommendations: bookRecommendations,
		IndexData:           indexData,
	}

	err = h.personalizedIndexTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
}

type bookRecommendation struct {
	Book       models.Book
	UserRating int
}

type topBook struct {
	Book       models.AverageBookRating
	UserRating int
}

type indexDataWithUserRatings struct {
	TopBooks      []topBook
	TopBooksByTag []topBook
	Tags          []models.Tag
	SelectedTag   string
}

func integrateUserRatings(usersRatings *models.User, books []models.Book, indexData IndexData) ([]bookRecommendation, indexDataWithUserRatings) {
	ratingsMap := make(map[string]models.BookRating)
	for _, userRating := range usersRatings.BookRatings {
		ratingsMap[userRating.Book.Title] = userRating
	}

	recommendations := make([]bookRecommendation, 0, len(books))
	for _, book := range books {
		rating := 0
		if userRating, ok := ratingsMap[book.Title]; ok {
			rating = userRating.Rating
		}
		recommendations = append(recommendations, bookRecommendation{
			Book:       book,
			UserRating: rating,
		})
	}

	topBooks := make([]topBook, 0, len(indexData.TopBooks))
	for _, book := range indexData.TopBooks {
		rating := 0
		if userRating, ok := ratingsMap[book.Book.Title]; ok {
			rating = userRating.Rating
		}
		topBooks = append(topBooks, topBook{
			Book:       book,
			UserRating: rating,
		})
	}

	topBooksByTag := make([]topBook, 0, len(indexData.TopBooksByTag))
	for _, book := range indexData.TopBooksByTag {
		rating := 0
		if userRating, ok := ratingsMap[book.Book.Title]; ok {
			rating = userRating.Rating
		}
		topBooksByTag = append(topBooksByTag, topBook{
			Book:       book,
			UserRating: rating,
		})
	}

	return recommendations, indexDataWithUserRatings{
		TopBooks:      topBooks,
		TopBooksByTag: topBooksByTag,
		Tags:          indexData.Tags,
		SelectedTag:   indexData.SelectedTag,
	}
}
