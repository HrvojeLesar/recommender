package models

type Tag struct {
	Id   int    `bson:"tag_id"`
	Name string `bson:"tag_name"`
}
