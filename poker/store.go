package poker

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// ======================= StubPlayerStore =======================
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func NewStubPlayerStore() *StubPlayerStore {
	return &StubPlayerStore{}
}
func (s *StubPlayerStore) GetWinCalls() []string {
	return s.winCalls
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

// ======================= InMemoryPlayerStore =======================
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

// ======================= FileSystemPlayerStore =======================

type FileSystemPlayerStore struct {
	Database *json.Encoder
	League   League
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

}

func (f *FileSystemPlayerStore) GetLeague() League {

	sort.Slice(f.League, func(i, j int) bool {
		return f.League[i].Wins > f.League[j].Wins
	})
	return f.League
}
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

	err := initialisePlayerDBFile(file)
	if err != nil {
		return nil, fmt.Errorf("problem initialising player db file, %v", err)
	}

	league, err := NewLeague(file)

	if err != nil {
		return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
	}

	return &FileSystemPlayerStore{
		Database: json.NewEncoder(&tape{file}),
		League:   league,
	}, nil
}
func initialisePlayerDBFile(file *os.File) error {
	file.Seek(0, 0)

	info, err := file.Stat()

	if err != nil {
		return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
	}

	if info.Size() == 0 {
		file.Write([]byte("[]"))
		file.Seek(0, 0)
	}

	return nil
}
