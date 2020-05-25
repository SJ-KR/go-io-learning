package main

import (
	"go-io-learning/server"
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(server.PlayerServer)
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
