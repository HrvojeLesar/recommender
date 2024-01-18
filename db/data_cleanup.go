package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tag struct {
	TagId   int    `bson:"tag_id"`
	TagName string `bson:"tag_name"`
}

func (m *MongoInstance) iterateMostPopularTags() {
	tagsCollection := m.database.Collection("most_popular_tags")
	tagsCursor, err := tagsCollection.Find(m.ctx, bson.D{})
	if err != nil {
		log.Panicln(err)
	}

	defer tagsCursor.Close(m.ctx)

	for tagsCursor.Next(m.ctx) {
		var tag Tag
		if err := tagsCursor.Decode(&tag); err != nil {
			log.Panicln(err)
		}

		ids := m.findAllGoodreadsBookIdsWithTag(tag)
		m.updateBooks(ids, tag)
	}
}

func (m *MongoInstance) findAllGoodreadsBookIdsWithTag(tag Tag) []int {
	bookTagsCollection := m.database.Collection("book_tags")
	matchStage := bson.D{
		{"$match", bson.D{
			{"tag_id", tag.TagId},
		}},
	}
	projectStage := bson.D{
		{"$project", bson.D{
			{"goodreads_book_id", 1},
		}},
	}

	goodreadsBooksIdsCursor, err := bookTagsCollection.Aggregate(m.ctx, mongo.Pipeline{matchStage, projectStage})
	if err != nil {
		log.Panicln(err)
	}
	defer goodreadsBooksIdsCursor.Close(m.ctx)

	type bookId struct {
		Id int `bson:"goodreads_book_id"`
	}

	ids := make([]int, 0, 50)

	for goodreadsBooksIdsCursor.Next(m.ctx) {
		var id bookId
		if err := goodreadsBooksIdsCursor.Decode(&id); err != nil {
			log.Panicln(err)
		}
		ids = append(ids, id.Id)
	}

	return ids
}

func (m *MongoInstance) updateBooks(ids []int, tag Tag) {
	booksCollection := m.database.Collection("books")
	result, err := booksCollection.UpdateMany(m.ctx, bson.D{
		{"goodreads_book_id", bson.D{
			{"$in", ids},
		}}},
		bson.D{{
			"$addToSet", bson.D{
				{"genres", tag.TagName},
			},
		}},
	)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(result.ModifiedCount)
}
