package blackjack

import "errors"

// Errors throwable by this module.
var (
	ErrHandNoPlayer             = errors.New("hand is not associated with a player")
	ErrHandLocked               = errors.New("hand is locked")
	ErrHandNotLocked            = errors.New("hand is not locked")
	ErrHandInvalid              = errors.New("hand is not correctly instantiated")
	ErrHandBust                 = errors.New("hand is bust")
	ErrInvalidCard              = errors.New("card in hand is invalid")
	ErrPlayerNoTable            = errors.New("player has no table assigned")
	ErrPlayerInvalid            = errors.New("player is invalid")
	ErrTableInPlay              = errors.New("table is in play")
	ErrTableNotInPlay           = errors.New("table is not in play")
	ErrTablePlayerAlreadyJoined = errors.New("player already on table")
)
