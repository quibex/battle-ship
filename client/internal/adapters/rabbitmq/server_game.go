package rabbitmq

import (
	"battlship/internal/service/game/domain"
	"context"
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"time"
)

type gameCreateRequest struct {
	UserName string `json:"user_name"`
}

type gameJoinRequest struct {
	CreatorUserName string `json:"creator_user_name"`
	JoiningUserName string `json:"joining_user_name"`
}

type gameDelRequest struct {
	UserName string `json:"user_name"`
}

type gameDelResponse struct {
	Err string `json:"error,omitempty"`
}

type getAvailableGamesRequest struct{}

type gameResultRequest struct {
	Winner string `json:"winner"`
	Loser  string `json:"loser"`
}

type getStatRequest struct {
	UserName string `json:"user_name"`
}

type getStatResponse struct {
	Rating int    `json:"rating"`
	Wins   int    `json:"wins"`
	Losses int    `json:"losses"`
	Err    string `json:"error,omitempty"`
}

type getAvailableGamesResponse struct {
	Games []string `json:"games"`
	Err   string   `json:"error,omitempty"`
}

type gameCreateResponse struct {
	Err   string `json:"error,omitempty"`
	User2 string `json:"user2,omitempty"`
}

type gameJoinResponse struct {
	Err string `json:"error,omitempty"`
}

type gameResultResponse struct {
	Err string `json:"error,omitempty"`
}

const ( //queue names
	gameCreate        = "game.create"
	gameJoin          = "game.join"
	getAvailableGames = "game.get_available"
	saveGameResult    = "game.save_result"
	getUserStat       = "game.get_user_stat"
	gameDel           = "game.del"
)

func (r *RabbitMQ) CreateGame(ctx context.Context) (user2 string, err error) {

	req := gameCreateRequest{
		UserName: r.player1Login,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	err = r.ch.Publish(
		"",         // exchange
		gameCreate, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return "", err
	}

	timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response gameCreateResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return "", err
		}
		if response.Err != "" {
			return "", errors.New(response.Err)
		}
		r.player2Login = response.User2
		return response.User2, nil
	case <-timer.C:
		err = r.DelGame()
		if err != nil {
			return "", err
		}
		return "", errors.New("timeout")
	case <-ctx.Done(): // cancel create game
		err = r.DelGame()
		if err != nil {
			return "", err
		}
		return "", errors.New("context done")
	}
}

func (r *RabbitMQ) DelGame() error {

	req := gameDelRequest{
		UserName: r.player1Login,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",      // exchange
		gameDel, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return err
	}

	timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response gameDelResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return err
		}
		if response.Err != "" {
			return errors.New(response.Err)
		}
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}

func (r *RabbitMQ) GetUserStat(username string) (domain.Statistics, error) {

	req := getStatRequest{
		UserName: username,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return domain.Statistics{}, err
	}

	err = r.ch.Publish(
		"",          // exchange
		getUserStat, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return domain.Statistics{}, err
	}

	//timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response getStatResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return domain.Statistics{}, err
		}
		if response.Err != "" {
			return domain.Statistics{}, errors.New(response.Err)
		}
		return domain.Statistics{
			Rating: response.Rating,
			Wins:   response.Wins,
			Losses: response.Losses,
		}, nil
		//case <-timer.C:
		//	return domain.Statistics{}, errors.New("timeout")
	}
}

func (r *RabbitMQ) JoinGame(creatorUserName string) error {

	req := gameJoinRequest{
		CreatorUserName: creatorUserName,
		JoiningUserName: r.player1Login,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",       // exchange
		gameJoin, // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return err
	}

	timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response gameJoinResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return err
		}
		if response.Err != "" {
			return errors.New(response.Err)
		}
		r.player2Login = creatorUserName
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}

func (r *RabbitMQ) GetAvailableGames() ([]string, error) {

	req := getAvailableGamesRequest{}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	err = r.ch.Publish(
		"",                // exchange
		getAvailableGames, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return nil, err
	}

	timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response getAvailableGamesResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return nil, err
		}
		if response.Err != "" {
			return nil, errors.New(response.Err)
		}
		return response.Games, nil
	case <-timer.C:
		return nil, errors.New("timeout")
	}
}

func (r *RabbitMQ) SaveGameResult(winner, loser string) error {

	req := gameResultRequest{
		Winner: winner,
		Loser:  loser,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",             // exchange
		saveGameResult, // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.que.Name,
		},
	)
	if err != nil {
		return err
	}

	timer := time.NewTimer(r.timeout)
	select {
	case d := <-r.msgs:
		var response gameResultResponse
		err = json.Unmarshal(d.Body, &response)
		if err != nil {
			return err
		}
		if response.Err != "" {
			return errors.New(response.Err)
		}
		return nil
	case <-timer.C:
		return errors.New("timeout")
	}
}

func (r *RabbitMQ) GetOpponentName() (string, error) {
	if r.player2Login == "" {
		return "", errors.New("no opponent")
	}
	return r.player2Login, nil
}
