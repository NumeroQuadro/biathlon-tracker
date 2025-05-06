package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig_GetStartTime(t *testing.T) {
	tests := []struct {
		name        string
		startTime   string
		expectError bool
	}{
		{
			name:        "valid time",
			startTime:   "10:00:00.000",
			expectError: false,
		},
		{
			name:        "invalid format",
			startTime:   "10:00:00",
			expectError: true,
		},
		{
			name:        "invalid time",
			startTime:   "25:00:00.000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{Start: tt.startTime}
			parsedTime, err := config.GetStartTime()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				expected, _ := time.Parse("15:04:05.000", tt.startTime)
				assert.Equal(t, expected, parsedTime)
			}
		})
	}
}

func TestConfig_GetStartDelta(t *testing.T) {
	tests := []struct {
		name        string
		startDelta  string
		expectError bool
	}{
		{
			name:        "valid delta",
			startDelta:  "00:01:30.000",
			expectError: false,
		},
		{
			name:        "invalid format",
			startDelta:  "00:01:30",
			expectError: true,
		},
		{
			name:        "invalid time",
			startDelta:  "00:61:30.000",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &Config{StartDelta: tt.startDelta}
			delta, err := config.GetStartDelta()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				expected, _ := time.Parse("15:04:05.000", tt.startDelta)
				assert.Equal(t, expected, delta)
			}
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectError bool
	}{
		{
			name: "valid config",
			config: &Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       "10:00:00.000",
				StartDelta:  "00:01:30.000",
			},
			expectError: false,
		},
		{
			name: "invalid laps",
			config: &Config{
				Laps:        0,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       "10:00:00.000",
				StartDelta:  "00:01:30.000",
			},
			expectError: true,
		},
		{
			name: "invalid lap length",
			config: &Config{
				Laps:        2,
				LapLen:      0,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       "10:00:00.000",
				StartDelta:  "00:01:30.000",
			},
			expectError: true,
		},
		{
			name: "invalid penalty length",
			config: &Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  0,
				FiringLines: 2,
				Start:       "10:00:00.000",
				StartDelta:  "00:01:30.000",
			},
			expectError: true,
		},
		{
			name: "invalid firing lines",
			config: &Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 0,
				Start:       "10:00:00.000",
				StartDelta:  "00:01:30.000",
			},
			expectError: true,
		},
		{
			name: "invalid start time",
			config: &Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       "invalid",
				StartDelta:  "00:01:30.000",
			},
			expectError: true,
		},
		{
			name: "invalid start delta",
			config: &Config{
				Laps:        2,
				LapLen:      3500,
				PenaltyLen:  150,
				FiringLines: 2,
				Start:       "10:00:00.000",
				StartDelta:  "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
