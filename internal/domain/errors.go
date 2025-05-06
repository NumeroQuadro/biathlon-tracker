package domain

import "errors"

var (
	ErrInvalidLaps        = errors.New("invalid number of laps")
	ErrInvalidLapLen      = errors.New("invalid lap length")
	ErrInvalidPenaltyLen  = errors.New("invalid penalty length")
	ErrInvalidFiringLines = errors.New("invalid number of firing lines")
)
