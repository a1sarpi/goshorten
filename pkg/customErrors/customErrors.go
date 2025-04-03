package customErrors

import "errors"

var (
	ErrURLNotFound          = errors.New("URL not found")
	ErrInvalidURL           = errors.New("Invalid URL")
	ErrClicksUninitialized  = errors.New("Uninitialized Clicks map")
	ErrExpiresUninitialized = errors.New("Uninitialized Expires map")
	ErrURLExpired           = errors.New("Attempt to use URL that expired")
)
