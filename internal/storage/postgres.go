package storage

import (
	"auth/internal/config"
	"auth/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfg config.Config) (*Storage, error) {
	const op = "storage.postgres.NewPostgresDB"

	db, err := sql.Open("postgres",
		fmt.Sprintf("port=%s user=%s host=%s password=%s dbname=%s sslmode=%s",
			cfg.DBConfig.Port, cfg.DBConfig.Username, cfg.DBConfig.Host, cfg.DBConfig.Password, cfg.DBConfig.DBName, cfg.DBConfig.SSLMode),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db}, nil
}

func (s *Storage) SaveUser(ctx context.Context, name string, username string, passHash []byte) (int, error) {
	const op = "storage.postgres.SaveUser"

	query := fmt.Sprintf(`INSERT INTO users (name, username, password_hash) VALUES ($1, $2, $3) RETURNING id`)

	var id int

	err := s.db.QueryRow(query, name, username, passHash).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Storage) User(ctx context.Context, username string) (models.User, error) {
	const op = "storage.postgres.User"

	query := fmt.Sprintf(`SELECT id, name, username, password_hash FROM users WHERE username=$1`)

	var user models.User

	err := s.db.QueryRow(query, username).Scan(&user.ID, &user.Name, &user.Username, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
