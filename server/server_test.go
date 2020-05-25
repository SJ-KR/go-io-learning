package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response, "20")
	})
	t.Run("returns Floyd's score", func(t *testing.T) {

		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response, "10")
	})
}
func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprint("/players/", name), nil)
	return req
}
func assertResponseBody(t *testing.T, response *httptest.ResponseRecorder, want string) {

	got := response.Body.String()

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
