package gameSrvs

import (
	"battlship/internal/service/game/domain"
	"context"
)

type serverMQ interface {
	CreateGame(ctx context.Context) (user2 string, err error)
	DelGame() error
	JoinGame(creatorUserName string) error
	GetAvailableGames() (games []string, err error)
	SaveGameResult(winner, loser string) error
	GetUserStat(username string) (domain.Statistics, error)
	GetOpponentName() (string, error)
}

func (b *BattleShip) CreateGame(ctx context.Context) (user2 string, err error) {
	return b.mq.CreateGame(ctx)
}

func (b *BattleShip) DelGame() error {
	return b.mq.DelGame()
}

func (b *BattleShip) JoinGame(creatorUserName string) error {
	return b.mq.JoinGame(creatorUserName)
}

func (b *BattleShip) GetAvailableGames() (games []string, err error) {
	return b.mq.GetAvailableGames()
}

func (b *BattleShip) SaveGameResult(winner, loser string) error {
	return b.mq.SaveGameResult(winner, loser)
}

func (b *BattleShip) GetUserStat(username string) (domain.Statistics, error) {
	return b.mq.GetUserStat(username)
}

func (b *BattleShip) GetOpponentName() (string, error) {
	return b.mq.GetOpponentName()
}
