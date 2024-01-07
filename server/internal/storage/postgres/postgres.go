package postgres

import (
	"battle-ship_server/internal/service/statistics"
	"battle-ship_server/internal/storage"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	db *pgx.Conn
}

// Prepared statements names
var (
	saveUser = "saveUser"
	getUserData = "getUserData"
	updateStat = "updateStat"
	getStat = "getStat"
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgx.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTablesStmt, err := db.Prepare(context.Background(), "",`
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		login TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);
	
	CREATE TABLE IF NOT EXISTS players_statistics(
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id),
		wins INTEGER NOT NULL DEFAULT 0,
		losses INTEGER NOT NULL DEFAULT 0,
		rating INTEGER NOT NULL DEFAULT 0
	);
	
	CREATE INDEX IF NOT EXISTS idx_player_statistics_user_id ON players_statistics(user_id);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Exec(context.Background(), createTablesStmt.SQL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Prepare statements for future use
	_, err = db.Prepare(context.Background(), saveUser, `
		INSERT INTO users(login, password_hash) VALUES ($1, $2) RETURNING id;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), getUserData, `
		SELECT id, password_hash FROM users WHERE login = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), updateStat, `
		INSERT INTO players_statistics(user_id, wins, losses, rating) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE SET wins = $2, losses = $3, rating = $4;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), getStat, `
		SELECT wins, losses, rating FROM players_statistics WHERE user_id = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(login string, password string) (int, error) {
	const op = "storage.postgres.Login"
	
	var id int
	err := s.db.QueryRow(context.Background(), saveUser, login, password).Scan(&id)
	if err != nil {
		// Check if user already exists
		var pgErr *pgconn.PgError
    	if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_login_key" {
				return -1, storage.ErrUserExists
			}
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUserData(login string) (int, []byte, error) {
	const op = "storage.postgres.GetUserData"

	var id int
	var passHash []byte
	err := s.db.QueryRow(context.Background(), getUserData, login).Scan(&id, &passHash)
	if err != nil {
		if err == pgx.ErrNoRows {
			return -1, nil, storage.ErrUserNotFound
		}
		return -1, nil, fmt.Errorf("%s: %w", op, err)
	}

	return id, passHash, nil
}

func (s *Storage) UpdateStat(id int, stat statistics.Statistics) error {
	const op = "storage.postgres.UpdateStat"

	_, err := s.db.Exec(context.Background(), updateStat, id, stat.Wins, stat.Losses, stat.Rating)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetStat(id int) (statistics.Statistics, error) {
	const op = "storage.postgres.GetStat"

	var stat statistics.Statistics
	err := s.db.QueryRow(context.Background(), getStat, id).Scan(&stat.Wins, &stat.Losses, &stat.Rating)
	if err != nil {
		if err == pgx.ErrNoRows {
			return statistics.Statistics{}, storage.ErrUserNotFound
		}
		return statistics.Statistics{}, fmt.Errorf("%s: %w", op, err)
	}

	return stat, nil
}

func (s *Storage) Close() error {
	return s.db.Close(context.Background())
}