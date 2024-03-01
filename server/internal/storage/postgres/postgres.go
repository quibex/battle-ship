package postgres

import (
	"battle-ship_server/internal/service/game"
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
	saveUser    = "saveUser"
	getUserData = "getUserData"
	updateStat  = "updateStat"
	getStat     = "getStat"
)

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := pgx.Connect(context.Background(), storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	createTablesStmt := []string{
		`CREATE TABLE IF NOT EXISTS users(
            id SERIAL PRIMARY KEY,
            login TEXT NOT NULL UNIQUE,
            password_hash TEXT NOT NULL
        );`,
		`CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);`,
		`CREATE TABLE IF NOT EXISTS players_statistics(
            id SERIAL PRIMARY KEY,
            user_login TEXT REFERENCES users(login),
            wins INTEGER NOT NULL DEFAULT 0,
            losses INTEGER NOT NULL DEFAULT 0,
            rating INTEGER NOT NULL DEFAULT 0
        );`,
		`CREATE INDEX IF NOT EXISTS idx_player_statistics_user_login ON players_statistics(user_login);`,
		`DROP TRIGGER IF EXISTS create_player_statistics_trigger ON users;

		CREATE OR REPLACE FUNCTION create_player_statistics() RETURNS TRIGGER AS $$
		BEGIN
		  INSERT INTO players_statistics(user_login, wins, losses, rating) 
		  VALUES (NEW.login, 0, 0, 0);
		  RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
		
		CREATE TRIGGER create_player_statistics_trigger
		AFTER INSERT ON users
		FOR EACH ROW
		EXECUTE FUNCTION create_player_statistics();`,
	}

	// batch := pgx.Batch{}

	// for i := range createTablesStmt {
	// 	batch.Queue(createTablesStmt[i])
	// }

	// fmt.Println("Creating tables...")

	// results := db.SendBatch(context.Background(), &batch)
	// defer results.Close()

	// for {
	// 	_, err := results.Exec()
	// 	if err == pgx.ErrNoRows {
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Printf("%s: %v\n", op, err)
	// 		return nil, fmt.Errorf("%s: %w", op, err)
	// 	}
	// }
	for i := range createTablesStmt {
		_, err := db.Exec(context.Background(), createTablesStmt[i])
		if err != nil {
			fmt.Printf("Error executing statement %d: %v\n", i, err)
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	// Prepare statements for future use
	_, err = db.Prepare(context.Background(), saveUser, `
		INSERT INTO users(login, password_hash) VALUES ($1, $2);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), getUserData, `
		SELECT password_hash FROM users WHERE login = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), updateStat, `
		INSERT INTO players_statistics(user_login, wins, losses, rating) VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_login) DO UPDATE SET wins = $2, losses = $3, rating = $4;
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = db.Prepare(context.Background(), getStat, `
		SELECT wins, losses, rating FROM players_statistics WHERE user_login = $1
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(login string, passHash []byte) error {
	const op = "storage.postgres.Login"

	_, err := s.db.Exec(context.Background(), saveUser, login, passHash)
	if err != nil {
		// Check if user already exists
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_login_key" {

				return storage.ErrUserExists
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserData(login string) ([]byte, error) {
	const op = "storage.postgres.GetUserData"

	var passHash []byte
	err := s.db.QueryRow(context.Background(), getUserData, login).Scan(&passHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return passHash, nil
}

func (s *Storage) UpdateStat(userLogin string, stat game.Statistics) error {
	const op = "storage.postgres.UpdateStat"

	_, err := s.db.Exec(context.Background(), updateStat, userLogin, stat.Wins, stat.Losses, stat.Rating)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetStat(login string) (game.Statistics, error) {
	const op = "storage.postgres.GetStat"

	var stat game.Statistics
	err := s.db.QueryRow(context.Background(), getStat, login).Scan(&stat.Wins, &stat.Losses, &stat.Rating)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return game.Statistics{}, storage.ErrUserNotFound
		}
		return game.Statistics{}, fmt.Errorf("%s: %w", op, err)
	}

	return stat, nil
}

func (s *Storage) Close() error {
	return s.db.Close(context.Background())
}
