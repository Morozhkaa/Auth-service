package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Login    string             `bson:"login"`
	Password string             `bson:"password"`
	Email    string             `bson:"email"`
}

type VerifyResponse struct {
	AccessToken  string
	RefreshToken string
	Login        string
	Email        string
}
