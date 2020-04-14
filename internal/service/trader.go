package service

type Trader struct {
	Name           string
	AvailableFunds int64
	Trades         []Trade
}

func NewTrader(name string, startingFunds int64) *Trader {
	return &Trader{
		Name:           name,
		AvailableFunds: startingFunds,
		Trades:         []Trade{},
	}

}
