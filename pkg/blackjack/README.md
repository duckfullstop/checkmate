# Blackjack

The _blackjack_ package exposes the core of this project, that being a Blackjack / Pontoon logic and scoring
solution.

## Why is this a package / what's here?

This package contains the core Blackjack logic, including concepts of _tables_, _players_, and _hands_.

### Table
A _Table_ stores the state of a Blackjack game. It has its own Deck of Cards (see `playdeck` package) to draw from, instead of just drawing them from thin air -
this way Players can only receive cards that are legitimately in the deck (no duplicates if you're only playing with one deck!).

_Tables_ have one or more _Players_ associated with them.

Tables can be in a few states, and these can be extended on later if more functionality is desired - pre-game, in-game, and post-game.

### Player
A _Player_ is the representation of a person (or bot!) that would be sat at a physical Table.
It must be associated with a _Table_ to work properly.

_Players_ may have multiple _hands_, though start with none - they are given a new one with two cards at the start of each new Game.
This implementation currently doesn't allow for splitting, but can easily be added by simply adding a new _Hand_ to the _Player_.

### Hand
A _Hand_ is, quite simply, a player's Hand of cards. Hands can be hit or stuck / stood (`hand.Hit()` and `hand.Stick()` respectively).
After each operation on a Hand, its score is re-evaluated, and either frozen out of play (if stick is called, or if the hand is bust),
or left open for further play.

## Thoughts on Implementation
This current implementation only provides for Hit and Stick, though the other decisions can be implemented easily as follows:

* Double Down: Bets are not currently implemented, though this would be as simple as adding a _Wager_ value to the _Hand_ then doubling it, doing `hand.Hit()` and then `hand.Stick()` to finalise.
* Split: Create a new _Hand_ associated with the same _Player_ via the `hand.Player` pointer, then move one _Card_ in the current _Hand_ to the new one. Immediately call `hand.Hit()` on both hands.
  * Game logic already handles this behaviour by checking that all _Hands_ belonging to all _Players_ are locked out of play.
* Surrender: Immediately lock the hand and render it invalid. Return half of the theoretical _Wager_ value on the hand to the player.

The concept of a _dealer_ is currently lacking, though could easily be added by using the reserved table game state `2`, and assigning a special _Player_ with a special _Dealer_ flag.

At present, dealing a new game discards all cards in the previous deck and starts again with new 52-card deck(s) from scratch, pulling random cards from the new deck to simulate a shuffle.
If a more authentic game allowing for advantage play (e.g. card-counting) is desired, then cards from all Hands would simply be reintroduced back into the Deck (using `deck.Push(card)`).

Everything in this package SHOULD be safe to call in goroutines due to the healthy usage of mutex locking throughout, though this isn't presently tested.

Laying the package out in this table-player-hand style allows for easy extension into things like multiplayer (discussion point, perhaps?), as well as adding bot support (simply by holding a _Player_ instance that is manipulated by a bot package).
