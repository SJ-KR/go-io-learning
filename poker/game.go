package poker

import (
	"time"
)

type GameIFC interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}

type GameSpy struct {
	StartedWith  int
	FinishedWith string

	StartCalled  bool
	FinishCalled bool
}

type Game struct {
	alerter BlindAlerter
	store   PlayerStore
}

func NewGame(alerter BlindAlerter, store PlayerStore) *Game {
	return &Game{
		alerter: alerter,
		store:   store,
	}
}

func (p *Game) Start(numberOfPlayers int) {
	blindIncrement := time.Duration(5+numberOfPlayers) * time.Minute

	blinds := []int{100, 200, 300, 400, 500, 600, 800, 1000, 2000, 4000, 8000}
	blindTime := 0 * time.Second
	for _, blind := range blinds {
		p.alerter.ScheduleAlertAt(blindTime, blind)
		blindTime = blindTime + blindIncrement
	}
}

func (p *Game) Finish(winner string) {
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
