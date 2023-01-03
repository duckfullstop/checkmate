package playdeck

import "testing"

func TestInvalidCard(t *testing.T) {
	// Bad suit and value
	c := Card{Suit: 50, Value: 42}
	if c.Suit.String() != "unknown" {
		t.Errorf("Bad card returned a suit of %s", c.Suit.String())
	}
	if c.Value.String() != "unknown" {
		t.Errorf("Bad card returned a value of %s", c.Value.String())
	}
	if c.Valid() {
		t.Errorf("Bad card is apparently valid")
	}
}

func TestValidCard(t *testing.T) {
	c := Card{Suit: SuitSpade, Value: ValueAce}
	if c.Suit.String() != "spade" {
		t.Errorf("Card returned invalid suit of %s", c.Suit.String())
	}
	if c.Suit.Value() != 4 {
		t.Errorf("Spade returned the incorrect suit value of %v", c.Value.Value())
	}
	if c.Value.String() != "ace" {
		t.Errorf("Card returned invalid value of %s", c.Value.String())
	}
	if c.String() != "ace of spades" {
		t.Errorf("Card returned invalid string representation %s", c.String())
	}

	if c.Value.Value() != 1 {
		t.Errorf("Ace returned the incorrect value of %v", c.Value.Value())
	}
	if c.Value.ValueAceHigh() != 14 {
		t.Errorf("Ace returned the incorrect Aces High value of %v", c.Value.ValueAceHigh())
	}

	if !c.Valid() {
		t.Errorf("Valid card is apparently invalid")
	}

	// Joker!
	c = Card{Suit: SuitJoker, Value: ValueJoker}
	if c.String() != "joker" {
		t.Errorf("Joker returned invalid string representation %s", c.String())
	}
	if c.Value.ValueAceHigh() != 0 {
		t.Errorf("Joker returned the incorrect Aces High value of %v", c.Value.ValueAceHigh())
	}

	if !c.Valid() {
		t.Errorf("Valid joker is apparently invalid")
	}
}
