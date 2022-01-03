package poker_test

import (
	"strings"
	"testing"

	"github.com/hilgardvr/go-with-tests/application/poker"
)

func TestCli(t *testing.T) {
	t.Run("record Chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCli(playerStore, in)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Chris")
	})
	t.Run("record Cloe win from user input", func(t *testing.T) {
		in := strings.NewReader("Cloe wins\n")
		playerStore := &poker.StubPlayerStore{}
		cli := poker.NewCli(playerStore, in)
		cli.PlayPoker()
		poker.AssertPlayerWins(t, playerStore, "Cloe")
	})
}