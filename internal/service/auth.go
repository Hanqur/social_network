package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"social/internal/database/postgresql"
	"social/internal/entity"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordGenerate = errors.New("password generation")
	ErrPasswordInvalid  = errors.New("invalid password")

	signingKey = []byte("AllYourBase")
)

const (
	tokenTTL     = 12 * time.Hour
	saltPassword = "as;lfj3iru38747938hraaksjhf"
)

type AuthSvc struct {
	storage *postgresql.Storage
}

type TokenClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID `json:"user_id"`
}

func NewAuthSvc(storage *postgresql.Storage) *AuthSvc {
	return &AuthSvc{storage: storage}
}

func (s *AuthSvc) Login(ctx context.Context, id uuid.UUID, password string) (string, error) {
	user, err := s.storage.Get(ctx, id)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), s.hashPassword(password))
	if err != nil {
		return "", ErrPasswordInvalid
	}

	return s.generateToken(user.ID)
}

func (s *AuthSvc) CreateUser(ctx context.Context, userOpts *entity.CreateUserOpts) (uuid.UUID, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(s.hashPassword(userOpts.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, ErrPasswordGenerate
	}

	id := uuid.New()

	user := &entity.User{
		ID:         id,
		FirstName:  userOpts.FirstName,
		SecondName: userOpts.SecondName,
		BirthDate:  userOpts.BirthDate,
		Sex:        userOpts.Sex,
		Biography:  userOpts.Biography,
		City:       userOpts.City,
		Password:   string(hashedPassword),
	}

	return s.storage.CreateUser(ctx, user)
}

func (s *AuthSvc) ParseToken(token string) (uuid.UUID, error) {
	t, err := jwt.ParseWithClaims(token, &TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := t.Claims.(*TokenClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid claims format")
	}

	return claims.UserID, nil
}

func (s *AuthSvc) generateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&TokenClaims{
			jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			userID,
		},
	)

	return token.SignedString(signingKey)
}

func (s *AuthSvc) hashPassword(password string) []byte {
	h := sha1.New()
	h.Write([]byte(password))
	h.Write([]byte(saltPassword))

	return h.Sum(nil)
}
