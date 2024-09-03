package domain

import "errors"

var (
	ErrRequired           = errors.New("required value")
	ErrNotFound           = errors.New("not found")
	ErrNil                = errors.New("nil data")
	ErrNegative           = errors.New("negative value")
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidUserLogin   = errors.New("invalid user login")
	ErrInvalidProviderIDs = errors.New("invalid provider IDs")
	ErrNoUserInContext    = errors.New("no user in context")
	ErrInvalidValue       = errors.New("invalid value")
	ErrInvalidFormatEmail = errors.New("must be email type")
	ErrInvalidLength      = errors.New("invalid length")
)
