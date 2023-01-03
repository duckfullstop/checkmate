package blackjack

import (
	"errors"
	"fmt"
	"github.com/duckfullstop/checkmate/pkg/playdeck"
	"testing"
)

func TestTableInitialisation(t *testing.T) {
	table := NewTable(1)
	if table.Deck == nil {
		t.Errorf("no deck!")
	}
	if table.sDecks != 1 {
		t.Errorf("deck count setting not set properly, expected 1 got %v", table.sDecks)
	}
	if len(*table.Deck.Cards) != 52 {
		t.Errorf("didn't get expected number of cards, expected 52 got %v", len(*table.Deck.Cards))
	}
}

func TestTablePlayOnePlayer(t *testing.T) {
	table := NewTable(1)
	player := Player{
		Table: table,
	}
	err := table.Join(&player)
	if err != nil {
		t.Error(err)
	}
	// We shouldn't be able to join again
	err = table.Join(&player)
	if !errors.Is(err, ErrTablePlayerAlreadyJoined) {
		t.Errorf("didn't get appropriate error when attempting rejoin, expected PlayerAlreadyJoined got %s", err)
	}

	// we shouldn't be able to end the round if it isn't in play
	err = table.EndRound()
	if !errors.Is(err, ErrTableNotInPlay) {
		t.Errorf("didn't get appropriate error when attempting to end round early, expected HandNotLocked got %s", err)
	}

	errs := table.Deal()
	if len(errs) != 0 {
		// SWEng: this isn't particularly pretty test wise, could probably improve this with more time
		t.Errorf("got errors when processing table dealout")
	}
	// We definitely shouldn't be able to deal twice
	errs = table.Deal()
	if len(errs) != 1 {
		// SWEng: this also isn't particularly pretty test wise
		t.Errorf("got invalid error count when dealing, expected 1 got %v", len(errs))
	}
	if len(player.Hands) != 1 {
		t.Errorf("invalid number of hands dealt, expected 1 got %v", len(player.Hands))
	}
	if len(player.Hands[0].Cards) != 2 {
		t.Errorf("invalid number of cards in hand, expected 2 got %v", len(player.Hands[0].Cards))
	}

	// a player shouldn't be able to join mid-game
	player2 := Player{
		Table: table,
	}
	err = table.Join(&player2)
	if !errors.Is(err, ErrTableInPlay) {
		t.Errorf("didn't get appropriate error when attempting join in progress, expected TableInPlay got %s", err)
	}

	// we shouldn't be able to end the round if a player is still playing
	err = table.EndRound()
	if !errors.Is(err, ErrHandNotLocked) {
		t.Errorf("didn't get appropriate error when attempting to end round early, expected HandNotLocked got %s", err)
	}

	// test to the end of the round
	err = player.Hands[0].Stick()
	if err != nil {
		t.Error(err)
	}
	err = table.EndRound()
	if err != nil {
		t.Errorf("could not end round, got error %s", err)
	}
	t.Logf("test ended with hand state: %s", fmt.Sprintln(player.Hands[0].Score()))
}

// Test a player being bound with no table, forcing the newHand call to fail
func TestTableDealPlayerNoTable(t *testing.T) {
	table := NewTable(1)
	player := NewPlayer()
	table.Players = append(table.Players, player)
	errs := table.Deal()
	if len(errs) != 1 {
		t.Error("invalid length of errors array")
	}
	if !errors.Is(errs[0], ErrPlayerNoTable) {
		t.Errorf("unexpected error thrown with empty player: %s", errs[0])
	}
}

// Test the player having a valid table, but no deck assigned to the table (i.e can't pull a card)
func TestTableDealPlayerNoDeckAssigned(t *testing.T) {
	table := NewTable(0)
	player := NewPlayer()
	err := table.Join(player)
	if err != nil {
		t.Error(err)
	}
	// Intentionally yeet the deck
	table.Deck = &playdeck.Deck{}
	errs := table.Deal()
	if len(errs) == 0 {
		t.Error("invalid length of errors array")
		t.FailNow()
	}
	if !errors.Is(errs[0], playdeck.ErrDeckEmpty) {
		t.Errorf("unexpected error thrown with bad deck: %s", errs[0])
	}
}
