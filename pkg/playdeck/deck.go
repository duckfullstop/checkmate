package playdeck

import (
	"math/rand"
	"sync"
	"time"
)

// A Deck is a helper struct that contains a quantity of Cards.
// Safe for use in asynchronous environments, so long as you don't directly mutate its state.
type Deck struct {
	sync.Mutex
	Cards *[]Card
}

// NewDeck returns a memory pointer to a new, standard, 52-card Deck.
// This is a shortcut to NewDeckOfDecks(1, joker)
func NewDeck(joker bool) (deck *Deck) {
	// SWEng: Instead of NewDeckOfDecks calling this method multiple times, we do it the other way around
	// this way we don't have to do a tonne of expensive, potentially quite large slice merges later on
	// (effectively this function becomes a shortcut)
	return NewDeckOfDecks(1, joker)
}

// NewDeckOfDecks returns a memory pointer to a Deck containing a specified number of new, standard, 52-card Decks.
func NewDeckOfDecks(count int, joker bool) *Deck {
	// Initialize new deck
	deck := Deck{Cards: new([]Card)}
	for i := 0; i < count; i++ {
		// For each suit...
		for s := CardSuit(1); s <= 4; s++ {
			// and each card in a suit...
			for v := CardValue(1); v <= 13; v++ {
				// append a new card
				*deck.Cards = append(*deck.Cards, Card{Suit: s, Value: v})
			}
		}
		// Add the Joker too, if it's requested.
		if joker {
			*deck.Cards = append(*deck.Cards, Card{Suit: 0, Value: 0})
		}
	}

	return &deck
}

// PullRandomCard returns a random card from the Deck, if possible.
// It returns an error if this is not possible for some reason (i.e the deck is empty).
func (d *Deck) PullRandomCard() (card Card, err error) {
	// Take a mutex lock, as this operation mutates the state of the deck
	d.Lock()
	defer d.Unlock()

	// Check that the Deck is valid to have a card drawn
	if d.Cards == nil {
		return card, ErrDeckUninitialized
	}
	if len(*d.Cards) == 0 {
		return card, ErrDeckEmpty
	}

	// Randomness could be seeded in other ways (discussion point?)
	// lint call out use of weak RNG, suggest crypto/rand (but is that necessary?)
	rand.Seed(time.Now().Unix()) // initialize PRNG
	indexToPull := rand.Intn(len(*d.Cards))
	card = (*d.Cards)[indexToPull]
	// Useful one-liner for deletion: https://github.com/golang/go/wiki/SliceTricks#delete
	*d.Cards = append((*d.Cards)[:indexToPull], (*d.Cards)[indexToPull+1:]...)

	return card, nil
}

// PullCard returns the first card on top of the Deck, if possible.
// It returns an error if this is not possible for some reason (i.e the deck is empty).
func (d *Deck) PullCard() (card Card, err error) {
	// Take a mutex lock, as this operation mutates the state of the deck
	d.Lock()
	defer d.Unlock()

	// Check that the Deck is valid to have a card drawn
	if d.Cards == nil {
		return card, ErrDeckUninitialized
	}
	if len(*d.Cards) == 0 {
		return card, ErrDeckEmpty
	}

	card = (*d.Cards)[0]

	*d.Cards = append((*d.Cards)[:0], (*d.Cards)[1:]...)

	return card, nil
}

// PushCard inserts the card back into the bottom of the Deck.
// It returns an error if this is not possible for some reason (i.e the deck is uninitialized)
func (d *Deck) PushCard(card Card) (err error) {
	// Take a mutex lock, as this operation mutates the state of the deck
	d.Lock()
	defer d.Unlock()

	// Check that the Deck is valid to have a card inserted
	if d.Cards == nil {
		return ErrDeckUninitialized
	}

	*d.Cards = append(*d.Cards, card)

	return
}

// Shuffle randomly repositions all cards in the Deck. This may be useful if your game depends on having a linear deck chronology.
// You might want to consider PullRandomCard() if you only need the deck to be pseudo-random.
// It returns an error if this is not possible for some reason (i.e the deck is uninitialized)
func (d *Deck) Shuffle() (err error) {
	// Take a mutex lock, as this operation mutates the state of the deck
	d.Lock()
	defer d.Unlock()

	// Check that the Deck is valid to have a card drawn
	if d.Cards == nil {
		return ErrDeckUninitialized
	}
	if len(*d.Cards) == 0 {
		return ErrDeckEmpty
	}
	// Randomness could be seeded in other ways (discussion point?)
	rand.Seed(time.Now().Unix()) // initialize PRNG
	rand.Shuffle(len(*d.Cards), func(i, j int) {
		(*d.Cards)[i], (*d.Cards)[j] = (*d.Cards)[j], (*d.Cards)[i]
	})
	return
}
