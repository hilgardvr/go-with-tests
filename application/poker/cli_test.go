package poker_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hilgardvr/go-with-tests/application/poker"
)

type SpyBlindAlerter struct {
	alerts []struct{
		scheduledAt time.Duration
		amount int
	}
}

func (s *SpyBlindAlerter) ScheduledAlertAt(duration time.Duration, amount int) {
	s.alerts = append(s.alerts, struct{scheduledAt time.Duration; amount int}{duration, amount})
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

		cases := []struct{
			expectedScheduledTime time.Duration
			expectedAmount int
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
			t.Run(fmt.Sprintf("%d scheduled for %v", c.expectedAmount, c.expectedScheduledTime), func(t *testing.T) {
				if len(blindAlerter.alerts) <= 1 {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}
			})

			alert := blindAlerter.alerts[i]

			amountGot := alert.amount
			if amountGot != c.expectedAmount {
				t.Errorf("got amount %d, want %d", amountGot, c.expectedAmount)
			}

			gotScheduledTime := alert.scheduledAt
			if gotScheduledTime != c.expectedScheduledTime {
				t.Errorf("got scheduled time %v, want %v", gotScheduledTime, c.expectedScheduledTime)
			}
		}
	})
}