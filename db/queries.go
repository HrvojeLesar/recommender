package db

import (
	"github.com/HrvojeLesar/recommender/db/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	roundProjectStage := bson.D{
		{"$project", bson.D{
			{"_id", 1},
			{"book", 1},
			{"totalRating", 1},
			{"count", 1},
			{"averageRating", bson.D{
				{"$round", bson.A{"$averageRating", 2}},
			}},
		}},
	}

	booksCursor, err := usersCollection.Aggregate(m.ctx, mongo.Pipeline{unwindStage, groupStage, projectStage, sortStage, roundProjectStage})
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

		if bookRating.Book.Rating1 == 0 {
			continue
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

	roundProjectStage := bson.D{
		{"$project", bson.D{
			{"_id", 1},
			{"book", 1},
			{"totalRating", 1},
			{"count", 1},
			{"averageRating", bson.D{
				{"$round", bson.A{"$averageRating", 2}},
			}},
		}},
	}

	booksCursor, err := usersCollection.Aggregate(m.ctx, mongo.Pipeline{unwindStage, groupStage, projectStage, matchStage, sortStage, roundProjectStage})
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

		if bookRating.Book.Rating1 == 0 {
			continue
		}

		if len(books) < 30 {
			books = append(books, bookRating)
		} else {
			break
		}
	}

	return books, nil
}

func (m *MongoInstance) NearestNeighbour(userId int64) ([]models.Book, error) {
	similarUsers, err := m.SimilarUsers(userId)
	if err != nil {
		return nil, err
	}
	return similarUsers.BookRecommendations(), nil
}

func (m *MongoInstance) Tags() ([]models.Tag, error) {
	tagsCollection := m.database.Collection(TAGSCOLLECTION)

	tagsCursor, err := tagsCollection.Find(m.ctx, bson.D{}, options.Find().SetSort(bson.D{{"tag_name", 1}}))
	if err != nil {
		return nil, err
	}
	defer tagsCursor.Close(m.ctx)

	tags := make([]models.Tag, 0, 10)
	for tagsCursor.Next(m.ctx) {
		var tag models.Tag
		if err := tagsCursor.Decode(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}

func (m *MongoInstance) MyRatings(userId int64) (*models.User, error) {
	usersCollectionns := m.database.Collection(USERSCOLLECTION)

	userResult := usersCollectionns.FindOne(m.ctx, bson.D{{"_id", userId}})
	var user models.User
	err := userResult.Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (m *MongoInstance) SimilarUsers(userId int64) (*models.SimilarUsers, error) {
	usersCollection := m.database.Collection(USERSCOLLECTION)
	var user models.User

	result := usersCollection.FindOne(m.ctx, bson.D{{"_id", userId}})
	err := result.Decode(&user)
	if err != nil {
		return nil, err
	}

	otherUsersCursor, err := usersCollection.Find(m.ctx, bson.D{
		{"_id", bson.D{{
			"$ne", userId,
		}}},
	})
	if err != nil {
		return nil, err
	}
	defer otherUsersCursor.Close(m.ctx)

	similarUsers := models.NewSimilarUsers(user)
	for otherUsersCursor.Next(m.ctx) {
		var otherUser models.User
		if err := otherUsersCursor.Decode(&otherUser); err != nil {
			return nil, err
		}

		similarUsers.TryAdd(otherUser)
	}

	return &similarUsers, nil
}

func (m *MongoInstance) InsertRating(newRating int64, userId int64, book models.Book) error {
	usersCollection := m.database.Collection(USERSCOLLECTION)
	book.Rating1 = 0

	_, err := usersCollection.UpdateOne(m.ctx, bson.D{{"_id", userId}}, bson.D{
		{"$push", bson.D{
			{"book_ratings", bson.D{
				{"rating", newRating},
				{"book", book},
			}},
		}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoInstance) UpdateRating(newRating int64, userId int64, bookId int64) error {
	usersCollection := m.database.Collection(USERSCOLLECTION)

	arrayFilters := options.ArrayFilters{
		Filters: []interface{}{bson.D{
			{"x.book.book_id", bookId},
		}},
	}

	_, err := usersCollection.UpdateOne(m.ctx, bson.D{{"_id", userId}}, bson.D{
		{"$set", bson.D{
			{"book_ratings.$[x].rating", newRating},
		}},
	}, options.Update().SetArrayFilters(arrayFilters))
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoInstance) RemoveRating(newRating int64, userId int64, bookId int64) error {
	usersCollection := m.database.Collection(USERSCOLLECTION)

	_, err := usersCollection.UpdateOne(m.ctx, bson.D{{"_id", userId}}, bson.D{
		{"$pull", bson.D{
			{"book_ratings", bson.D{
				{"book.book_id", bookId},
			}},
		}},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MongoInstance) Book(bookId int64) (*models.Book, error) {
	booksCollection := m.database.Collection(BOOKSCOLLECTION)

	var book models.Book

	result := booksCollection.FindOne(m.ctx, bson.D{{"book_id", bookId}})
	err := result.Decode(&book)
	if err != nil {
		return nil, err
	}

	return &book, nil
}

func (m *MongoInstance) Search(searchTerm string) ([]models.Book, error) {
	booksCollection := m.database.Collection(BOOKSCOLLECTION)

	booksCursor, err := booksCollection.Find(m.ctx, bson.D{{"title", bson.D{
		{"$regex", searchTerm},
		{"$options", "i"},
	}}}, options.Find().SetLimit(100))
	if err != nil {
		return nil, err
	}
	defer booksCursor.Close(m.ctx)

	books := make([]models.Book, 0, 10)
	for booksCursor.Next(m.ctx) {
		var book models.Book
		err := booksCursor.Decode(&book)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}
