package domain

type State int

const (
	Unregistered State = iota
	Registered
	Scheduled
	LapStarted
	Started
	FiringRangeEntered
	FiringRangeLeft
	PenaltyLapEntered
	PenaltyLapLeft
	LapEnded
)
