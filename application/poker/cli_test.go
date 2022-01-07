package poker_test

import (
	"bytes"
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

type GameSpy struct {
	StartedWith int
	FinishedWith string
}

func (g *GameSpy) Start(numberOfPlayers int) {
	g.StartedWith = numberOfPlayers
}

func (g *GameSpy) Finish(winner string) {
	g.FinishedWith = winner
}

func (s scheduledAlert) String() string {
	return fmt.Sprintf("%d chips at %v", s.amount, s.at)
}

func (s *SpyBlindAlerter) ScheduledAlertAt(at time.Duration, amount int) {
	s.alerts = append(s.alerts, scheduledAlert{at, amount})
}

func TestCli(t *testing.T) {
	dummyAlerter := &SpyBlindAlerter{}
	dummyPlayerStore := &poker.StubPlayerStore{}
	dummyStdOut := &bytes.Buffer{}
	t.Run("record Chris win from user input", func(t *testing.T) {
		in := strings.NewReader("5\nChris wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummyAlerter, playerStore)
		cli := poker.NewCli(in, dummyStdOut, game)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Chris")
	})
	t.Run("record Cloe win from user input", func(t *testing.T) {
		in := strings.NewReader("1\nCloe wins\n")
		playerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(dummyAlerter, playerStore)
		cli := poker.NewCli(in, dummyStdOut, game)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Cloe")
	})
	t.Run("it prompts the user to enter the number of players", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		in := strings.NewReader("7\n")
		blindAlerter := &SpyBlindAlerter{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)

		cli := poker.NewCli(in, stdout, game)
		cli.PlayPoker()
		got := stdout.String()
		want := poker.PlayerPrompt
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}

		for i, want := range cases {
			t.Run(fmt.Sprint(want), func(t *testing.T) {
				if len(blindAlerter.alerts) <= i {
					t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
				}

				got := blindAlerter.alerts[i]
				assertScheduledAlert(t, got, want)
			})
		}
	})
}

func TestGameStart(t *testing.T) {
	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		dummyPlayerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)
		game.Start(5)

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
		checkSchedulingCases(cases, t, blindAlerter)
	})

	t.Run("schedules alerts on game start for 5 players", func(t *testing.T) {
		blindAlerter := &SpyBlindAlerter{}
		dummyPlayerStore := &poker.StubPlayerStore{}
		game := poker.NewGame(blindAlerter, dummyPlayerStore)
		game.Start(7)

		cases := []scheduledAlert{
			{0 * time.Second, 100},
			{12 * time.Minute, 200},
			{24 * time.Minute, 300},
			{36 * time.Minute, 400},
		}
		checkSchedulingCases(cases, t, blindAlerter)
	})
}

func TestGameFinish(t *testing.T) {
	store := &poker.StubPlayerStore{}
	dummyAlerter := &SpyBlindAlerter{}
	game := poker.NewGame(dummyAlerter, store)
	winner := "Ruth"
	game.Finish(winner)
	poker.AssertPlayerWins(t, store, winner)
}

func checkSchedulingCases(cases []scheduledAlert, t *testing.T, blindAlerter *SpyBlindAlerter) {
	for i, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			if len(blindAlerter.alerts) <= 1 {
				t.Fatalf("alert %d was not scheduled %v", i, blindAlerter.alerts)
			}
		})

		got := blindAlerter.alerts[i]
		assertScheduledAlert(t, got, c)
	}
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