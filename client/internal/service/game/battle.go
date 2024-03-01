package gameSrvs

import (
	"battlship/internal/adapters/rabbitmq"
	"errors"
	"fmt"
	"time"
)

var (
	InternalError = errors.New("internal error")
)

type battleMQ interface {
	GetterMessages() (<-chan rabbitmq.Message, error)
	SendMessage(msg rabbitmq.Message) error
}

func (b *BattleShip) StartBattle() error {
	msgs, err := b.mq.GetterMessages()
	if err != nil {
		return err
	}
	b.opponentMsgs = msgs
	return nil
}

func (b *BattleShip) Ready() (iFirst bool, err error) {
	ticker := time.NewTimer(2 * time.Second)
	select {
	case msg := <-b.opponentMsgs:
		fmt.Println("Get message: ", msg)
		if msg.Type != rabbitmq.Ready {
			return false, InternalError
		}
		return false, nil
	case <-ticker.C:
		err = b.mq.SendMessage(rabbitmq.Message{Type: rabbitmq.Ready})
		if err != nil {
			return false, InternalError
		}
		return true, nil
	}
}

func (b *BattleShip) Attack(x, y int) (msgToUser string, err error) {
	msg := rabbitmq.Message{
		Type: rabbitmq.Attack,
		X:    x,
		Y:    y,
	}
	err = b.mq.SendMessage(msg)
	if err != nil {
		return "", InternalError
	}
	resultMsg := <-b.opponentMsgs
	switch resultMsg.Type {
	case rabbitmq.Result:
		b.markHitOrMiss(x, y, resultMsg.Hit, resultMsg.Destroy, false)
		if resultMsg.Destroy {
			msgToUser = "You destroyed the ship!"
		} else if resultMsg.Hit {
			msgToUser = "To the point!"
		} else {
			msgToUser = "You missed"
		}
	case rabbitmq.End:
		msgToUser = Win
	default:
		return "", InternalError
	}
	return msgToUser, nil
}

const (
	Win  = "You win!"
	Lose = "You lose"
)

func (b *BattleShip) Defend() (msgToUser string, err error) {
	msg := <-b.opponentMsgs
	if msg.Type != rabbitmq.Attack {
		return "", InternalError
	}
	x, y := msg.X, msg.Y
	hit, destroy := b.hit(x, y)

	b.markHitOrMiss(x, y, hit, destroy, true)

	var ansMsg rabbitmq.Message
	if b.AllShipsDestroyed() {
		ansMsg = rabbitmq.Message{
			Type: rabbitmq.End,
		}
		msgToUser = Lose
	} else {
		ansMsg = rabbitmq.Message{
			Type:    rabbitmq.Result,
			Hit:     hit,
			Destroy: destroy,
		}
		if destroy {
			msgToUser = "The enemy destroyed the ship"
		} else if hit {
			msgToUser = "The enemy got hit"
		} else {
			msgToUser = "The enemy missed!"
		}
	}
	err = b.mq.SendMessage(ansMsg)
	if err != nil {
		return "", InternalError
	}
	return msgToUser, nil
}
