package main

import (
	"go-io-learning/server"
	"log"
	"net/http"
)

type InMemoryPlayerStore struct {
	store server.PlayerStore
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {

	return 123
}

func main() {
	server := server.NewPlayerServer(&InMemoryPlayerStore{})

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
