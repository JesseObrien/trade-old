package service

import (
	"github.com/google/uuid"
)

type Trader struct {
	ID             string
	Name           string
	AvailableFunds int64
	Trades         []Trade
}

func NewTrader(name string, startingFunds int64) *Trader {
	return &Trader{
		ID:             uuid.New().String(),
		Name:           name,
		AvailableFunds: startingFunds,
		Trades:         []Trade{},
	}

}
