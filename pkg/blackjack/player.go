package blackjack

import "sync"

// A Player is a representation of someone sat at a Table, possibly with a nice cup of tea or similar beverage.
type Player struct {
	sync.RWMutex
	// A Player belongs to a Table.
	Table *Table

	// A Player has one (or possibly two) Hands.
	Hands []*Hand
}

// NewPlayer returns a new Player instance.
// SWEng: This is a separate function to just being bound to a Table to support future concepts e.g. Players moving between tables mid-life
func NewPlayer() (player *Player) {
	return new(Player)
}

// clearHands empties out this Player's hands.
// All hands must be unlocked for write, or this will stall.
// SWEng: Should this stall? Should we return an error if we can't get the lock?
func (p *Player) clearHands() {
	p.Lock()
	defer p.Unlock()
	for _, h := range p.Hands {
		// Ensure write lock is available (no need to worry about releasing, we'll be destroying it shortly)
		h.Lock()
	}
	// Create empty hands table (garbage collector should take care of the orphans)
	p.Hands = []*Hand{}
	return
}

// newHand instantiates a new Hand on the given Player.
func (p *Player) newHand() (hand *Hand, err error) {
	p.Lock()
	defer p.Unlock()
	h, err := newHand(p)
	p.Hands = append(p.Hands, &h)
	return &h, err
}
