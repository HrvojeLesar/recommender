package models

import "math"

type Book struct {
	BookId                  int      `bson:"book_id"`
	GoodreadsBookId         int      `bson:"goodreads_book_id"`
	Authors                 []string `bson:"authors"`
	OriginalPublicationYear int      `bson:"original_publication_year"`
	OriginalTitle           string   `bson:"original_title"`
	Title                   string   `bson:"title"`
	Genres                  []string `bson:"genres"`
	Rating                  float64  `bson:"average_rating"`
	ImageUrl                string   `bson:"image_url"`
}

type AverageBookRating struct {
	Book          Book    `bson:"book"`
	AverageRating float64 `bson:"averageRating"`
}

type BookRating struct {
	Book   Book `bson:"book"`
	Rating int  `bson:"rating"`
}

func (br *BookRating) Distance(other *BookRating) float64 {
	distance := int64(math.Abs(float64(br.Rating) - float64(other.Rating)))
	switch distance {
	case 0:
		return 1.0
	case 1:
		return 0.8
	case 2:
		return 0.6
	case 3:
		return 0.4
	case 4:
		return 0.2
	default:
		return 0
	}
}
