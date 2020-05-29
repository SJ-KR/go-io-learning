package poker

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

var dummySpyAlerter = &SpyBlindAlerter{}

func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		cli := NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := NewCLI(playerStore, in, dummySpyAlerter)
		cli.PlayPoker()

		AssertPlayerWin(t, playerStore, "Cleo")
	})
}

func assertPlayerWin(t *testing.T, store *StubPlayerStore, winner string) {
	t.Helper()
	wincall := store.GetWinCalls()
	if len(wincall) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(wincall), 1)
	}

	if wincall[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", wincall[0], winner)
	}
}

func TestTime(t *testing.T) {
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		if len(blindAlerter.Alerts) != 1 {
			t.Fatal("expected a blind alert to be scheduled")
		}
	})

	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}

		cli := NewCLI(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []struct {
			expectedScheduleTime time.Duration
			expectedAmount       int
		}{
			{0 * time.Second, 100},
			{10 * time.Minute, 200},
			{20 * time.Minute, 300},
			{30 * time.Minute, 400},
			{40 * time.Minute, 500},
			{50 * time.Minute, 600},
			{60 * time.Minute, 800},
			{70 * time.Minute, 1000},
			{80 * time.Minute, 2000},
			{90 * time.Minute, 4000},
			{100 * time.Minute, 8000},
		}

		for i, c := range cases {
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduleTime), func(t *testing.T) {

				if len(blindAlerter.Alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.Alerts)
				}

				alert := blindAlerter.Alerts[i]

				amountGot := alert.amount
				if amountGot != c.expectedAmount {
					t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
				}

				gotScheduledTime := alert.scheduledAt
				if gotScheduledTime != c.expectedScheduleTime {
					t.Errorf("got scheduled time of %v, want %v", gotScheduledTime, c.expectedScheduleTime)
				}
			})
		}
	})
}
