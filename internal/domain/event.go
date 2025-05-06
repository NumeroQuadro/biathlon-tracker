package domain

import "time"

// EventType represents the type of event (incoming or outgoing)
type EventType int

const (
	EventTypeIncoming EventType = iota
	EventTypeOutgoing
)

// Event represents a competition event
type Event struct {
	Time         time.Time
	Type         EventType
	EventID      int
	CompetitorID int
	ExtraParams  string
}

type IncomingEventID int

const (
	EventRegistered IncomingEventID = iota + 1
	EventStartTimeSet
	EventOnStartLine
	EventStarted
	EventOnFiringRange
	EventTargetHit
	EventLeftFiringRange
	EventEnteredPenaltyLaps
	EventLeftPenaltyLaps
	EventEndedMainLap
	EventCannotContinue
)

type OutgoingEventID int

const (
	EventDisqualified OutgoingEventID = iota + 32
	EventFinished
)

func NewEvent(time time.Time, eventType EventType, eventID int, competitorID int, extraParams string) *Event {
	return &Event{
		Time:         time,
		Type:         eventType,
		EventID:      eventID,
		CompetitorID: competitorID,
		ExtraParams:  extraParams,
	}
}
