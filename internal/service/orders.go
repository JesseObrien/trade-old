package service

type Order interface {
	Symbol() string
	Trader() *Trader
	Shares() int64
}

type MarketOrder struct {
	symbol string
	trader *Trader
	shares int64
}

func NewMarketOrder(symbol string, trader *Trader, shares int64) *MarketOrder {
	return &MarketOrder{
		symbol,
		trader,
		shares,
	}
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
