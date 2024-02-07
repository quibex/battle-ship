package game

import (
	"log/slog"
	"sync"

	"github.com/streadway/amqp"
)

type Service struct {
	Storage StatStorage
	log     *slog.Logger
	games   map[string]game // user name -> game
	mu      sync.RWMutex
}

type gameStatus int

const (
	wait       gameStatus = iota // 0
	inProgress                   // 1
)

type Statistics struct {
	Wins   int
	Losses int
	Rating int
}

type StatStorage interface {
	UpdateStat(login string, stat Statistics) error
	GetStat(login string) (Statistics, error)
}

type game struct {
	user1  string
	dUser1 *amqp.Delivery
	user2  string
	dUser2 *amqp.Delivery
	status gameStatus
}

func New(log *slog.Logger) *Service {
	return &Service{log: log, games: make(map[string]game)}
}

func (s *Service) CreateGame(userName string, dUser *amqp.Delivery) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.games[userName] = game{
		user1:  userName,
		dUser1: dUser,
		status: wait,
	}
	return nil
}

func (s *Service) DelGame(userName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.games, userName)
	return nil
}

func (s *Service) GetAvailableGames() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	games := make([]string, 0)
	for userName, game := range s.games {
		if game.status == wait {
			games = append(games, userName)
		}
	}
	return games, nil
}

func (s *Service) JoinGame(creatorUserName, joiningUserName string, dJoiningUser *amqp.Delivery) (dCreatorUserName *amqp.Delivery, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ugame := s.games[creatorUserName]
	s.games[creatorUserName] = game{
		user1:  ugame.user1,
		dUser1: ugame.dUser1,
		user2:  joiningUserName,
		dUser2: dJoiningUser,
		status: inProgress,
	}
	return ugame.dUser1, nil
}

func (s *Service) GameResult(winner string, loser string) error {
	const op = "Service.ProcessGameResult"

	log := s.log.With(
		slog.String("op", op),
		slog.String("winner", winner),
		slog.String("loser", loser),
	)

	winnerStat, err := s.Storage.GetStat(winner)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	loserStat, err := s.Storage.GetStat(loser)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	winnerStat.Wins++
	winnerStat.Rating += loserStat.Rating / 10
	loserStat.Losses++
	loserStat.Rating -= winnerStat.Rating / 10

	err = s.Storage.UpdateStat(winner, winnerStat)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = s.Storage.UpdateStat(loser, loserStat)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

func (s *Service) GetUserStat(userName string) (Statistics, error) {
	const op = "Service.GetUserStat"

	log := s.log.With(
		slog.String("op", op),
		slog.String("user_name", userName),
	)

	stat, err := s.Storage.GetStat(userName)
	if err != nil {
		log.Error(err.Error())
		return Statistics{}, err
	}

	return stat, nil
}
