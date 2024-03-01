package rabbitmq

import (
	"github.com/streadway/amqp"
	"time"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	msgs <-chan amqp.Delivery
	que  amqp.Queue

	timeout time.Duration

	player1Login string
	player2Login string
}

func New(url string, timeout time.Duration) *RabbitMQ {
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
		ch:      ch,
		timeout: timeout,
	}
}

func (r *RabbitMQ) initQueue() {
	q, err := r.ch.QueueDeclare(
		r.player1Login, // name
		false,          // durable
		false,          // delete when unused
		true,           // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	if err != nil {
		panic(err)
	}

	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	r.que = q
	r.msgs = msgs
}

func (r *RabbitMQ) Close() {
	_ = r.ch.Cancel(r.que.Name, false)
	_, _ = r.ch.QueueDelete(r.que.Name, false, false, false)
	_ = r.ch.Close()
	_ = r.conn.Close()

}
