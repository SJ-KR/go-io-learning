package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}

type PlayerServer struct {
	Store PlayerStore
	http.Handler
}
type Player struct {
	Name string
	Wins int
}

func (p *PlayerServer) GetLeague() League {
	return nil
}
func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)
	p.Store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playerHandler))

	p.Handler = router
	return p
}

func (p *PlayerServer) GetPlayerScore(name string) int {
	if name == "Pepper" {
		return 20
	}
	if name == "Floyd" {
		return 10
	}
	return 0
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	_ = json.NewEncoder(w).Encode(p.Store.GetLeague())
	w.WriteHeader(http.StatusOK)

}
func (p *PlayerServer) getLeagueTable() League {

	leagueTable := p.Store.GetLeague()
	return leagueTable
}
func (p *PlayerServer) playerHandler(w http.ResponseWriter, r *http.Request) {

	player := r.URL.Path[len("/players/"):]
	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}
func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {

	p.Store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)

}
func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {

	score := p.Store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	fmt.Fprintf(w, "%d", score)
}

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

/*
type InMemoryPlayerStore struct {
	store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}
func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}
func (i *InMemoryPlayerStore) GetLeague() League {
	var league League

	for name, wins := range i.store {
		league = append(league, Player{name, wins})
	}
	return league
}
*/
type FileSystemPlayerStore struct {
	Database *json.Encoder
	League   League
}

func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
	file.Seek(0, 0)
	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		Database: json.NewEncoder(&tape{file}),
		League:   league,
	}, nil
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
	player := f.League.Find(name)

	if player != nil {
		return player.Wins
	}

	return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {

	player := f.League.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.League = append(f.League, Player{name, 1})
	}
	/*
		f.Database.Seek(0, 0)
		json.NewEncoder(f.Database).Encode(f.League)
	*/
}

func (f *FileSystemPlayerStore) GetLeague() League {

	return f.League
}

func NewLeague(file *os.File) (League, error) {
	var league League
	err := json.NewDecoder(file).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}

type League []Player

func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

type tape struct {
	file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
	t.file.Truncate(0)
	t.file.Seek(0, 0)
	return t.file.Write(p)
}
