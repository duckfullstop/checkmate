// This package implements integration and functional tests against the entire package.
package test

import (
	"errors"
	"github.com/duckfullstop/checkmate/pkg/blackjack"
	"github.com/duckfullstop/checkmate/pkg/playdeck"
	"testing"
)

func initTestGame() (table *blackjack.Table, player *blackjack.Player, err error) {
	table = blackjack.NewTable(1)
	if err != nil {
		return nil, nil, err
	}
	player = blackjack.NewPlayer()
	err = table.Join(player)
	// Shorthand on the return, saves a needless if err check
	return table, player, err
}

func TestBBC1(t *testing.T) {
	// Given I play a game of blackjack
	// When I am dealt my opening hand
	// Then I have two cards
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}
	if len(player.Hands[0].Cards) != 2 {
		t.Errorf("player does not have 2 cards after dealing opening hand, got %v", len(player.Hands[0].Cards))
	}
}

func TestBBC2(t *testing.T) {
	// Given I have a valid hand of cards
	// When I choose to 'hit'
	// Then I receive another card
	// And my score is updated
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}
	// Make a note of the existing score to ensure score does update.
	// Proving the minimum score updated is sufficient to ensure that the routine runs.
	_, minScore, _, _ := player.Hands[0].Score()
	// When I choose to hit...
	err = player.Hands[0].Hit()
	if err != nil {
		t.Error(err)
	}
	// I receive another card...
	if len(player.Hands[0].Cards) != 3 {
		t.Errorf("player does not have 3 cards after hit, got %v", len(player.Hands[0].Cards))
	}
	// and my score is updated.
	_, nMinScore, _, _ := player.Hands[0].Score()
	if nMinScore == minScore {
		t.Error("hand score was not updated")
	}
}

func TestBBC3(t *testing.T) {
	// Given I have a valid hand of cards
	// When I choose to 'stand'
	// Then I receive no further cards
	// And my score is updated
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}
	// When I choose to stand...
	err = player.Hands[0].Stick()
	if err != nil {
		t.Error(err)
	}
	// Then I (can) receive no further cards...
	err = player.Hands[0].Hit()
	if !errors.Is(err, blackjack.ErrHandLocked) {
		t.Errorf("unexpected state when attempting to pull further cards after stick, expected HandLocked got %v", err)
	}
	// and my score is updated.
	// This isn't really provable due to how the code is written. Discussion point?
}

func TestBBC4(t *testing.T) {
	// Given my score is updated or evaluated
	// When it is 21 or less
	// Then I have a valid hand
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}

	// We manually inject cards into the hand here to test the integration
	player.Hands[0].Cards = []playdeck.Card{
		{
			Suit:  playdeck.SuitClub,
			Value: playdeck.ValueAce,
		},
		{
			Suit:  playdeck.SuitClub,
			Value: playdeck.ValueKing,
		},
	}
	err = player.Hands[0].EvalScore()
	if err != nil {
		t.Error(err)
	}
	// We should now have a score of exactly 21
	score, _, _, valid := player.Hands[0].Score()
	if score != 21 {
		t.Errorf("expected hand score of 21, got %v", score)
	}
	// and, most importantly, it should be valid
	if !valid {
		t.Error("hand not valid, even though it should be")
	}
}

func TestBBC5(t *testing.T) {
	// Given my score is updated
	// When it is 22 or more
	// Then I am 'bust' and do not have a valid hand
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}

	// We manually inject cards into the hand here to test the integration
	player.Hands[0].Cards = []playdeck.Card{
		{
			Suit:  playdeck.SuitClub,
			Value: playdeck.ValueTen,
		},
		{
			Suit:  playdeck.SuitClub,
			Value: playdeck.ValueTen,
		},
		{
			Suit:  playdeck.SuitDiamond,
			Value: playdeck.ValueTwo,
		},
	}
	err = player.Hands[0].EvalScore()
	if err != nil {
		t.Error(err)
	}
	// We should now have a score of exactly 22
	score, _, _, valid := player.Hands[0].Score()
	if score != 22 {
		t.Errorf("expected hand score of 22, got %v", score)
	}
	// and, most importantly, it should not be valid
	if valid {
		t.Error("hand valid, even though it should not be")
	}
}

func TestBBC6(t *testing.T) {
	// Given I have a king and an ace
	// When my score is evaluated
	// Then my score is 21
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}

	// We manually inject cards into the hand here to test the integration
	player.Hands[0].Cards = []playdeck.Card{
		{
			Suit:  playdeck.SuitSpade,
			Value: playdeck.ValueKing,
		},
		{
			Suit:  playdeck.SuitSpade,
			Value: playdeck.ValueAce,
		},
	}
	err = player.Hands[0].EvalScore()
	if err != nil {
		t.Error(err)
	}
	// We should now have a score of exactly 21
	score, _, _, valid := player.Hands[0].Score()
	if score != 21 {
		t.Errorf("expected hand score of 21, got %v", score)
	}
	// and, additionally, it should be valid
	if !valid {
		t.Error("hand not valid, even though it should be")
	}
}

func TestBBC7(t *testing.T) {
	// Given I have a king, a queen, and an ace
	// When my score is evaluated
	// Then my score is 21
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}

	// We manually inject cards into the hand here to test the integration
	player.Hands[0].Cards = []playdeck.Card{
		{
			Suit:  playdeck.SuitClub,
			Value: playdeck.ValueKing,
		},
		{
			Suit:  playdeck.SuitSpade,
			Value: playdeck.ValueQueen,
		},
		{
			Suit:  playdeck.SuitDiamond,
			Value: playdeck.ValueAce,
		},
	}
	err = player.Hands[0].EvalScore()
	if err != nil {
		t.Error(err)
	}
	// We should now have a score of exactly 21
	score, _, _, valid := player.Hands[0].Score()
	if score != 21 {
		t.Errorf("expected hand score of 21, got %v", score)
	}
	// and, additionally, it should be valid
	if !valid {
		t.Error("hand not valid, even though it should be")
	}
}

func TestBBC8(t *testing.T) {
	// Given that I have a nine, an ace, and another ace
	// When my score is evaluated
	// Then my score is 21
	table, player, err := initTestGame()
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error("table.Deal returned errors")
	}

	// We manually inject cards into the hand here to test the integration
	player.Hands[0].Cards = []playdeck.Card{
		{
			Suit:  playdeck.SuitSpade,
			Value: playdeck.ValueNine,
		},
		{
			Suit:  playdeck.SuitSpade,
			Value: playdeck.ValueAce,
		},
		{
			Suit:  playdeck.SuitHeart,
			Value: playdeck.ValueAce,
		},
	}
	err = player.Hands[0].EvalScore()
	if err != nil {
		t.Error(err)
	}
	// We should now have a score of exactly 21
	score, _, _, valid := player.Hands[0].Score()
	if score != 21 {
		t.Errorf("expected hand score of 21, got %v", score)
	}
	// and, additionally, it should be valid
	if !valid {
		t.Error("hand not valid, even though it should be")
	}
}
