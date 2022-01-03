package poker

import "testing"

type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
	league []Player
}

func (s *StubPlayerStore)GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore)RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore)GetLeague() League {
	return s.league
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response code is wrong, got %d, want %d", got, want)
	}
}

func AssertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}

func AssertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func AssertPlayerWins(t testing.TB, store *StubPlayerStore, winner string) {
	if len(store.winCalls) != 1 {
		t.Fatalf("got %d, want %d win calls to store", len(store.winCalls), 1)
	}
	if store.winCalls[0] != winner {
		t.Errorf("got %s, want %s winner in store", store.winCalls[0], winner)
	}
}