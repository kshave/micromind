package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuestionModel struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	Question string             `json:"question" bson:"question"`
}

type QuoteModel struct {
	ID     primitive.ObjectID `json:"id" bson:"_id"`
	Quote  string             `json:"quote" bson:"quote"`
	Author string             `json:"author" bson:"author"`
}
