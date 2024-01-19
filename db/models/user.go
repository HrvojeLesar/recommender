package models

import (
	"log"
	"sort"
)

const (
	RECOMMENDATIONLIMIT = 20
)

type User struct {
	Id          int          `bson:"_id"`
	BookRatings []BookRating `bson:"book_ratings"`
	GenreScores []GenreScore
}

func (u *User) NeighbourSimilarity(other *User) float64 {
	similarityScore := 0.0
	for i := range u.BookRatings {
		bookRating := &u.BookRatings[i]
		for j := range other.BookRatings {
			otherBookRating := &other.BookRatings[j]
			if bookRating.Book.BookId == otherBookRating.Book.BookId {
				similarityScore += bookRating.Distance(otherBookRating)
			}
		}
	}

	return similarityScore / float64(len(u.BookRatings))
}

type UserSimilarity struct {
	User            User
	SimilarityScore float64
}

type SimilarUsers struct {
	user       User
	otherUsers []UserSimilarity
}

func NewSimilarUsers(user User) SimilarUsers {
	return SimilarUsers{
		user:       user,
		otherUsers: make([]UserSimilarity, 0, 10),
	}
}

func (su *SimilarUsers) TryAdd(otherUser User) bool {
	if su.user.Id == otherUser.Id {
		return false
	}
	score := su.user.NeighbourSimilarity(&otherUser)
	if score > 0.0 {
		su.otherUsers = append(su.otherUsers, UserSimilarity{
			User:            otherUser,
			SimilarityScore: score,
		})
		return true
	}
	return false
}

func (su *SimilarUsers) sortOtherUsers() {
	sort.SliceStable(su.otherUsers, func(i, j int) bool {
		return su.otherUsers[i].SimilarityScore > su.otherUsers[j].SimilarityScore
	})
}

func (su *SimilarUsers) OtherUsers() []UserSimilarity {
	return su.otherUsers
}

func (su *SimilarUsers) MostSimilarUser() {
	su.sortOtherUsers()
	log.Println(su.otherUsers[0].User.Id)
}

func (su *SimilarUsers) BookRecommendations() []Book {
	recommendations := su.filterRecommendations()
	sort.SliceStable(recommendations, func(i, j int) bool {
		return recommendations[i].Rating > recommendations[j].Rating
	})

	return recommendations
}

func (su *SimilarUsers) isBookRatedByUser(bookId int) bool {
	for i := range su.user.BookRatings {
		bookRating := &su.user.BookRatings[i]
		if bookRating.Book.BookId == bookId {
			return true
		}
	}
	return false
}

func (su *SimilarUsers) filterRecommendations() []Book {
	su.sortOtherUsers()

	recommendations := make([]Book, 0, RECOMMENDATIONLIMIT)
	for i := range su.otherUsers {
		otherUser := &su.otherUsers[i]
		for j := range otherUser.User.BookRatings {
			bookRating := &otherUser.User.BookRatings[j]
			if !su.isBookRatedByUser(bookRating.Book.BookId) {
				recommendations = append(recommendations, bookRating.Book)
				if len(recommendations) == RECOMMENDATIONLIMIT {
					return recommendations
				}
			}
		}
	}

	return recommendations
}
