package rabbitmq

import (
	"encoding/json"
	"log/slog"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	log *slog.Logger

	auth authService
	authExchange string

	game gameService
	gameExchange string
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

func sendError(ch *amqp.Channel, d amqp.Delivery, response any) {
    body, err := json.Marshal(response)
    if err != nil {
        slog.Error("Failed to marshal error response: %v", err)
        return
    }

    err = ch.Publish(
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
        slog.Error("Failed to publish error message: %v", err)
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