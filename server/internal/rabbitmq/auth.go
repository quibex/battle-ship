package rabbitmq

import (
	"battle-ship_server/internal/storage"
	"battle-ship_server/internal/service/auth"
	"encoding/json"
	"errors"
	"log/slog"

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
	ClientID int `json:"client_id"`
	Err string `json:"error,omitempty"`
}

type registerResponse struct {
	ClientID int `json:"client_id"`
	Err string `json:"error,omitempty"`
}

var (
	ErrBadRequest = errors.New("bad request")
	ErrInternal  = errors.New("internal error")
)

type authService interface {
	Login(username string, password string) (int, error)
	Register(username, password string) (int, error)
}

func (r *RabbitMQ) Login() {
	const op = "RabbitMQ.Login"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"login", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
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
			sendError(r.ch, d, ErrBadRequest)
			continue
		}

		userID, err := r.auth.Login(request.Username, request.Password)
		if err == auth.ErrWrongPass {
			sendError(r.ch, d, loginResponse{Err: err.Error()})
			continue
		} else if err == storage.ErrUserNotFound {
			sendError(r.ch, d, loginResponse{Err: err.Error()})
			continue
		} else if err != nil {
			sendError(r.ch, d, loginResponse{Err: ErrInternal.Error()})
			continue
		}

		clientID := r.game.GetClientID(userID)
		response := loginResponse{ClientID: clientID}
		body, err := json.Marshal(response)
		if err != nil {
			r.log.Error("Failed to marshal login response: %v", err)
			sendError(r.ch, d, ErrInternal.Error())
			continue
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
			r.log.Error("Failed to publish response: %v", err)
			sendError(r.ch, d, ErrInternal.Error())
			continue
		}
	}
}

func (r *RabbitMQ) Register() {
	const op = "RabbitMQ.Register"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"register", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
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
			sendError(r.ch, d, ErrBadRequest)
			continue
		}

		userID, err := r.auth.Register(request.Username, request.Password)
		if err == storage.ErrUserExists {
			sendError(r.ch, d, registerResponse{Err: err.Error()})
			continue
		} else if err != nil {
			sendError(r.ch, d, registerResponse{Err: ErrInternal.Error()})
			continue
		}

		clientID := r.game.GetClientID(userID)

		response := registerResponse{ClientID: clientID}
		body, err := json.Marshal(response)
		if err != nil {
			r.log.Error("Failed to marshal response: %v", err)
			sendError(r.ch, d, ErrInternal.Error())
			continue
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
			r.log.Error("Failed to publish response: %v", err)
			sendError(r.ch, d, ErrInternal.Error())
			continue
		}
	}
}