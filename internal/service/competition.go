package service

import (
	"fmt"
	"time"

	"github.com/numero_quadro/biathlon-tracker/internal/domain"
)

// CompetitionService handles the biathlon competition logic
type CompetitionService struct {
	config      *domain.Config
	competitors map[int]*domain.Competitor
	events      []*domain.Event
	log         []string
}

// NewCompetitionService creates a new competition service
func NewCompetitionService(config *domain.Config) *CompetitionService {
	return &CompetitionService{
		config:      config,
		competitors: make(map[int]*domain.Competitor),
		events:      make([]*domain.Event, 0),
		log:         make([]string, 0),
	}
}

func (s *CompetitionService) formatEventMessage(event *domain.Event) string {
	timeStr := event.Time.Format("15:04:05.000")
	switch domain.IncomingEventID(event.EventID) {
	case domain.EventRegistered:
		return fmt.Sprintf("[%s] The competitor(%d) registered", timeStr, event.CompetitorID)
	case domain.EventStartTimeSet:
		return fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s", timeStr, event.CompetitorID, event.ExtraParams)
	case domain.EventOnStartLine:
		return fmt.Sprintf("[%s] The competitor(%d) is on the start line", timeStr, event.CompetitorID)
	case domain.EventStarted:
		return fmt.Sprintf("[%s] The competitor(%d) has started", timeStr, event.CompetitorID)
	case domain.EventOnFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) is on the firing range(%s)", timeStr, event.CompetitorID, event.ExtraParams)
	case domain.EventTargetHit:
		return fmt.Sprintf("[%s] The target(%s) has been hit by competitor(%d)", timeStr, event.ExtraParams, event.CompetitorID)
	case domain.EventLeftFiringRange:
		return fmt.Sprintf("[%s] The competitor(%d) left the firing range", timeStr, event.CompetitorID)
	case domain.EventEnteredPenaltyLaps:
		return fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps", timeStr, event.CompetitorID)
	case domain.EventLeftPenaltyLaps:
		return fmt.Sprintf("[%s] The competitor(%d) left the penalty laps", timeStr, event.CompetitorID)
	case domain.EventEndedMainLap:
		return fmt.Sprintf("[%s] The competitor(%d) ended the main lap", timeStr, event.CompetitorID)
	case domain.EventCannotContinue:
		return fmt.Sprintf("[%s] The competitor(%d) can't continue: %s", timeStr, event.CompetitorID, event.ExtraParams)
	}
	return ""
}

// ProcessEvent processes an incoming event
func (s *CompetitionService) ProcessEvent(event *domain.Event) error {
	if event.Type != domain.EventTypeIncoming {
		return fmt.Errorf("invalid event type: %v", event.Type)
	}

	competitor, exists := s.competitors[event.CompetitorID]
	if !exists && event.EventID != int(domain.EventRegistered) {
		return fmt.Errorf("competitor %d not registered", event.CompetitorID)
	}

	// Log the event
	if msg := s.formatEventMessage(event); msg != "" {
		s.log = append(s.log, msg)
	}

	switch domain.IncomingEventID(event.EventID) {
	case domain.EventRegistered:
		s.competitors[event.CompetitorID] = domain.NewCompetitor(event.CompetitorID)
	case domain.EventStartTimeSet:
		competitor.PlannedStart = parseTime(event.ExtraParams)
	case domain.EventOnStartLine:
		competitor.Status = domain.StatusOnStartLine
	case domain.EventStarted:
		competitor.Status = domain.StatusRacing
		competitor.StartTime = event.Time
	case domain.EventOnFiringRange:
		competitor.Status = domain.StatusOnFiringRange
	case domain.EventTargetHit:
		competitor.RecordShot(true)
	case domain.EventLeftFiringRange:
		competitor.Status = domain.StatusRacing
	case domain.EventEnteredPenaltyLaps:
		competitor.Status = domain.StatusOnPenaltyLaps
		competitor.TotalTime = event.Time.Sub(competitor.StartTime)
	case domain.EventLeftPenaltyLaps:
		competitor.Status = domain.StatusRacing
		penaltyTime := event.Time.Sub(competitor.StartTime) - competitor.TotalTime
		speed := float64(s.config.PenaltyLen) / penaltyTime.Seconds()
		competitor.AddPenalty(penaltyTime, speed)
	case domain.EventEndedMainLap:
		competitor.Status = domain.StatusRacing
		if len(competitor.Laps) > 0 {
			lastLapEnd := competitor.TotalTime
			lapTime := event.Time.Sub(competitor.StartTime) - lastLapEnd
			speed := float64(s.config.LapLen) / lapTime.Seconds()
			competitor.AddLap(lapTime, speed)
		} else {
			lapTime := event.Time.Sub(competitor.StartTime)
			speed := float64(s.config.LapLen) / lapTime.Seconds()
			competitor.AddLap(lapTime, speed)
		}
		if competitor.CurrentLap == s.config.Laps {
			competitor.Status = domain.StatusFinished
			finishEvent := domain.NewEvent(event.Time, domain.EventTypeOutgoing, int(domain.EventFinished), event.CompetitorID, "")
			s.events = append(s.events, finishEvent)
			s.log = append(s.log, fmt.Sprintf("[%s] The competitor(%d) has finished", event.Time.Format("15:04:05.000"), event.CompetitorID))
		}
	case domain.EventCannotContinue:
		competitor.Status = domain.StatusNotFinished
		competitor.Comment = event.ExtraParams
		disqualifyEvent := domain.NewEvent(event.Time, domain.EventTypeOutgoing, int(domain.EventDisqualified), event.CompetitorID, event.ExtraParams)
		s.events = append(s.events, disqualifyEvent)
	}

	s.events = append(s.events, event)
	return nil
}

// GetEventLog returns the formatted event log
func (s *CompetitionService) GetEventLog() string {
	log := ""
	for _, entry := range s.log {
		log += entry + "\n"
	}
	return log
}

// GetFinalReport generates the final report for all competitors
func (s *CompetitionService) GetFinalReport() string {
	report := ""
	for _, competitor := range s.competitors {
		status := getStatusString(competitor.Status)
		laps := formatLaps(competitor.Laps)
		penalties := formatPenalties(competitor.Penalties)
		shots := fmt.Sprintf("%d/%d", competitor.Hits, competitor.Shots)
		report += fmt.Sprintf("[%s] %d %s %s %s\n", status, competitor.ID, laps, penalties, shots)
	}
	return report
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse("15:04:05.000", timeStr)
	return t
}

func getStatusString(status domain.CompetitorStatus) string {
	switch status {
	case domain.StatusNotStarted:
		return "NotStarted"
	case domain.StatusNotFinished:
		return "NotFinished"
	case domain.StatusDisqualified:
		return "Disqualified"
	default:
		return "Finished"
	}
}

func formatDuration(d time.Duration) string {
	totalSeconds := int(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60
	ms := d.Milliseconds() % 1000
	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, seconds, ms)
}

func formatLaps(laps []domain.LapInfo) string {
	if len(laps) == 0 {
		return "[]"
	}
	result := "["
	for i, lap := range laps {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("{%s, %.3f}", formatDuration(lap.Time), lap.Speed)
	}
	result += "]"
	return result
}

func formatPenalties(penalties []domain.PenaltyInfo) string {
	if len(penalties) == 0 {
		return "{}"
	}
	result := "{"
	for i, penalty := range penalties {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("{%s, %.3f}", formatDuration(penalty.Time), penalty.Speed)
	}
	result += "}"
	return result
}
