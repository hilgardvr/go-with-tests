package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hilgardvr/go-with-tests/application/poker"
)

type SpyBlindAlerter struct {
	alerts []scheduledAlert
}


type scheduledAlert struct {
	at time.Duration
	amount int
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduledAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

func TestCli(t *testing.T) {
	dummyAlerter := &SpyBlindAlerter{}
	t.Run("record Chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCli(playerStore, in, dummyAlerter)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Chris")
	})
	t.Run("record Cloe win from user input", func(t *testing.T) {
		in := strings.NewReader("Cloe wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCli(playerStore, in, dummyAlerter)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Cloe")
	})
	t.Run("it schedules printing of blind values", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		blindAlerter := &SpyBlindAlerter{}
		cli := poker.NewCli(playerStore, in, blindAlerter)
		cli.PlayPoker()

		cases := []scheduledAlert{
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
			t.Run(fmt.Sprint(c), func(t *testing.T) {
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
			})

			got := blindAlerter.alerts[i]
			assertScheduledAlert(t, got, c)
		}
	})
}

func assertScheduledAlert(t testing.TB, got scheduledAlert, want scheduledAlert) {
		amountGot := got.amount
		if amountGot != want.amount {
			t.Errorf("got amount %d, want %d", amountGot, want.amount)
		}

		gotScheduledTime := got.at
		if gotScheduledTime != want.at {
			t.Errorf("got scheduled time %v, want %v", gotScheduledTime, want.at)
		}
}