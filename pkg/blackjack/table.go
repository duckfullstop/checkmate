package blackjack

import (
	"github.com/duckfullstop/checkmate/pkg/playdeck"
	"sync"
)

// A Table represents an object that houses the overall state of a game of Blackjack.
// Tables have a Deck of Cards to deal from, and a set of Players (with many Hands) playing at it.
type Table struct {
	sync.Mutex
	// The deck of cards to play from. Refilled every new game based on sDecks.
	Deck *playdeck.Deck
	// The setting for the number of decks to refill the Deck with every new game.
	sDecks int

	// Pointers to all the Players currently playing on this Table.
	Players []*Player

	// The current state of play of the table. 0 = not in play, 1 = play in progress, 2 = (reserved, dealer playing), 3 = endgame (payouts, etc)
	playState uint
}

// NewTable initializes a new Table for further use.
func NewTable(decks int) (table *Table) {
	table = new(Table)
	table.Deck = playdeck.NewDeckOfDecks(decks, false)
	table.sDecks = decks
	return table
}

// Join seats the given Player at the Table.
// This function cannot be used and will return an error if the table is either in play, or if the Player is already known to this Table.
func (t *Table) Join(p *Player) (err error) {
	// Lock the table
	t.Lock()
	defer t.Unlock()

	// Lock the player
	p.Lock()
	defer p.Unlock()

	// Can't join a table that's in progress
	// SwEng: this could be changed fairly easily, but it's a safety check for this implementation
	if t.playState != 0 {
		return ErrTableInPlay
	}

	for _, tp := range t.Players {
		if p == tp {
			return ErrTablePlayerAlreadyJoined
		}
	}
	t.Players = append(t.Players, p)
	p.Table = t
	return
}

// Reset sets the game state of a table back to 0 (pre-game).
// It revokes all hands that each Player has.
// It can only be called successfully if the game is not in play (SWEng: this could be changed).
func (t *Table) Reset() (err error) {
	t.Lock()
	defer t.Unlock()

	// SWEng: Two ways to write this, the other one is (if (not valid state) and (not valid state))
	// Both are equally tricky reads, so a comment is probably a good idea, such as:
	// Check if the table is in a valid state to reset, throwing if it's in play or otherwise invalid
	if !(t.playState == 3 || t.playState == 0) {
		return ErrTableInPlay
	}
	for _, p := range t.Players {
		// Clear their hands
		p.clearHands()
	}

	t.playState = 0
	return
}

// Deal starts the game by dealing 2 cards from a new Deck into a new hand for each Player.
// This function can only be used if the game is not in play (gameState 0 or 3).
// SWEng: This is a function I'd honestly like to completely reengineer because returning a slice of errors is silly
func (t *Table) Deal() (err []error) {
	// First ensure the game state is clean
	e := t.Reset()
	if e != nil {
		// kinda silly way of returning a single error, but it works
		// SWEng: does it? Good discussion point around perhaps putting an error state on each Player object
		return append(err, e)
	}

	// Now we can take the lock and do stuff ourselves
	t.Lock()
	defer t.Unlock()

	// Throw out the current deck, and create a new one.
	// Yes, this is the equivalent of just throwing an entire pack of cards into the shredder and pulling a new one out of the box,
	// but it works for pseudo-randomness.
	// See README.md for further discussion.
	t.Deck = playdeck.NewDeckOfDecks(t.sDecks, false)

	// This would be pretty straight forward to switch to a goroutine for speed,
	// and to also handle errors better.
	for _, p := range t.Players {
		// Create a new Hand for the player
		h, e := p.newHand()
		if e != nil {
			err = append(err, e)
			continue
		}
		// Pull two cards
		for i := 0; i < 2; i++ {
			e := h.addCard()
			if e != nil {
				err = append(err, e)
				continue
			}
		}
		// Evaluate this hand's score, so it's ready to go when the player looks at their cards
		e = h.EvalScore()
		if e != nil {
			err = append(err, e)
			continue
		}
	}
	t.playState = 1
	return
}

// EndRound moves the game to the end phase.
// It can only be used if the game is in the Play state, and no Players have any Hands that are not locked.
func (t *Table) EndRound() (err error) {
	t.Lock()
	defer t.Unlock()
	if t.playState != 1 {
		// kinda silly way of returning a single error, but it works
		return ErrTableNotInPlay
	}
	// SWEng: This can easily be moved to a goroutine. It's a simple check that all players have finished playing with their hands.
	for _, p := range t.Players {
		for _, h := range p.Hands {
			if !h.locked {
				return ErrHandNotLocked
			}
		}
	}

	// move to endgame phase
	t.playState = 3
	return
}
