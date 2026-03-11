package internal

import (
	"sort"
	"time"
)

type ProjectSummary struct {
	Name     string
	Duration time.Duration
}

type LanguageSummary struct {
	Name     string
	Duration time.Duration
}

type Summary struct {
	Total     time.Duration
	Projects  []ProjectSummary  // sorted by duration descending
	Languages []LanguageSummary // sorted by duration descending
}

// Summarize aggregates a slice of sessions into a Summary.
func Summarize(sessions []Session) Summary {
	var total time.Duration
	projects := make(map[string]time.Duration)
	languages := make(map[string]time.Duration)

	for _, s := range sessions {
		total += s.Duration
		projects[s.Project] += s.Duration
		languages[s.Language] += s.Duration
	}

	var ps []ProjectSummary
	for name, dur := range projects {
		ps = append(ps, ProjectSummary{Name: name, Duration: dur})
	}
	sort.Slice(ps, func(i, j int) bool {
		return ps[i].Duration > ps[j].Duration
	})

	var ls []LanguageSummary
	for name, dur := range languages {
		ls = append(ls, LanguageSummary{Name: name, Duration: dur})
	}
	sort.Slice(ls, func(i, j int) bool {
		return ls[i].Duration > ls[j].Duration
	})

	return Summary{
		Total:     total,
		Projects:  ps,
		Languages: ls,
	}
}

// FilterByProject returns only sessions matching the given project name.
func FilterByProject(sessions []Session, project string) []Session {
	var filtered []Session
	for _, s := range sessions {
		if s.Project == project {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// FilterByLanguage returns only sessions matching the given language.
func FilterByLanguage(sessions []Session, lang string) []Session {
	var filtered []Session
	for _, s := range sessions {
		if s.Language == lang {
			filtered = append(filtered, s)
		}
	}
	return filtered
}

// TodayRange returns start/end of today (local time).
func TodayRange() (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return start, now
}

// WeekRange returns start of Monday / end of today (local time).
func WeekRange() (time.Time, time.Time) {
	now := time.Now()
	weekday := now.Weekday()
	// Go's Sunday=0, Monday=1, ...
	// We want Monday as start of week
	daysSinceMonday := int(weekday+6) % 7
	start := time.Date(now.Year(), now.Month(), now.Day()-daysSinceMonday, 0, 0, 0, 0, now.Location())
	return start, now
}

// MonthRange returns start of month / end of today (local time).
func MonthRange() (time.Time, time.Time) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return start, now
}
