package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	Timestamp time.Time
	Type      string // "heartbeat", "focus", "blur"
	Project   string
	Language  string
	Editor    string
}

type rawEvent struct {
	Ts      string `json:"ts"`
	Event   string `json:"event"`
	Project string `json:"project"`
	Lang    string `json:"lang"`
	Editor  string `json:"editor"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var raw rawEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC3339, raw.Ts)
	if err != nil {
		return err
	}
	e.Timestamp = t
	e.Type = raw.Event
	e.Project = raw.Project
	e.Language = raw.Lang
	e.Editor = raw.Editor
	return nil
}

// EventsDir returns ~/.devtime/
func EventsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".devtime"), nil
}

// EventFilePath returns the path for a given year/month.
// e.g. ~/.devtime/events-2026-03.jsonl
func EventFilePath(dir string, year int, month time.Month) string {
	return filepath.Join(dir, fmt.Sprintf("events-%04d-%02d.jsonl", year, int(month)))
}

// ReadEvents reads all events from a single JSONL file.
// Skips malformed lines silently. Returns empty slice if file doesn't exist.
func ReadEvents(path string) ([]Event, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer f.Close()

	var events []Event
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // skip malformed lines
		}
		events = append(events, ev)
	}
	return events, scanner.Err()
}

// ReadEventsForRange reads events from all month files that overlap
// with the given time range [start, end].
func ReadEventsForRange(start, end time.Time) ([]Event, error) {
	dir, err := EventsDir()
	if err != nil {
		return nil, err
	}

	var allEvents []Event

	// Iterate month by month from start to end
	cursor := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, start.Location())
	endMonth := time.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, end.Location())

	for !cursor.After(endMonth) {
		path := EventFilePath(dir, cursor.Year(), cursor.Month())
		events, err := ReadEvents(path)
		if err != nil {
			return nil, err
		}
		for _, ev := range events {
			if !ev.Timestamp.Before(start) && !ev.Timestamp.After(end) {
				allEvents = append(allEvents, ev)
			}
		}
		cursor = cursor.AddDate(0, 1, 0)
	}

	return allEvents, nil
}

// ReadAllEvents reads events from all event files in ~/.devtime/
func ReadAllEvents() ([]Event, error) {
	dir, err := EventsDir()
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "events-*.jsonl"))
	if err != nil {
		return nil, err
	}

	var allEvents []Event
	for _, path := range matches {
		events, err := ReadEvents(path)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, events...)
	}

	return allEvents, nil
}
