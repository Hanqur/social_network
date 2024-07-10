package service

import (
	"context"
	"social/internal/database/postgresql"
	"social/internal/entity"

	"github.com/google/uuid"
)

type UserSvc struct {
	storage *postgresql.Storage
}

func NewUserSvc(storage *postgresql.Storage) *UserSvc {
	return &UserSvc{
		storage: storage,
	}
}

func (s *UserSvc) Get(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	return s.storage.Get(ctx, userID)
}
