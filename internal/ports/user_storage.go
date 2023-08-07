package ports

import (
	"auth/internal/domain/models"
	"context"
)

type UserStorage interface {
	GetUser(ctx context.Context, login string) (models.User, error)
	SaveUser(ctx context.Context, user models.User) error
}
