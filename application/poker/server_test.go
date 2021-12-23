package poker

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

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

func TestGetPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd": 10,
		},
		[]string{},
		nil,
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
	data := `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`
			
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
		nil,
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

	t.Run("store wins for an existing player", func(t *testing.T){
		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()
		// store := FileSystemPlayerStore{database}
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		store.RecordWin("Chris")
		got := store.GetPlayerScore("Chris")
		want := 34
		assertScoreEquals(t, got, want)
	})
	t.Run("store wins for new players", func(t *testing.T) {
		database, _ := createTempFile(t, data)
		player := "Pepper"
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		store.RecordWin(player)
		league := store.GetLeague()
		player2 := league.FindName(player)
		println(player2)
		got := store.GetPlayerScore(player)
		want := 1
		assertScoreEquals(t, got, want)
	})
}

func TestLeague(t *testing.T) {
	t.Run("returns the league table as json", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}
		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)
		request := newLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := getLeagueFromResponse(t, response.Body)
		assertLeague(t, got, wantedLeague)
		assertStatus(t, response.Code, http.StatusOK)
		assertContentType(t, response, jsonContentType)
	})
}

func TestFileSystemStore(t *testing.T) {
	data := `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`
	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		got := store.GetLeague()
		want := []Player{
			{"Chris", 33},
			{"Cleo", 10},
		}
		assertLeague(t, got, want)
		//read again
		got = store.GetLeague()
		assertLeague(t, got, want)
	})
	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, data)
		defer cleanDatabase()
		store, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
		got := store.GetPlayerScore("Chris")
		want := 33
		assertScoreEquals(t, got, want)
	})
	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := createTempFile(t, "")
		defer cleanDatabase()
		_, err := NewFileSystemPlayerStore(database)
		assertNoError(t, err)
	})
}

func TestTapeWrite(t *testing.T) {
	file, clean := createTempFile(t, "12345")
	defer clean()
	tape := &tape{file}
	tape.Write([]byte("abc"))
	file.Seek(0,0)
	newFileContents, _ := ioutil.ReadAll(file)
	got := string(newFileContents)
	want := "abc"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertScoreEquals(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d want %d", got, want)
	}
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response %v did not have a content-type of %v", response.Result().Header.Get("content-type"), want)
	}
}

func getLeagueFromResponse(t testing.TB, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)
	if err != nil {
		t.Fatalf("unable to parse response from server %q into slice of Player, '%v'", body, err)
	}
	return
}

func assertLeague(t testing.TB, got, want []Player) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("wanted %v got %v", want, got)
	}
}

func newLeagueRequest() *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return request
}

func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	database, cleanDatabase := createTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	assertNoError(t, err)
	server := NewPlayerServer(store)
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), getPostWinRequest(player))

	t.Run("get score", func(t *testing.T){
		response := httptest.NewRecorder()
		server.ServeHTTP(response, getScoreRequest(player))
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, newLeagueRequest())
		assertStatus(t, response.Code, http.StatusOK)
		got := getLeagueFromResponse(t, response.Body)
		want := []Player{
			{"Pepper", 3},
		}
		assertLeague(t, got, want)
	})
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

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("didn't expect an error but got one, %v", err)
	}
}

func createTempFile(t testing.TB, initialData string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := ioutil.TempFile("", "db")
	if err != nil {
		t.Fatalf("could not create temp file %v", err)
	}
	tmpFile.Write([]byte(initialData))
	removeFile := func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
	return tmpFile, removeFile
}