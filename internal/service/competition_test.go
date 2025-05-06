package service

import (
	"testing"
	"time"

	"github.com/numero_quadro/biathlon-tracker/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewCompetitionService(t *testing.T) {
	config := &domain.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30.000",
	}

	service := NewCompetitionService(config)
	assert.NotNil(t, service)
	assert.Equal(t, config, service.config)
	assert.NotNil(t, service.competitors)
	assert.NotNil(t, service.events)
	assert.NotNil(t, service.log)
}

func TestProcessEvent_Registration(t *testing.T) {
	config := &domain.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30.000",
	}

	service := NewCompetitionService(config)
	eventTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	event := domain.NewEvent(eventTime, domain.EventTypeIncoming, int(domain.EventRegistered), 1, "")

	err := service.ProcessEvent(event)
	assert.NoError(t, err)
	assert.Contains(t, service.competitors, 1)
	assert.Equal(t, domain.StatusRegistered, service.competitors[1].Status)
	assert.Len(t, service.log, 1)
	assert.Contains(t, service.log[0], "The competitor(1) registered")
}

func TestProcessEvent_StartTimeSet(t *testing.T) {
	config := &domain.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30.000",
	}

	service := NewCompetitionService(config)

	// Register competitor first
	registerTime := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	registerEvent := domain.NewEvent(registerTime, domain.EventTypeIncoming, int(domain.EventRegistered), 1, "")
	err := service.ProcessEvent(registerEvent)
	assert.NoError(t, err)

	// Set start time
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	startEvent := domain.NewEvent(startTime, domain.EventTypeIncoming, int(domain.EventStartTimeSet), 1, "10:00:00.000")
	err = service.ProcessEvent(startEvent)
	assert.NoError(t, err)

	expectedTime, _ := time.Parse("15:04:05.000", "10:00:00.000")
	assert.Equal(t, expectedTime, service.competitors[1].PlannedStart)
	assert.Len(t, service.log, 2)
	assert.Contains(t, service.log[1], "The start time for the competitor(1) was set by a draw to 10:00:00.000")
}

func TestProcessEvent_TargetHit(t *testing.T) {
	config := &domain.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30.000",
	}

	service := NewCompetitionService(config)

	// Register competitor
	registerTime := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	registerEvent := domain.NewEvent(registerTime, domain.EventTypeIncoming, int(domain.EventRegistered), 1, "")
	err := service.ProcessEvent(registerEvent)
	assert.NoError(t, err)

	// Hit target
	hitTime := time.Date(2024, 1, 1, 10, 5, 0, 0, time.UTC)
	hitEvent := domain.NewEvent(hitTime, domain.EventTypeIncoming, int(domain.EventTargetHit), 1, "1")
	err = service.ProcessEvent(hitEvent)
	assert.NoError(t, err)

	assert.Equal(t, 1, service.competitors[1].Hits)
	assert.Equal(t, 1, service.competitors[1].Shots)
	assert.Len(t, service.log, 2)
	assert.Contains(t, service.log[1], "The target(1) has been hit by competitor(1)")
}

func TestProcessEvent_Finish(t *testing.T) {
	config := &domain.Config{
		Laps:        2,
		LapLen:      3500,
		PenaltyLen:  150,
		FiringLines: 2,
		Start:       "10:00:00.000",
		StartDelta:  "00:01:30.000",
	}

	service := NewCompetitionService(config)

	// Register competitor
	registerTime := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	registerEvent := domain.NewEvent(registerTime, domain.EventTypeIncoming, int(domain.EventRegistered), 1, "")
	err := service.ProcessEvent(registerEvent)
	assert.NoError(t, err)

	// Set start time
	startTime := time.Date(2024, 1, 1, 10, 0, 0, 0, time.UTC)
	startEvent := domain.NewEvent(startTime, domain.EventTypeIncoming, int(domain.EventStartTimeSet), 1, "10:00:00.000")
	err = service.ProcessEvent(startEvent)
	assert.NoError(t, err)

	// Start race
	raceStartTime := time.Date(2024, 1, 1, 10, 0, 1, 0, time.UTC)
	raceStartEvent := domain.NewEvent(raceStartTime, domain.EventTypeIncoming, int(domain.EventStarted), 1, "")
	err = service.ProcessEvent(raceStartEvent)
	assert.NoError(t, err)

	// Complete first lap
	lap1Time := time.Date(2024, 1, 1, 10, 10, 0, 0, time.UTC)
	lap1Event := domain.NewEvent(lap1Time, domain.EventTypeIncoming, int(domain.EventEndedMainLap), 1, "")
	err = service.ProcessEvent(lap1Event)
	assert.NoError(t, err)

	// Complete second lap (finish)
	lap2Time := time.Date(2024, 1, 1, 10, 20, 0, 0, time.UTC)
	lap2Event := domain.NewEvent(lap2Time, domain.EventTypeIncoming, int(domain.EventEndedMainLap), 1, "")
	err = service.ProcessEvent(lap2Event)
	assert.NoError(t, err)

	assert.Equal(t, domain.StatusFinished, service.competitors[1].Status)
	assert.Len(t, service.competitors[1].Laps, 2)
	assert.Len(t, service.log, 6)
	assert.Contains(t, service.log[5], "The competitor(1) has finished")
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{
			name:     "zero duration",
			duration: 0,
			expected: "00:00:00.000",
		},
		{
			name:     "one hour",
			duration: time.Hour,
			expected: "01:00:00.000",
		},
		{
			name:     "one minute",
			duration: time.Minute,
			expected: "00:01:00.000",
		},
		{
			name:     "one second",
			duration: time.Second,
			expected: "00:00:01.000",
		},
		{
			name:     "one millisecond",
			duration: time.Millisecond,
			expected: "00:00:00.001",
		},
		{
			name:     "complex duration",
			duration: 2*time.Hour + 30*time.Minute + 15*time.Second + 500*time.Millisecond,
			expected: "02:30:15.500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStatusString(t *testing.T) {
	tests := []struct {
		name     string
		status   domain.CompetitorStatus
		expected string
	}{
		{
			name:     "not started",
			status:   domain.StatusNotStarted,
			expected: "NotStarted",
		},
		{
			name:     "not finished",
			status:   domain.StatusNotFinished,
			expected: "NotFinished",
		},
		{
			name:     "disqualified",
			status:   domain.StatusDisqualified,
			expected: "Disqualified",
		},
		{
			name:     "finished",
			status:   domain.StatusFinished,
			expected: "Finished",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getStatusString(tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}
