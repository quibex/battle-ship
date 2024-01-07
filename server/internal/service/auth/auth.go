package auth

import (
	"battle-ship_server/internal/storage"
	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrWrongPass = errors.New("wrong password")
)

type UserStorage interface {
	SaveUser(login string, password []byte) (int, error)
	GetUserData(login string) (int, []byte, error)
}

type Service struct {
	Storage UserStorage
	log     *slog.Logger
}

func New(Storage UserStorage, log *slog.Logger) *Service {
	return &Service{Storage: Storage, log: log}
}

func (s *Service) Register(login string, password string) (int, error) {
	const op = "Service.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	id, err := s.Storage.SaveUser(login, passHash)
	if err == storage.ErrUserExists {
		log.Info("user already exists")
		return -1, err
	} else if err != nil {
		log.Error(err.Error())
		return -1, err
	}
	return id, nil
}

func (s *Service) Login(login string, password string) (int, error) {
	const op = "Service.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	id, passHash, err := s.Storage.GetUserData(login)
	if err == storage.ErrUserNotFound {
		log.Info("user not found")
		return -1, err
	} else if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	err = bcrypt.CompareHashAndPassword(passHash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		log.Info("wrong password")
		return -1, ErrWrongPass
	} else if err != nil {
		log.Error(err.Error())
		return -1, err
	}

	return id, nil
}