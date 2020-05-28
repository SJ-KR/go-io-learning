package poker_test

import (
	"go-io-learning/poker"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})
}

func assertPlayerWin(t *testing.T, store *poker.StubPlayerStore, winner string) {
	t.Helper()
	wincall := store.GetWinCalls()
	if len(wincall) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(wincall), 1)
	}

	if wincall[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", wincall[0], winner)
	}
}
