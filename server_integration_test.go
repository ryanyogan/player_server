package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecordingWinsAndRetreivingThem(t *testing.T) {
	database, cleanDatabase := CreateTempFile(t, `[]`)
	defer cleanDatabase()
	store, err := NewFileSystemPlayerStore(database)
	server := NewPlayerServer(store)
	player := "Ryan"

	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), NewPostWinRequest(player))

	t.Run("get score", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewGetScoreRequest(player))
		AssertNoError(t, err)
		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "3")
	})

	t.Run("get league", func(t *testing.T) {
		response := httptest.NewRecorder()
		server.ServeHTTP(response, NewLeagueRequest())
		AssertNoError(t, err)
		AssertStatus(t, response.Code, http.StatusOK)

		got := GetLeagueFromRequest(t, response.Body)
		want := []Player{
			{"Ryan", 3},
		}
		AssertLeague(t, got, want)
	})
}
