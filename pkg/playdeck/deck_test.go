package playdeck

import (
	"errors"
	"testing"
)

func TestNewSingleDeck(t *testing.T) {
	deck := NewDeck(false)
	if len(*deck.Cards) != 52 {
		t.Errorf("deck of cards is not correct size! expected 52 cards, got %v", len(*deck.Cards))
	}

	deck = NewDeck(true)
	if len(*deck.Cards) != 53 {
		t.Errorf("deck of cards is not correct size! expected 53 cards, got %v", len(*deck.Cards))
	}
}

func TestDeckShuffle(t *testing.T) {
	deck := NewDeck(false)

	// Testing that the deck has indeed been shuffled is technically possible, but there's an (incredibly small!) chance
	// that the deck gets shuffled back into its original order (1 in 8.06581752 x 10⁶⁷), causing a failure
	// This seems like a great way to cause CI to mysteriously fail on this test, so it's omitted
	err := deck.Shuffle()
	if err != nil {
		t.Error(err)
	}

	// Check error conditions
	deck.Cards = new([]Card)
	err = deck.Shuffle()
	if !errors.Is(err, ErrDeckEmpty) {
		t.Errorf("wrong error when deck empty, got %s", err)
	}

	deck.Cards = nil
	err = deck.Shuffle()
	if !errors.Is(err, ErrDeckUninitialized) {
		t.Errorf("wrong error when deck empty, got %s", err)
	}
}

func TestDeckMutation(t *testing.T) {
	deck := NewDeck(false)

	if len(*deck.Cards) != 52 {
		t.Errorf("deck of cards is not correct size! expected 52 cards, got %v", len(*deck.Cards))
	}

	_, err := deck.PullRandomCard()
	if err != nil {
		t.Error(err)
	}

	if len(*deck.Cards) != 51 {
		t.Errorf("deck of cards was not reduced in size! expected 51 cards, got %v", len(*deck.Cards))
	}

	card, err := deck.PullCard()
	if err != nil {
		t.Error(err)
	}

	if len(*deck.Cards) != 50 {
		t.Errorf("deck of cards was not reduced in size! expected 50 cards, got %v", len(*deck.Cards))
	}

	err = deck.PushCard(card)

	if err != nil {
		t.Error(err)
	}

	if len(*deck.Cards) != 51 {
		t.Errorf("deck of cards did not increase in size! expected 51 cards, got %v", len(*deck.Cards))
	}
}

func TestDeckMutationErrorHandling(t *testing.T) {
	deck := new(Deck)
	_, err := deck.PullRandomCard()
	if !errors.Is(err, ErrDeckUninitialized) {
		t.Errorf("wrong error when deck uninitialized, got %s", err)
	}
	_, err = deck.PullCard()
	if !errors.Is(err, ErrDeckUninitialized) {
		t.Errorf("wrong error when deck uninitialized, got %s", err)
	}
	err = deck.PushCard(Card{})
	if !errors.Is(err, ErrDeckUninitialized) {
		t.Errorf("wrong error when deck uninitialized, got %s", err)
	}

	// Now add an empty hand and check that pulling cards fails correctly
	deck.Cards = new([]Card)

	_, err = deck.PullRandomCard()
	if !errors.Is(err, ErrDeckEmpty) {
		t.Errorf("wrong error when deck empty, got %s", err)
	}
	_, err = deck.PullCard()
	if !errors.Is(err, ErrDeckEmpty) {
		t.Errorf("wrong error when deck empty, got %s", err)
	}
}
