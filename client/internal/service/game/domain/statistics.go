package domain

import "fmt"

type Statistics struct {
	Wins   int
	Losses int
	Rating int
}

func (s Statistics) String() string {
	return fmt.Sprintf("Wins: %d, Losses: %d, Rating: %d", s.Wins, s.Losses, s.Rating)
}
