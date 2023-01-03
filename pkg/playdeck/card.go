package playdeck

import "fmt"

// CardValue represents the value of a Card.
// For example, the **NINE** of Spades.
type CardValue uint8

// A Card's Value is one of these tokens.
//
//goland:noinspection GoUnusedConst
const (
	ValueJoker CardValue = iota
	ValueAce
	ValueTwo
	ValueThree
	ValueFour
	ValueFive
	ValueSix
	ValueSeven
	ValueEight
	ValueNine
	ValueTen
	ValueJack
	ValueQueen
	ValueKing
)

var cardNames = map[CardValue]string{
	0:  "joker",
	1:  "ace",
	2:  "2",
	3:  "3",
	4:  "4",
	5:  "5",
	6:  "6",
	7:  "7",
	8:  "8",
	9:  "9",
	10: "10",
	11: "jack",
	12: "queen",
	13: "king",
}

func (v CardValue) String() string {
	valueName, exists := cardNames[v]
	if !exists {
		return "unknown"
	}
	// Note: optionally capitalise first letter using something like strings.ToTitle() (or just modify the string directly)
	return valueName
}

// ValueAceHigh is a helper function that shifts Aces to slot 14, for use in games with direct card comparison that treat the Ace as being high.
// Not used in Blackjack.
func (v CardValue) ValueAceHigh() (value int) {
	if v == ValueAce {
		return 14
	}
	return v.Value()
}

func (v CardValue) Value() int {
	return int(v)
}

// CardSuit represents the suit of a Card.
// For example, the Ten of **HEARTS**.
type CardSuit uint8

// A Card's Suit is one of these tokens.
const (
	SuitJoker CardSuit = iota
	SuitClub
	SuitDiamond
	SuitHeart
	SuitSpade
)

var suitNames = map[CardSuit]string{
	0: "joker",
	1: "club",
	2: "diamond",
	3: "heart",
	4: "spade",
}

// String returns the non-pluralised string representation for the name of this suit (e.g. "heart").
// Unknown or invalid suits return "unknown".
func (s CardSuit) String() (name string) {
	suitName, exists := suitNames[s]
	if !exists {
		return "unknown"
	}
	// Note: optionally capitalise first letter using something like strings.ToTitle() (or just modify the string directly)
	return suitName
}

// Value returns the integer value for this suit.
// You probably don't need to use this, compare directly against playdeck.Suit*.
func (s CardSuit) Value() int {
	return int(s)
}

// A Card represents a single playing card.
type Card struct {
	// SWEng: We choose to store suit and value here as integers instead of strings to make interpolation easier,
	// and also to cut down on memory usage when a large number of cards are in play.

	// The Suit of the card is an integer from 1 to 4, where 1 is "club", 2 is "diamond", 3 is "heart", and 4 is "spade". 0 is reserved for the Joker.
	Suit CardSuit

	// The Value of the card is an integer from 1 to 12, where 1 is an "ace" and 13 is a "king". 0 is reserved for the Joker.
	Value CardValue
}

// String returns a human-readable representation of the card (e.g "Nine of Diamonds").
func (c *Card) String() (name string) {
	// Special case for the Joker.
	// (We catch both cases here just to account for any edge cases)
	if c.Suit == SuitJoker || c.Value == ValueJoker {
		return "joker"
	}
	// Feel like handling pluralisation could be done better here?
	return fmt.Sprintf("%s of %ss", c.Value.String(), c.Suit.String())
}

// Valid returns whether this Card struct represents a valid playing card of any type.
func (c *Card) Valid() (valid bool) {
	// By virtue of being unsigned, the value is guaranteed to be > 0.
	return c.Value <= 13 && c.Suit <= 4
}
