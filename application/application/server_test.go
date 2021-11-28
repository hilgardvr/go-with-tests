package application

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)
type StubPlayerStore struct {
	scores map[string]int
	winCalls []string
}

func (s *StubPlayerStore)GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore)RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
		[]string{},
	}

	server := NewPlayerServer(&store)

	t.Run("return players score", func(t *testing.T) {
		request := getScoreRequest("Pepper")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})
	t.Run("return Floyd's score", func(t *testing.T) {
		request := getScoreRequest("Floyd")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})
	t.Run("return 404 on missing players", func(t * testing.T) {
		request := getScoreRequest("Apollo")
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := response.Code
		want := http.StatusNotFound
		assertStatus(t, got, want)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
	}
	server := NewPlayerServer(&store)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Pepper"
		request := getPostWinRequest(player)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusAccepted)
		if len(store.winCalls) != 1 {
			t.Fatalf("got %d, want %d win calls to store", len(store.winCalls), 1)
		}
		if store.winCalls[0] != player {
			t.Errorf("got %s, want %s winner in store", store.winCalls[0], player)
		}
	})
}

func TestLeague(t *testing.T) {
	store := StubPlayerStore{}
	server := NewPlayerServer(&store)

	t.Run("returns 200 on /league", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusOK)
	})
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{map[string]int{}}
	server := NewPlayerServer(&store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))

	response := httptest.NewRecorder()

	server.ServeHTTP(response, getScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}

func getPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func getScoreRequest(name string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return request
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("response code is wrong, got %d, want %d", got, want)
	}
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q, want %q", got, want)
	}
}