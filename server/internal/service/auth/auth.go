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
	SaveUser(login string, password []byte) error
	GetUserData(login string) ([]byte, error)
}

type Service struct {
	Storage UserStorage
	log     *slog.Logger
}

func New(Storage UserStorage, log *slog.Logger) *Service {
	return &Service{Storage: Storage, log: log}
}

func (s *Service) Register(login string, password string) error {
	const op = "Service.Register"

	log := s.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = s.Storage.SaveUser(login, passHash)
	if errors.Is(err, storage.ErrUserExists) {
		log.Info("user already exists")
		return err
	} else if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("user registered")
	return nil
}

func (s *Service) Login(login string, password string) error {
	const op = "Service.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("login", login),
	)

	passHash, err := s.Storage.GetUserData(login)
	if errors.Is(err, storage.ErrUserNotFound) {
		log.Info("user not found")
		return err
	} else if err != nil {
		log.Error(err.Error())
		return err
	}

	err = bcrypt.CompareHashAndPassword(passHash, []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		log.Info("wrong password")
		return ErrWrongPass
	} else if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("user logged in")
	return nil
}

func (s *Service) Logout(login string) error {
	// nothing to do here
	return nil
}
