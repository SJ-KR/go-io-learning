package poker

import (
	"io"
	"time"
)

type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string

	StartCalled  bool
	FinishCalled bool
}

type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}

func NewTexasHoldem(alerter BlindAlerter, store PlayerStore) *TexasHoldem {
	return &TexasHoldem{
		alerter: alerter,
		store:   store,
	}
}

func (p *TexasHoldem) Start(numberOfPlayers int) {
	//blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *TexasHoldem) Finish(winner string) {
	p.store.RecordWin(winner)
}

func (p *GameSpy) Start(numberOfPlayers int) {
	p.StartedWith = numberOfPlayers
	p.StartCalled = true
}

func (p *GameSpy) Finish(winner string) {
	p.FinishedWith = winner
	p.FinishCalled = true
}
