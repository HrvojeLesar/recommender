package handler

import (
	"log"
	"text/template"

	"github.com/HrvojeLesar/recommender/db/models"
	"github.com/HrvojeLesar/recommender/global"
)

type Handler struct {
	globalInstances           *global.Instance
	indexTemplate             *template.Template
	personalizedIndexTemplate *template.Template
	myRatingsTemplate         *template.Template
	Tags                      []models.Tag
}

func NewWebserverHandler(globalInstances *global.Instance) Handler {
	indexTemplate := template.Must(template.ParseFiles("templates/index.html", "templates/header.html", "templates/book_bar_average_rating.html"))
	personalizedIndexTemplate := template.Must(template.ParseFiles("templates/index_perso.html", "templates/header.html", "templates/book_bar.html", "templates/book_bar_average_rating.html", "templates/book_bar_index_average_rating.html"))
	myRatingsTemplate := template.Must(template.ParseFiles("templates/my_ratings.html", "templates/header.html", "templates/book_bar.html", "templates/book_bar_average_rating.html", "templates/book_bar_my_ratings.html"))
	tags, err := globalInstances.Mongo.Tags()
	if err != nil {
		log.Panicln(err)
	}
	return Handler{
		globalInstances:           globalInstances,
		indexTemplate:             indexTemplate,
		personalizedIndexTemplate: personalizedIndexTemplate,
		myRatingsTemplate:         myRatingsTemplate,
		Tags:                      tags,
	}
}
