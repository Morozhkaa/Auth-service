package memory

import (
	"auth/internal/domain/models"
	"context"
)

func (s *MemoryStorage) SaveUser(ctx context.Context, user models.User) error {
	s.storage[user.Login] = user
	return nil
}

func (s *MemoryStorage) GetUser(ctx context.Context, login string) (models.User, error) {
	if user, ok := s.storage[login]; ok {
		return user, nil
	}
	return models.User{}, models.ErrNotFound
}
