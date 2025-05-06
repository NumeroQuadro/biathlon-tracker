package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/numero_quadro/biathlon-tracker/internal/domain"
	"github.com/numero_quadro/biathlon-tracker/internal/service"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: biathlon-tracker <config_file> <events_file>")
		os.Exit(1)
	}

	config, err := loadConfig(os.Args[1])
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	competition := service.NewCompetitionService(config)
	events, err := loadEvents(os.Args[2])
	if err != nil {
		fmt.Printf("Error loading events: %v\n", err)
		os.Exit(1)
	}

	for _, event := range events {
		if err := competition.ProcessEvent(event); err != nil {
			fmt.Printf("Error processing event: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println(competition.GetFinalReport())
}

func loadConfig(path string) (*domain.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config domain.Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %v", err)
	}

	return &config, nil
}

func loadEvents(path string) ([]*domain.Event, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening events file: %v", err)
	}
	defer file.Close()

	var events []*domain.Event
	scanner := bufio.NewScanner(file)
	eventRegex := regexp.MustCompile(`\[(\d{2}:\d{2}:\d{2}\.\d{3})\] (\d+) (\d+)(?: (.+))?`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := eventRegex.FindStringSubmatch(line)
		if len(matches) < 4 {
			continue
		}

		time, err := time.Parse("15:04:05.000", matches[1])
		if err != nil {
			return nil, fmt.Errorf("error parsing time: %v", err)
		}

		eventID, err := strconv.Atoi(matches[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing event ID: %v", err)
		}

		competitorID, err := strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("error parsing competitor ID: %v", err)
		}

		extraParams := ""
		if len(matches) > 4 {
			extraParams = matches[4]
		}

		events = append(events, domain.NewEvent(time, domain.EventTypeIncoming, eventID, competitorID, extraParams))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading events file: %v", err)
	}

	return events, nil
}
