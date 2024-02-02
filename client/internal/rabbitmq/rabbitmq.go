package rabbitmq

import "github.com/streadway/amqp"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	id      int // id of this client
}

func New(url string) *RabbitMQ {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return &RabbitMQ{
		conn:    conn,
		channel: ch,
	}
}