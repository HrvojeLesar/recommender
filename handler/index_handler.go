package handler

import (
	"log"
	"net/http"

	"github.com/HrvojeLesar/recommender/db/models"
)

type IndexData struct {
	TopBooks      []models.AverageBookRating
	TopBooksByTag []models.AverageBookRating
	Tags          []models.Tag
	SelectedTag   string
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	data := IndexData{
		Tags: h.Tags,
	}
	topBooks, err := h.globalInstances.Mongo.TopBooksByUserRating()
	if err != nil {
		log.Panicln(err)
	}
	data.TopBooks = topBooks

	data.SelectedTag = r.URL.Query().Get("tag")
	if len(data.SelectedTag) == 0 {
		data.SelectedTag = "classics"
	}
	topBooksByTag, err := h.globalInstances.Mongo.TopBooksByGenre(data.SelectedTag)
	if err != nil {
		log.Panicln(err)
	}
	data.TopBooksByTag = topBooksByTag

	err = h.indexTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
