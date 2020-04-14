package service

import "sync"

type Ask struct {
	Price  int64
	Shares int64
	Owner  string
}

type Bid struct {
	Price  int64
	Shares int64
	Owner  string
}

type Symbol struct {
	IPOPrice     int64
	Name         string
	CurrentPrice int64
	Bids         []*Bid
	Asks         []*Ask
	IssuedShares int64
	Volume       int64
	Trades       []*Trade
	Mux          sync.Mutex
}

func NewSymbol(name string, IPOPrice, IssuedShares int64) *Symbol {
	return &Symbol{
		Name:         name,
		IPOPrice:     IPOPrice,
		IssuedShares: IssuedShares,
		Bids:         []*Bid{},
		Asks: []*Ask{
			&Ask{Price: IPOPrice, Shares: IssuedShares, Owner: name},
		},
		Trades: []*Trade{},
	}
}

func (s *Symbol) IPOPriceCurrency() float64 {
	return float64(s.IPOPrice / 100)
}
