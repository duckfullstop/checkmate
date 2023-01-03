package blackjack

import "testing"

func TestPlayerClearHands(t *testing.T) {
	table := NewTable(1)
	player := NewPlayer()
	err := table.Join(player)
	if err != nil {
		t.Error(err)
	}
	errors := table.Deal()
	if len(errors) != 0 {
		// SWEng: this isn't particularly pretty test wise, could probably improve this with more time
		t.Errorf("got errors when processing table dealout")
	}
	// Basically, this shouldn't panic.
	player.clearHands()
}

// newHand is tested by other files.
