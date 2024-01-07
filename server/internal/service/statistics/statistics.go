package statistics

import (
	"log/slog"
)

type Statistics struct {
	Wins   int
	Losses int
	Rating int
}

type StatStorage interface {
	UpdateStat(id int, stat Statistics) error
	GetStat(id int) (Statistics, error)
}

type Service struct {
	Storage StatStorage
	log     *slog.Logger
}

func New(Storage StatStorage, log *slog.Logger) *Service {
	return &Service{Storage: Storage, log: log}
}

func (s *Service) ProcessGameResult(winnerId int, loserId int) error {
	const op = "Service.ProcessGameResult"

	log := s.log.With(
		slog.String("op", op),
		slog.Int("winnerID", winnerId),
		slog.Int("loserID", loserId),
	)

	winnerStat, err := s.Storage.GetStat(winnerId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	loserStat, err := s.Storage.GetStat(loserId)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	winnerStat.Wins++
	winnerStat.Rating += loserStat.Rating / 10
	loserStat.Losses++
	loserStat.Rating -= winnerStat.Rating / 10

	err = s.Storage.UpdateStat(winnerId, winnerStat)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = s.Storage.UpdateStat(loserId, loserStat)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
