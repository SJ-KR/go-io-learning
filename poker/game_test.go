package poker

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGame_Start(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(5, ioutil.Discard)

		cases := []scheduledAlert{
			{at: 0 * time.Second, amount: 100},
			{at: 10 * time.Second, amount: 200},
			{at: 20 * time.Second, amount: 300},
			{at: 30 * time.Second, amount: 400},
			{at: 40 * time.Second, amount: 500},
			{at: 50 * time.Second, amount: 600},
			{at: 60 * time.Second, amount: 800},
			{at: 70 * time.Second, amount: 1000},
			{at: 80 * time.Second, amount: 2000},
			{at: 90 * time.Second, amount: 4000},
			{at: 100 * time.Second, amount: 8000},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 7 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		game := NewTexasHoldem(blindAlerter, dummyPlayerStore)

		game.Start(7, ioutil.Discard)

		cases := []scheduledAlert{
			{at: 0 * time.Second, amount: 100},
			{at: 12 * time.Second, amount: 200},
			{at: 24 * time.Second, amount: 300},
			{at: 36 * time.Second, amount: 400},
		}

		checkSchedulingCases(cases, t, blindAlerter)
	})
	t.Run("it prompts the user to enter the number of players and starts the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()
		assertMessagesSentToUser(t, stdout, PlayerPrompt)

		if game.StartedWith != 7 {
			t.Errorf("wanted Start called with 7 but got %d", game.StartedWith)
		}
	})
	t.Run("it prints an error when a non numeric value is entered and does not start the game", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("Pies\n")
		game := &GameSpy{}

		cli := NewCLI(in, stdout, game)
		cli.PlayPoker()

		if game.StartCalled {
			t.Errorf("game should not have started")
		}

	})
}
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server, _ := NewPlayerServer(&StubPlayerStore{}, dummyGame)

		request := newGameRequest()

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})

}
func newGameRequest() *http.Request {
	request, err := http.NewRequest(http.MethodGet, "/game", nil)

	if err != nil {
		fmt.Errorf("%q", err)
		return nil
	}
	return request
}
func checkSchedulingCases(cases []scheduledAlert, t *testing.T, blindAlerter *SpyBlindAlerter) {
	a := blindAlerter.Alerts
	for i, _ := range cases {
		if a[i].amount != cases[i].amount {
			t.Errorf("got %d, want %d", a[i].amount, cases[i].amount)
		}
		if a[i].at != cases[i].at {
			t.Errorf("got %d, want %d", a[i].at, cases[i].at)
		}
	}
}
func TestGame_Finish(t *testing.T) {
	store := &StubPlayerStore{}
	blindAlerter := &SpyBlindAlerter{}
	game := NewTexasHoldem(blindAlerter, store)
	winner := "Ruth"

	game.Finish(winner)
	AssertPlayerWin(t, store, winner)
}
func assertMessagesSentToUser(t *testing.T, stdout *bytes.Buffer, messages ...string) {
	t.Helper()
	want := strings.Join(messages, "")
	got := stdout.String()
	if got != want {
		t.Errorf("got %q sent to stdout but expected %+v", got, messages)
	}
}
