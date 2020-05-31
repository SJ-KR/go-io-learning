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
	StartCalled bool
	StartedWith int

	BlindAlert []byte

	FinishCalled bool
	FinishedWith string
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

func (p *TexasHoldem) Start(numberOfPlayers int, alertsDestination io.Writer) {
	//blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Second

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind, alertsDestination)
		blindTime = blindTime + blindIncrement
	}
}

func (p *TexasHoldem) Finish(winner string) {
	p.store.RecordWin(winner)
}

func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartedWith = numberOfPlayers
	out.Write(g.BlindAlert)
}

func (p *GameSpy) Finish(winner string) {
	p.FinishedWith = winner
	p.FinishCalled = true
}
