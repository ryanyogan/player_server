package poker

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const jsonContentType = "application/json"

func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Ryan":  20,
			"Floyd": 10,
		},
		[]string{},
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("returns Ryan's score", func(t *testing.T) {
		request := NewGetScoreRequest("Ryan")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusOK)
		AssertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := NewGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
		[]string{},
		nil,
	}
	server := NewPlayerServer(&store)

	t.Run("it records wins on POST", func(t *testing.T) {
		player := "Ryan"

		request := NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		AssertStatus(t, response.Code, http.StatusAccepted)
		AssertPlayerWin(t, &store, "Ryan")
	})
}

func TestLeague(t *testing.T) {
	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []Player{
			{"Cleo", 32},
			{"Ryan", 44},
			{"Mike", 12},
		}

		store := StubPlayerStore{nil, nil, wantedLeague}
		server := NewPlayerServer(&store)

		request := NewLeagueRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := GetLeagueFromRequest(t, response.Body)
		AssertStatus(t, response.Code, http.StatusOK)
		AssertLeague(t, got, wantedLeague)
		AssertContentType(t, response, jsonContentType)
	})
}

func TestFileSystemStore(t *testing.T) {
	t.Run("/league from a reader", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		got := store.GetLeague()
		want := []Player{
			{"Chris", 33},
			{"Ryan", 10},
		}

		AssertNoError(t, err)
		AssertLeague(t, got, want)

		// Check to ensure the file is seeked from 0-char
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})

	t.Run("get player score", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Chris", "Wins": 22}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		got := store.GetPlayerScore("Ryan")

		AssertNoError(t, err)
		AssertScoreEquals(t, got, 10)
	})

	t.Run("store wins for existing player", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		store.RecordWin("Ryan")

		got := store.GetPlayerScore("Ryan")
		AssertNoError(t, err)
		AssertScoreEquals(t, got, 11)
	})

	t.Run("store wins for new players", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Chris", "Wins": 33}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)
		store.RecordWin("Pepper")

		got := store.GetPlayerScore("Pepper")
		AssertNoError(t, err)
		AssertScoreEquals(t, got, 1)
	})

	t.Run("works with an empty file", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, "")
		defer cleanDatabase()

		_, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)
	})

	t.Run("league sorted", func(t *testing.T) {
		database, cleanDatabase := CreateTempFile(t, `[
			{"Name": "Ryan", "Wins": 10},
			{"Name": "Chris", "Wins": 22}
		]`)
		defer cleanDatabase()

		store, err := NewFileSystemPlayerStore(database)

		AssertNoError(t, err)
		got := store.GetLeague()

		want := []Player{
			{"Chris", 22},
			{"Ryan", 10},
		}

		AssertLeague(t, got, want)

		// Read again (seek to ensure order)
		got = store.GetLeague()
		AssertLeague(t, got, want)
	})
}
