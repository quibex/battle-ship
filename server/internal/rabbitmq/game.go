package rabbitmq

import (
	"battle-ship_server/internal/service/game"
	"encoding/json"
	"log/slog"

	"github.com/streadway/amqp"
)

type gameService interface {
	CreateGame(userName string, dUser *amqp.Delivery) error
	DelGame(userName string) error
	GetAvailableGames() (games []string, err error)
	JoinGame(creatorUserName, joiningUserName string, dJoiningUser *amqp.Delivery) (dCreatorUserName *amqp.Delivery, err error)
	GameResult(winner, loser string) error
	GetUserStat(userName string) (game.Statistics, error)
}

func (r *RabbitMQ) CreateGame() {
	const op = "RabbitMQ.CreateGame"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"game.create", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
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
		var req gameCreateRequest
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, gameCreateResponse{Err: ErrBadRequest.Error()})
			continue
		}

		err = r.game.CreateGame(req.UserName, &d)
		if err != nil {
			r.sendResp(d, gameCreateResponse{Err: ErrInternal.Error()})
			continue
		}
		// the creator is waiting for another user to join
	}
}

func (r *RabbitMQ) JoinGame() {
	const op = "RabbitMQ.JoinGame"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"game.join", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
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
		var req gameJoinRequest
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, gameJoinResponse{Err: ErrBadRequest.Error()})
			continue
		}

		dUser1, err := r.game.JoinGame(req.CreatorUserName, req.JoiningUserName, &d)
		if err != nil {
			r.sendResp(d, gameJoinResponse{Err: err.Error()})
			continue
		}

		r.sendResp(*dUser1, gameJoinResponse{})
		r.sendResp(d, gameJoinResponse{}) //
	}
}

func (r *RabbitMQ) GetAvailableGames() {
	const op = "RabbitMQ.GetAvailableGames"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"game.get_available_games", // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
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
		var req getAvailableGamesRequest
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, getAvailableGamesResponse{Err: ErrBadRequest.Error()})
			continue
		}

		games, err := r.game.GetAvailableGames()
		if err != nil {
			r.sendResp(d, getAvailableGamesResponse{Err: ErrInternal.Error()})
			continue
		}

		r.sendResp(d, getAvailableGamesResponse{Games: games})
	}
}

func (r *RabbitMQ) GameResult() {
	const op = "RabbitMQ.GameResult"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"game.result", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
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
		var req gameResultRequest
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, gameResultResponse{Err: ErrBadRequest.Error()})
			continue
		}

		err = r.game.GameResult(req.Winner, req.Loser)
		if err != nil {
			r.sendResp(d, gameResultResponse{Err: ErrInternal.Error()})
			continue
		}

		r.sendResp(d, gameResultResponse{})
	}
}

func (r *RabbitMQ) GetUserStat() {
	const op = "RabbitMQ.GetUserStat"

	log := r.log.With(
		slog.String("op", op),
	)

	q, err := r.ch.QueueDeclare(
		"user.get_stat", // name
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
		var req getStatRequest
		err := json.Unmarshal(d.Body, &req)
		if err != nil {
			log.Error("Failed to unmarshal request: %v", err)
			r.sendResp(d, getStatResponse{Err: ErrBadRequest.Error()})
			continue
		}

		stat, err := r.game.GetUserStat(req.UserName)
		if err != nil {
			r.sendResp(d, getStatResponse{Err: ErrInternal.Error()})
			continue
		}

		r.sendResp(d, getStatResponse{
			Rating: stat.Rating,
			Wins:   stat.Wins,
			Losses: stat.Losses,
		})
	}
}
