package models

type Genre = string

type GenreScore struct {
	Genre  Genre   `bson:"genre"`
	Weight float64 `bson:"weight"`
}
