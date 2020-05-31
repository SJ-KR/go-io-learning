package poker

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLeague(t *testing.T) {

	t.Run("it returns 200 on /league", func(t *testing.T) {
		store := StubPlayerStore{}
		server, _ := NewPlayerServer(&store, dummyGame)

		request, _ := http.NewRequest(http.MethodGet, "/league", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got []Player

		err := json.NewDecoder(response.Body).Decode(&got)

		if err != nil {
			t.Fatalf("Unable to parse, %q, %v", response.Body, got)
		}
		assertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Chris", 20},
			{"Tiest", 14},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server, _ := NewPlayerServer(&store, dummyGame)

		request := newLeagueRequest()

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getLeagueFromResponse(t, response.Body)
		assertStatus(t, response.Code, http.StatusOK)
		assertLeague(t, got, wantedLeague)
		assertContentType(t, response, jsonContentType)

	})
}
func getLeagueFromResponse(t *testing.T, body io.Reader) (league []Player) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&league)

	if err != nil {
		t.Fatalf("Unable to parse response from cmd %q into slice of Player, '%v'", body, err)
	}
	return
}
func assertLeague(t *testing.T, got, wantedLeague []Player) {
	t.Helper()

	if !reflect.DeepEqual(got, wantedLeague) {
		t.Errorf("got %v want %v", got, wantedLeague)
	}
}
func newLeagueRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "/league", nil)
	return req
}
func assertContentType(t *testing.T, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}
