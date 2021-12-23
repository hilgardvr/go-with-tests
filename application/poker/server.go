package poker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const jsonContentType = "application/json"

type Player struct {
	Name string	
	Wins int
}

type League []Player

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	store PlayerStore
	http.Handler
}

type InMemoryPlayerStore struct {
	db map[string]int
}

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.store = store
	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	p.Handler = router
	return p
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

func (i *InMemoryPlayerStore)GetPlayerScore(name string) int {
	return i.db[name]
}

func (i *InMemoryPlayerStore)RecordWin(name string) {
	i.db[name]++
}

func (i *InMemoryPlayerStore)GetLeague() []Player {
	var league []Player
	for name, wins := range i.db {
		league = append(league, Player{name, wins})
	}
	return league
}

func (p *PlayerServer)leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer)getLeagueTable() []Player {
	leagueTable := []Player{
		{"Chris", 20},
	}
	return leagueTable
}

func (p *PlayerServer)playersHandler(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	switch r.Method {
	case http.MethodPost: 
		p.ProcessWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer)ProcessWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}

func (p *PlayerServer)showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)
	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprint(w, score)
}