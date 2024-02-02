package rabbitmq

import (
	"encoding/json"
	"errors"

	"github.com/streadway/amqp"
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
	ClientID int    `json:"client_id"`
	Err      string `json:"error,omitempty"`
}

type registerResponse struct {
	ClientID int    `json:"client_id"`
	Err      string `json:"error,omitempty"`
}


func (r *RabbitMQ) Login(login, password string) error {
	q, err := r.channel.QueueDeclare(
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

	err = r.channel.Publish(
		"", // exchange
		"login",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     q.Name,
		},
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return err
	}

	d, ok := <-msgs
	if !ok {
		return errors.New("client login channel closed")
	}

	var response loginResponse
	err = json.Unmarshal(d.Body, &response)
	if err != nil {
		return err
	}
	if response.Err != "" {
		return errors.New(response.Err)
	}

	r.id = response.ClientID
	r.channel.Cancel(q.Name, false)
	return nil
}

func (r *RabbitMQ) Register(login, password string) error {
	q, err := r.channel.QueueDeclare(
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

	err = r.channel.Publish(
		"", // exchange
		"register",      // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     q.Name,
		},
	)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		q.Name, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return err
	}

	d, ok := <-msgs
	if !ok {
		return errors.New("client register channel closed")
	}

	var response registerResponse
	err = json.Unmarshal(d.Body, &response)
	if err != nil {
		return err
	}
	if response.Err != "" {
		return errors.New(response.Err)
	}

	r.id = response.ClientID
	r.channel.Cancel(q.Name, false)
	return nil
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}
