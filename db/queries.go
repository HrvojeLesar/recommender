package db

import (
	"log"

	"github.com/HrvojeLesar/recommender/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (m *MongoInstance) TopBooksByUserRating() ([]models.AverageBookRating, error) {
	usersCollection := m.database.Collection(USERSCOLLECTION)

	unwindStage := bson.D{
		{"$unwind", "$book_ratings"},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$book_ratings.book"},
			{"totalRating", bson.D{
				{"$sum", "$book_ratings.rating"},
			}},
			{"count", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", "$_id._id"},
			{"book", "$_id"},
			{"totalRating", 1},
			{"count", 1},
			{"averageRating", bson.D{
				{"$divide", bson.A{
					"$totalRating", "$count",
				}},
			}},
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.D{
			{"averageRating", -1},
		}},
	}

	booksCursor, err := usersCollection.Aggregate(m.ctx, mongo.Pipeline{unwindStage, groupStage, projectStage, sortStage})
	if err != nil {
		return nil, err
	}
	defer booksCursor.Close(m.ctx)

	books := make([]models.AverageBookRating, 0, 30)
	for booksCursor.Next(m.ctx) {
		var bookRating models.AverageBookRating
		if err := booksCursor.Decode(&bookRating); err != nil {
			return nil, err
		}

		if len(books) < 30 {
			books = append(books, bookRating)
		} else {
			break
		}
	}

	return books, nil
}

func (m *MongoInstance) TopBooksByGenre(genre string) ([]models.AverageBookRating, error) {
	usersCollection := m.database.Collection(USERSCOLLECTION)

	unwindStage := bson.D{
		{"$unwind", "$book_ratings"},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", "$book_ratings.book"},
			{"totalRating", bson.D{
				{"$sum", "$book_ratings.rating"},
			}},
			{"count", bson.D{
				{"$sum", 1},
			}},
		}},
	}

	projectStage := bson.D{
		{"$project", bson.D{
			{"_id", "$_id._id"},
			{"book", "$_id"},
			{"totalRating", 1},
			{"count", 1},
			{"averageRating", bson.D{
				{"$divide", bson.A{
					"$totalRating", "$count",
				}},
			}},
		}},
	}

	matchStage := bson.D{
		{"$match", bson.D{
			{"book.genres", bson.D{
				{"$in", bson.A{genre}},
			}},
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.D{
			{"averageRating", -1},
		}},
	}

	booksCursor, err := usersCollection.Aggregate(m.ctx, mongo.Pipeline{unwindStage, groupStage, projectStage, matchStage, sortStage})
	if err != nil {
		return nil, err
	}
	defer booksCursor.Close(m.ctx)

	books := make([]models.AverageBookRating, 0, 30)
	for booksCursor.Next(m.ctx) {
		var bookRating models.AverageBookRating
		if err := booksCursor.Decode(&bookRating); err != nil {
			return nil, err
		}

		if len(books) < 30 {
			books = append(books, bookRating)
		} else {
			break
		}
	}

	return books, nil
}

func (m *MongoInstance) NearestNeighbour(userId int) []models.Book {
	usersCollection := m.database.Collection(USERSCOLLECTION)
	var user models.User

	result := usersCollection.FindOne(m.ctx, bson.D{{"_id", userId}})
	err := result.Decode(&user)
	if err != nil {
		log.Panicln(err)
	}

	otherUsersCursor, err := usersCollection.Find(m.ctx, bson.D{
		{"_id", bson.D{{
			"$ne", userId,
		}}},
	})
	if err != nil {
		log.Panicln(err)
	}
	defer otherUsersCursor.Close(m.ctx)

	similarUsers := models.NewSimilarUsers(user)
	for otherUsersCursor.Next(m.ctx) {
		var otherUser models.User
		if err := otherUsersCursor.Decode(&otherUser); err != nil {
			log.Panicln(err)
		}

		similarUsers.TryAdd(otherUser)
	}
	return similarUsers.BookRecommendations()
}
