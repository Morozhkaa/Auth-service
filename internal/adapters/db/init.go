package db

import (
	"auth/internal/domain/models"
	"auth/internal/ports"
	"context"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

type MongoStorage struct {
	*mongo.Client
}

func New(ctx context.Context, conn string) (ports.UserStorage, error) {
	opts := options.Client().ApplyURI(conn)
	opts.Monitor = otelmongo.NewMonitor()

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	s := &MongoStorage{Client: client}

	// initial users
	user1 := models.User{Login: "Olenka", Password: "Z55+a4XKM9tBClCi4/z9soOEK1/th6bWGveVqhgZTth3uoXAt+afxpy9m77Mo+y7LHJVKipxbOQL1u90V9oceaHaQATc9DH5UB8SEtYg/I6NKyrrnjQasdy7NBN++6834ZErQEsA6+9DmIr4ER3H2ecnQbXiRjBHQ5M2hzvTqc8=", Email: "Olya@mail.ru"}
	s.SaveUser(ctx, user1)

	user2 := models.User{Login: "Katya", Password: "b7t2PnSfjgF7/7Rr+gvOd5whra5HP7q9bV6AXp5sdRfQN0R4ashgfSr6hXi8KxkWQVf3ebmOAngocSc6Wo9HOX/I6OxIACEptozQ4eOwC0PR15ZO3w5SlOWMe6+wyjaJwOdOjhcPHQ1cP5DxxkWlIY+p/7XjcqHUNMzYdYQss8I=", Email: "Katya@mail.ru"}
	s.SaveUser(ctx, user2)

	return s, err
}
