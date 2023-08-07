package ports

import (
	"auth/internal/domain/models"
	"context"
)

type Auth interface {
	Login(ctx context.Context, login, password string) (string, string, error)
	Verify(ctx context.Context, access, refresh string) (models.VerifyResponse, error)
}
