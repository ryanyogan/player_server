package poker

import (
	"encoding/json"
	"fmt"
	"io"
)

// League is an alias for a list of Player's
type League []Player

// Find searches a League given a name and returns a player pointer
func (l League) Find(name string) *Player {
	for i, p := range l {
		if p.Name == name {
			return &l[i]
		}
	}
	return nil
}

// NewLeague creates a new League and returns it as json
func NewLeague(rdr io.Reader) ([]Player, error) {
	var league []Player
	err := json.NewDecoder(rdr).Decode(&league)
	if err != nil {
		err = fmt.Errorf("problem parsing league, %v", err)
	}

	return league, err
}
