package poker_test

import (
	"strings"
	"testing"

	poker "github.com/ryanyogan/tdd/player_server"
)

func TestCLI(t *testing.T) {
	t.Run("record ryan win from user input", func(t *testing.T) {
		in := strings.NewReader("Ryan wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Ryan")
	})

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := poker.NewCLI(playerStore, in)
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})
}
