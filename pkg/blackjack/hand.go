package blackjack

import (
	"github.com/duckfullstop/checkmate/pkg/playdeck"
	"sync"
)

// A Hand is a representation of the state of a player's hand of cards. It stores cards, handles hitting and sticking, and calculates scores.
// Safe for use in asynchronous environments.
type Hand struct {
	sync.RWMutex

	// Cards stores this hand's drawn Cards.
	Cards []playdeck.Card

	// Table is a reference to the Table that this hand is being played on.
	// This is a shortcut to Hand.Player.Table
	// SWEng: discussion point?
	Table *Table

	// Player is a reference to the Player that owns this Hand.
	Player *Player

	// score is the current best possible score of this Hand. This is calculated with aces being worth 10 if it makes sense for them to be, otherwise 1.
	score int
	// minScore is the lowest possible score of this Hand, if all aces are being scored as 1.
	minScore int
	// locked is true if the hand cannot be played further, for example when it's been stuck or is bust
	locked bool
	// valid is true if the hand is not bust.
	valid bool
}

// newHand creates a new Hand object. For internal use.
// Not thread-safe! Use with caution (or better still, with the Player lock)
func newHand(player *Player) (hand Hand, err error) {
	if player == nil {
		return Hand{}, ErrPlayerInvalid
	}
	if player.Table == nil {
		return Hand{}, ErrPlayerNoTable
	}

	return Hand{
		Table:  player.Table,
		Player: player,
		valid:  true,
	}, nil
}

// Score returns the current state of this Hand. Thread-safe.
// score represents the maximum score of the hand, taking into account aces being reduced in value to attempt to avoid a bust (possibly in vain)
// minScore represents the minimum possible score of the hand, taking into account all aces being reduced in value.
// locked is true if the hand cannot be played any further (i.e. it's been stuck or is bust)
// valid is true if the hand is not bust.
func (h *Hand) Score() (score int, minScore int, locked bool, valid bool) {
	h.RLock()
	defer h.RUnlock()
	return h.score, h.minScore, h.locked, h.valid
}

// Hit adds a card to this Hand, if possible. Automatically re-evaluates score, ending play on the hand if Bust occurs.
func (h *Hand) Hit() (err error) {
	err = h.canPlay()
	if err != nil {
		return err
	}

	err = h.addCard()
	if err != nil {
		return err
	}
	err = h.EvalScore()
	return err
}

// Stick ends play on this hand. Locks the hand for further play.
func (h *Hand) Stick() (err error) {
	err = h.lockHand()
	if err != nil {
		return err
	}
	return h.EvalScore()
}

// canPlay returns an error if the hand cannot be played further, otherwise nil.
func (h *Hand) canPlay() (err error) {
	h.RLock()
	defer h.RUnlock()
	// Check to ensure we may proceed with the hit
	if h.Table == nil || h.Table.Deck == nil {
		return ErrHandInvalid
	}
	if !h.valid {
		return ErrHandBust
	}
	if h.locked {
		return ErrHandLocked
	}

	if h.Player == nil {
		return ErrHandNoPlayer
	}
	h.Player.RLock()
	defer h.Player.RUnlock()
	if h.Player.Table == nil {
		return ErrPlayerNoTable
	}
	// SwEng: discussion point, should this be handled elsewhere / not as a direct check?
	if h.Player.Table.playState != 1 {
		return ErrTableNotInPlay
	}
	return
}

// addCard is an internal function for adding a card to the hand. Prefer Hit().
func (h *Hand) addCard() (err error) {
	h.Lock()
	defer h.Unlock()

	newCard, err := h.Table.Deck.PullRandomCard()
	if err != nil {
		return err
	}

	h.Cards = append(h.Cards, newCard)
	return nil
}

// lockHand is an internal function that locks the hand from being played further. Prefer Stick().
func (h *Hand) lockHand() (err error) {
	h.Lock()
	defer h.Unlock()

	if h.locked {
		return ErrHandLocked
	}
	h.locked = true
	return nil
}

// EvalScore forcefully evaluates the score of the given hand, storing the current maximum and minimum score in the object.
// This is called automatically by all actions (e.g. Hit() and Stick(), etc.), and is left public for use in integration tests.
// Prefer Score() for thread-safe access of the score, hand validity, and lock status.
func (h *Hand) EvalScore() (err error) {
	h.Lock()
	defer h.Unlock()

	// Start scoring from 0.
	h.score = 0
	h.minScore = 0

	// A valid hand of cards must always have at least 2 cards.
	if len(h.Cards) < 2 {
		return ErrHandInvalid
	}

	// First iteration loop is just for simple cards and basic sanity checking of the hand.
	for _, c := range h.Cards {
		// Catch invalid cards in hand. (This should never happen?)
		if c.Value == playdeck.ValueJoker || !c.Valid() {
			return ErrInvalidCard
		}

		// SWEng: Using a switch here isn't useful, as we'd have to use conditionals on each case
		// It might be a little more readable, but IMO using if-then here is just as reasonable

		// Skip aces early on.
		if c.Value == 1 {
			continue
		}
		// Normal scored cards, do these first.
		if c.Value >= 2 && c.Value <= 10 {
			// Just add the score, nothing more to do.
			h.score += c.Value.Value()
		}
		if c.Value >= 11 && c.Value <= 13 {
			// It's a face card, they're worth 10.
			h.score += 10
		}
	}

	// The minimum score after the simple cards is the same as the full score.
	h.minScore = h.score

	// Second iteration is for aces.
	for _, c := range h.Cards {
		if c.Value == 1 {
			// It's an ace
			// The minimum value is always 1, so let's get that out of the way first
			h.minScore++
			// Now determine how to score it.
			if h.score <= 10 {
				// It wouldn't bust us, treat as 11.
				h.score += 11
			} else {
				// It would bust us, treat as 1.
				h.score++
			}
		}
	}

	// Now perform the bust check.
	if h.score > 21 {
		// We're bust! D:
		// The hand is invalid, and immediately locked out of play.
		// SWEng: This is done separately from Hand.lockHand() as we already have the mutex lock
		h.valid = false
		h.locked = true
	}

	return
}
