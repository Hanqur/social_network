package entity

import (
	"time"

	"github.com/google/uuid"
)

type SexType string

const (
	MaleSexType   SexType = "male"
	FemaleSexType SexType = "female"
)

func (s SexType) String() string {
	return string(s)
}

type User struct {
	ID         uuid.UUID `db:"id"`
	FirstName  string    `db:"first_name"`
	SecondName string    `db:"second_name"`
	BirthDate  time.Time `db:"birthdate"`
	Sex        SexType   `db:"sex"`
	Biography  string    `db:"biography"`
	City       string    `db:"city"`
	Password   string    `db:"password"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type CreateUserOpts struct {
	FirstName  string
	SecondName string
	BirthDate  time.Time
	Sex        SexType
	Biography  string
	City       string
	Password   string
}
