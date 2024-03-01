package rabbitmq

import (
	"battle-ship_server/internal/service/auth"
	"battle-ship_server/internal/storage"
	"encoding/json"
	"errors"
	"log/slog"
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

type authService interface {
	Login(username, password string) error
	Register(username, password string) error
	Logout(username string) error
}

func (r *RabbitMQ) Login() {
	const op = "RabbitMQ.Login"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"auth.login", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Error("Failed to declare a queue: %v", err)
		return
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
		log.Error("Failed to register a login consumer: %v", err)
		return
	}

	for d := range msgs {
		var request loginRequest
		err := json.Unmarshal(d.Body, &request)
		if err != nil {
			r.log.Error("Failed to unmarshal login request: %v", err)
			r.sendResp(d, ErrBadRequest)
			continue
		}

		err = r.auth.Login(request.Username, request.Password)
		if errors.Is(err, auth.ErrWrongPass) {
			r.sendResp(d, loginResponse{Err: err.Error()})
			continue
		} else if errors.Is(err, storage.ErrUserNotFound) {
			r.sendResp(d, loginResponse{Err: err.Error()})
			continue
		} else if err != nil {
			r.sendResp(d, loginResponse{Err: ErrInternal.Error()})
			continue
		}

		r.sendResp(d, loginResponse{})
	}
}

func (r *RabbitMQ) Register() {
	const op = "RabbitMQ.Register"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"auth.register", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		log.Error("Failed to declare a queue: %v", err)
		return
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
		log.Error("Failed to register a consumer: %v", err)
		return
	}

	for d := range msgs {
		var request registerRequest
		err := json.Unmarshal(d.Body, &request)
		if err != nil {
			r.log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, ErrBadRequest)
			continue
		}

		err = r.auth.Register(request.Username, request.Password)
		if errors.Is(err, storage.ErrUserExists) {
			r.sendResp(d, registerResponse{Err: err.Error()})
			continue
		} else if err != nil {
			r.sendResp(d, registerResponse{Err: ErrInternal.Error()})
			continue
		}

		r.sendResp(d, registerResponse{})
	}
}

func (r *RabbitMQ) Logout() {
	const op = "RabbitMQ.Logout"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"logout", // name
		false,    // durable
		false,    // delete when unused
		false,    // exclusive
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Error("Failed to declare a queue: %v", err)
		return
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
		log.Error("Failed to register a consumer: %v", err)
		return
	}

	for d := range msgs {
		var clientID int
		err := json.Unmarshal(d.Body, &clientID)
		if err != nil {
			r.log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, ErrBadRequest)
			continue
		}
	}
}
