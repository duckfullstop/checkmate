package blackjack

import (
	"errors"
	"github.com/duckfullstop/checkmate/pkg/playdeck"
	"testing"
)

// Test that the hand fails to create in safe ways.
func TestHandInvalidPlayerInitialisation(t *testing.T) {
	_, err := newHand(nil)
	if !errors.Is(err, ErrPlayerInvalid) {
		t.Errorf("didn't get appropriate error when creating invalid hand, expected PlayerInvalid got %s", err)
	}
	// This player intentionally left blank.
	player := Player{}
	_, err = newHand(&player)
	if !errors.Is(err, ErrPlayerNoTable) {
		t.Errorf("didn't get appropriate error when creating invalid hand, expected PlayerNoTable got %s", err)
	}
}

// Test that the basics of initialising a hand work as intended.
func TestHandInitialisation(t *testing.T) {
	table := NewTable(1)
	player := Player{
		Table: table,
	}

	hand, err := newHand(&player)
	if err != nil {
		t.Error(err)
	}
	if len(hand.Cards) != 0 {
		t.Errorf("newly initialised Hand has %v cards", len(hand.Cards))
	}
	if hand.Player != &player {
		t.Errorf("hand does not have player properly assigned")
	}
	score, minScore, locked, valid := hand.Score()
	if score != 0 {
		t.Errorf("score not zero, got %v", score)
	}
	if minScore != 0 {
		t.Errorf("minScore not zero, got %v", minScore)
	}
	if locked {
		t.Errorf("hand locked, should be unlocked")
	}
	if !valid {
		t.Errorf("hand invalid, should be valid")
	}
	if hand.canPlay() == nil {
		t.Errorf("")
	}
}

// Test that invalid hand handling (heh) performs as expected.
func TestHandInvalid(t *testing.T) {
	// Start by testing cases where the hand is blank
	h := Hand{}
	err := h.canPlay()
	if !errors.Is(err, ErrHandInvalid) {
		t.Errorf("unexpected state when deck not initialised, expected HandInvalid got %v", err)
	}

	// Now test an otherwise valid hand with an invalid player reference
	h = Hand{
		Cards: []playdeck.Card{
			{Suit: playdeck.SuitClub, Value: playdeck.ValueTen},
			{Suit: playdeck.SuitHeart, Value: playdeck.ValueTen},
		},
		// Null deck is fine for this test
		Table:    &Table{Deck: &playdeck.Deck{}},
		score:    10,
		minScore: 10,
		locked:   false,
		valid:    true,
	}
	err = h.canPlay()
	if !errors.Is(err, ErrHandNoPlayer) {
		t.Errorf("unexpected state when deck not bound to player, expected HandNoPlayer got %v", err)
	}
	player := Player{}
	h.Player = &player
	err = h.canPlay()
	if !errors.Is(err, ErrPlayerNoTable) {
		t.Errorf("unexpected state when player not bound to table, expected PlayerInvalid got %v", err)
	}
}

// Test that the deck running out of cards is handled as expected.
func TestDeckOutOfCards(t *testing.T) {
	table := NewTable(1)
	player := Player{
		Table: table,
	}

	hand, err := newHand(&player)
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error(errs)
	}
	// Let's now intentionally run out of cards by creating a new, empty deck...
	table.Deck.Cards = &[]playdeck.Card{}
	err = hand.Hit()
	if !errors.Is(err, playdeck.ErrDeckEmpty) {
		t.Errorf("unexpected state when deck out of cards, expected DeckEmpty got %v", err)
	}
}

// Test that sticking works as expected, and sets the hand to a completed state.
func TestHandStick(t *testing.T) {
	table := NewTable(1)
	player := Player{
		Table: table,
	}
	err := table.Join(&player)
	if err != nil {
		t.Error(err)
	}
	errs := table.Deal()
	if len(errs) != 0 {
		t.Error(errs)
	}
	hand := player.Hands[0]
	err = hand.Stick()
	if err != nil {
		t.Error(err)
	}
	// We shouldn't be able to stick again.
	err = hand.Stick()
	if !errors.Is(err, ErrHandLocked) {
		t.Errorf("didn't get appropriate error when attempting re-stick, expected HandLocked got %s", err)
	}

	// And we definitely shouldn't be able to hit again.
	err = hand.Hit()
	if !errors.Is(err, ErrHandLocked) {
		t.Errorf("didn't get appropriate error when attempting hit on locked hand, expected HandLocked got %s", err)
	}
}

// Test that busting the hand out works as intended.
func TestHandIntentionalBust(t *testing.T) {
	table := NewTable(1)
	player := Player{
		Table: table,
	}

	err := table.Join(&player)
	if err != nil {
		t.Error(err)
	}
	hand, err := newHand(&player)
	if err != nil {
		t.Error(err)
	}
	hand.Cards = []playdeck.Card{
		{Suit: playdeck.SuitSpade, Value: playdeck.ValueAce},
		{Suit: playdeck.SuitSpade, Value: playdeck.ValueAce},
	}
	err = hand.EvalScore()
	if err != nil {
		t.Error(err)
	}
	// Force state to play
	table.playState = 1
	t.Logf("score: %v, hand: %v", hand.score, hand.Cards)
	// we should bust within 21 cards, surely
	for i := 0; i < 21; i++ {
		err = hand.Hit()
		if err != nil {
			t.Error(err)
		}
		score, _, _, valid := hand.Score()
		if valid {
			t.Logf("score: %v, hand: %v", score, hand.Cards)
		} else {
			t.Logf("bust! score: %v, hand: %v", score, hand.Cards)
			break
		}
	}
	if hand.valid {
		t.Errorf("the hand is still valid after 21 cards")
	}
	// We shouldn't be able to hit again.
	err = hand.Hit()
	if !errors.Is(err, ErrHandBust) {
		t.Errorf("didn't get appropriate error when attempting re-hit, expected HandBust got %s", err)
	}
}

// Test that the scoring algorithm works as expected.
func TestHandEvalScore(t *testing.T) {
	hand := Hand{}
	// Check that a hand with no Cards isn't valid to score
	err := hand.EvalScore()
	if !errors.Is(err, ErrHandInvalid) {
		t.Errorf("hand with less than two cards got inappropriate response when scoring, expected HandInvalid got %v", err)
	}
	// Make the hand invalid, with an ace (valid) and a joker (definitely not valid)
	hand.Cards = []playdeck.Card{
		{Suit: playdeck.SuitSpade, Value: playdeck.ValueAce},
		{Suit: playdeck.SuitJoker, Value: playdeck.ValueJoker},
	}
	err = hand.EvalScore()
	if !errors.Is(err, ErrInvalidCard) {
		t.Errorf("hand with an invalid card not caught, expected InvalidCard got %v", err)
	}
	// Test that face cards are scored at 10
	hand.Cards = []playdeck.Card{
		{Suit: playdeck.SuitSpade, Value: playdeck.ValueKing},
		{Suit: playdeck.SuitDiamond, Value: playdeck.ValueQueen},
	}
	err = hand.EvalScore()
	if err != nil {
		t.Error(err)
	}
	if hand.score != 20 {
		t.Errorf("hand with face cards should be worth 20, got %v", hand.score)
	}
}
