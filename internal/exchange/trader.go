package exchange

import (
	"errors"
	"fmt"

	"github.com/apex/log"
	"github.com/google/uuid"
)

type Trader struct {
	Identifier     string `json:"id"`
	Name           string `json:"name"`
	AvailableFunds int64  `json:"funds"`
}

type TraderStore struct {
	logger   log.Logger
	database map[string]*Trader
}

func NewTraderStore(logger log.Logger) *TraderStore {
	return &TraderStore{
		logger:   logger,
		database: make(map[string]*Trader),
	}
}

func (ts *TraderStore) Write(t *Trader) error {
	ts.logger.Infof("🤑 A new trader appears: %s", t.Name)

	t.Identifier = uuid.New().String()
	ts.database[t.Identifier] = t

	return nil
}

func (ts *TraderStore) Read(identifier string, t *Trader) error {
	trader, ok := ts.database[identifier]

	if !ok {
		return errors.New(fmt.Sprintf("trader with id %s not found", identifier))
	}

	t = trader

	return nil
}
