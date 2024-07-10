package postgresql

import (
	"embed"
	"fmt"
	"social/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	dialectPostgres    = "postgres"
	gooseMigrationsDir = "migrations"

	usersTable = "users"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Storage struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (*Storage, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.DatabaseName,
	)

	dbx, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	if err := dbx.Ping(); err != nil {
		panic(err)
	}

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(dialectPostgres); err != nil {
		return nil, err
	}

	if err := goose.Up(dbx.DB, gooseMigrationsDir); err != nil {
		return nil, err
	}

	return &Storage{db: dbx}, nil
}
