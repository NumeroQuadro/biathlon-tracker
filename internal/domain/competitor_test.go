package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCompetitor(t *testing.T) {
	competitor := NewCompetitor(1)
	assert.NotNil(t, competitor)
	assert.Equal(t, 1, competitor.ID)
	assert.Equal(t, StatusRegistered, competitor.Status)
	assert.NotNil(t, competitor.Laps)
	assert.NotNil(t, competitor.Penalties)
	assert.Equal(t, 0, competitor.CurrentLap)
	assert.Equal(t, 0, competitor.Hits)
	assert.Equal(t, 0, competitor.Shots)
}

func TestAddLap(t *testing.T) {
	competitor := NewCompetitor(1)
	lapTime := 10 * time.Minute
	speed := 5.0

	competitor.AddLap(lapTime, speed)
	assert.Len(t, competitor.Laps, 1)
	assert.Equal(t, 1, competitor.CurrentLap)
	assert.Equal(t, lapTime, competitor.TotalTime)

	// Add second lap
	competitor.AddLap(lapTime, speed)
	assert.Len(t, competitor.Laps, 2)
	assert.Equal(t, 2, competitor.CurrentLap)
	assert.Equal(t, 2*lapTime, competitor.TotalTime)
}

func TestAddPenalty(t *testing.T) {
	competitor := NewCompetitor(1)
	penaltyTime := 1 * time.Minute
	speed := 2.5

	competitor.AddPenalty(penaltyTime, speed)
	assert.Len(t, competitor.Penalties, 1)
	assert.Equal(t, penaltyTime, competitor.TotalTime)

	// Add second penalty
	competitor.AddPenalty(penaltyTime, speed)
	assert.Len(t, competitor.Penalties, 2)
	assert.Equal(t, 2*penaltyTime, competitor.TotalTime)
}

func TestRecordShot(t *testing.T) {
	competitor := NewCompetitor(1)

	// Record a hit
	competitor.RecordShot(true)
	assert.Equal(t, 1, competitor.Hits)
	assert.Equal(t, 1, competitor.Shots)

	// Record a miss
	competitor.RecordShot(false)
	assert.Equal(t, 1, competitor.Hits)
	assert.Equal(t, 2, competitor.Shots)
}
