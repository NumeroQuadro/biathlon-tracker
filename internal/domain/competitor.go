package domain

import "time"

// CompetitorStatus represents the current status of a competitor
type CompetitorStatus int

const (
	StatusRegistered CompetitorStatus = iota
	StatusOnStartLine
	StatusRacing
	StatusOnFiringRange
	StatusOnPenaltyLaps
	StatusFinished
	StatusDisqualified
	StatusNotStarted
	StatusNotFinished
)

// LapInfo represents information about a completed lap
type LapInfo struct {
	Time     time.Duration
	Speed    float64 // m/s
	LapIndex int
}

// PenaltyInfo represents information about penalty laps
type PenaltyInfo struct {
	Time  time.Duration
	Speed float64 // m/s
}

// Competitor represents a biathlon competitor
type Competitor struct {
	ID            int
	Status        CompetitorStatus
	StartTime     time.Time
	PlannedStart  time.Time
	Laps          []LapInfo
	Penalties     []PenaltyInfo
	CurrentLap    int
	Hits          int
	Shots         int
	Comment       string
	DisqualReason string
	TotalTime     time.Duration
}

// NewCompetitor creates a new competitor
func NewCompetitor(id int) *Competitor {
	return &Competitor{
		ID:         id,
		Status:     StatusRegistered,
		Laps:       make([]LapInfo, 0),
		Penalties:  make([]PenaltyInfo, 0),
		CurrentLap: 0,
	}
}

// AddLap adds a completed lap to the competitor's record
func (c *Competitor) AddLap(lapTime time.Duration, speed float64) {
	c.Laps = append(c.Laps, LapInfo{
		Time:     lapTime,
		Speed:    speed,
		LapIndex: len(c.Laps),
	})
	c.CurrentLap++
	c.TotalTime += lapTime
}

// AddPenalty adds penalty lap information
func (c *Competitor) AddPenalty(penaltyTime time.Duration, speed float64) {
	c.Penalties = append(c.Penalties, PenaltyInfo{
		Time:  penaltyTime,
		Speed: speed,
	})
	c.TotalTime += penaltyTime
}

// RecordShot records a shot attempt
func (c *Competitor) RecordShot(hit bool) {
	c.Shots++
	if hit {
		c.Hits++
	}
}
