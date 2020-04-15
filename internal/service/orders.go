package service

import "github.com/google/uuid"

type Order interface {
	GetID() string
	Symbol() string
	Trader() *Trader
	Shares() int64
}

type MarketOrder struct {
	ID     string
	symbol string
	trader *Trader
	shares int64
}

func NewMarketOrder(symbol string, trader *Trader, shares int64) *MarketOrder {
	return &MarketOrder{
		ID:     uuid.New().String(),
		symbol: symbol,
		trader: trader,
		shares: shares,
	}
}

func (m *MarketOrder) GetID() string {
	return m.ID
}

func (m *MarketOrder) Symbol() string {
	return m.symbol
}

func (m *MarketOrder) Trader() *Trader {
	return m.trader
}

func (m *MarketOrder) Shares() int64 {
	return m.shares
}
