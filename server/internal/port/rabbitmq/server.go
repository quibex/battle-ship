package rabbitmq

import (
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/streadway/amqp"
)

var (
	ErrBadRequest = errors.New("bad request")
	ErrInternal   = errors.New("internal error")
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	log  *slog.Logger

	auth authService
	game gameService
}

func New(urlRmq string, log *slog.Logger, auth authService, game gameService) *RabbitMQ {
	conn, err := amqp.Dial(urlRmq)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMQ{conn: conn, ch: ch, log: log, auth: auth, game: game}
}

func (r *RabbitMQ) Run() {
	go r.Login()
	go r.Register()
	go r.CreateGame()
	go r.DelGame()
	go r.GetAvailableGames()
	go r.JoinGame()
	go r.GameResult()
	go r.GetUserStat()
}

func (r *RabbitMQ) sendResp(d amqp.Delivery, response any) {
	const op = "RabbitMQ.sendResp"

	log := r.log.With(
		slog.String("op", op),
	)

	body, err := json.Marshal(response)
	if err != nil {
		log.Error("Failed to marshal response: %v", err)
		return
	}

	err = r.ch.Publish(
		"",        // exchange
		d.ReplyTo, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			CorrelationId: d.CorrelationId,
			ContentType:   "application/json",
			Body:          body,
		})
	if err != nil {
		log.Error("Failed to publish response: %v", err)
		return
	}
}

func (r *RabbitMQ) Close() error {
	if err := r.ch.Close(); err != nil {
		return err
	}
	if err := r.conn.Close(); err != nil {
		return err
	}
	return nil
}
