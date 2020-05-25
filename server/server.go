package server

import (
	"fmt"
	"net/http"
	"strings"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Println(r.URL.Path)
	fmt.Println(player)
	if player == "Pepper" {
		fmt.Fprint(w, "20")
	}
	if player == "Floyd" {
		fmt.Fprint(w, "10")
	}
}
