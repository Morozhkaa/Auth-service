package db

import (
	"auth/internal/domain/models"
	"auth/pkg/infra/metrics"
	"context"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *MongoStorage) GetUser(ctx context.Context, login string) (models.User, error) {
	ctx, span := metrics.FollowSpan(ctx)
	defer span.End()

	var u models.User
	users := s.Database("auth").Collection("users")
	err := users.FindOne(ctx, bson.D{{Key: "login", Value: login}}).Decode(&u)
	if err != nil {
		return models.User{}, models.ErrNotFound
	}
	return u, nil
}

func (s *MongoStorage) SaveUser(ctx context.Context, user models.User) error {
	ctx, span := metrics.FollowSpan(ctx)
	defer span.End()

	users := s.Database("auth").Collection("users")
	_, err := users.InsertOne(ctx, bson.M{
		"login":    user.Login,
		"password": user.Password,
		"email":    user.Email,
	})
	return err
}
