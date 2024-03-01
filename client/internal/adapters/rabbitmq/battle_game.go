package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type MessageType int

const (
	Ready MessageType = iota
	Attack
	Result
	End // send by losing user
)

// Message structure:
// Ready { empty };
// Attack {X, Y int }
// Result { Hit, Destroy bool }
type Message struct {
	Type    MessageType `json:"type"`
	X       int         `json:"x,omitempty"`
	Y       int         `json:"y,omitempty"`
	Hit     bool        `json:"hit,omitempty"`
	Destroy bool        `json:"destroy,omitempty"`
}

func (r *RabbitMQ) GetterMessages() (<-chan Message, error) {
	msgs := make(chan Message)
	go func() {
		for d := range r.msgs {
			var msg Message
			err := json.Unmarshal(d.Body, &msg)
			if err != nil {
				continue
			}
			msgs <- msg
			fmt.Println("Get message: ", msg)
		}
		close(msgs)
	}()
	fmt.Println("Get messages ", r.player1Login)
	return msgs, nil
}

func (r *RabbitMQ) SendMessage(msg Message) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = r.ch.Publish(
		"",
		r.player2Login,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			ReplyTo:     r.player1Login,
		},
	)
	fmt.Println("Sending message to ", r.player2Login, " ", msg)
	if err != nil {
		return err
	}
	return nil
}
