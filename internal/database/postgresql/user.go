package postgresql

import (
	"context"
	"fmt"
	"social/internal/entity"
	"time"

	"github.com/google/uuid"
)

func (s *Storage) Get(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user := new(entity.User)

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1`, usersTable)

	if err := s.db.QueryRowxContext(ctx, query, id.String()).StructScan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) GetByPass(ctx context.Context, id uuid.UUID, pass []byte) (*entity.User, error) {
	var user *entity.User

	query := fmt.Sprintf(`SELECT * FROM %s WHERE id = $1 AND password = $2`, usersTable)

	if err := s.db.QueryRowContext(ctx, query, id.String(), string(pass)).Scan(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *entity.User) (uuid.UUID, error) {
	creationTime := time.Now()

	query := fmt.Sprintf(`INSERT INTO %s (id, created_at, updated_at, first_name, 
	second_name, birthdate, sex, biography, city, password) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`, usersTable)

	_, err := s.db.ExecContext(ctx, query, user.ID, creationTime, creationTime, user.FirstName,
		user.SecondName, user.BirthDate, user.Sex, user.Biography, user.City, user.Password)
	if err != nil {
		return uuid.Nil, err
	}

	return user.ID, nil
}
