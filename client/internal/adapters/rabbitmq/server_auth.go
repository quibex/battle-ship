package rabbitmq

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"time"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Err string `json:"error,omitempty"`
}

type registerResponse struct {
	Err string `json:"error,omitempty"`
}

func (r *RabbitMQ) Login(login, password string) error {
	q, err := r.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	req := loginRequest{
		Username: login,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",           // exchange
		"auth.login", // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     q.Name,
		},
	)
	if err != nil {
		return err
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
		return err
	}

	timer := time.NewTimer(r.timeout)

	select {
	case d, ok := <-msgs:
		if !ok {
			return errors.New("client playerLogin ch closed")
		}
		var response loginResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return err
		}
		if response.Err != "" {
			return errors.New(response.Err)
		}

		err := r.ch.Cancel(q.Name, false)
		if err != nil {
			return err
		}

		r.player1Login = login
		r.initQueue()
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}

func (r *RabbitMQ) Register(login, password string) error {
	q, err := r.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	req := registerRequest{
		Username: login,
		Password: password,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",              // exchange
		"auth.register", // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     q.Name,
		},
	)
	if err != nil {
		return err
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
		return err
	}

	timer := time.NewTimer(r.timeout)

	select {
	case d, ok := <-msgs:
		if !ok {
			return errors.New("client register ch closed")
		}

		var response registerResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return err
		}
		if response.Err != "" {
			return errors.New(response.Err)
		}

		err = r.ch.Cancel(q.Name, false)
		if err != nil {
			return err
		}

		r.player1Login = login
		r.initQueue()
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}
