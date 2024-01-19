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
	Tags                      []models.Tag
}

func NewWebserverHandler(globalInstances *global.Instance) Handler {
	indexTemplate := template.Must(template.ParseFiles("templates/index.html", "templates/header.html"))
	personalizedIndexTemplate := template.Must(template.ParseFiles("templates/index_perso.html", "templates/header.html"))
	tags, err := globalInstances.Mongo.Tags()
	if err != nil {
		log.Panicln(err)
	}
	return Handler{
		globalInstances:           globalInstances,
		indexTemplate:             indexTemplate,
		personalizedIndexTemplate: personalizedIndexTemplate,
		Tags:                      tags,
	}
}
