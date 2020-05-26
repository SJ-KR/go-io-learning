package main

import (
	sv "go-io-learning/server"
	"log"
	"net/http"
)

func main() {
	store := sv.NewInMemoryPlayerStore()

	server := &sv.PlayerServer{Store: store}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
