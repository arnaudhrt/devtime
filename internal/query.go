package internal

import (
	"path/filepath"
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

// readRawEventSummaries reads all remaining (non-compacted) event files,
// computes sessions, and returns a slice of per-file sessions.
func readRawEventSessions() ([][]Session, error) {
	dir, err := EventsDir()
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "events-*.jsonl"))
	if err != nil {
		return nil, err
	}

	var allSessions [][]Session
	for _, path := range matches {
		events, err := ReadEvents(path)
		if err != nil {
			return nil, err
		}
		if len(events) > 0 {
			allSessions = append(allSessions, ComputeSessions(events))
		}
	}
	return allSessions, nil
}

// AllTimeSummary returns a combined Summary from all compacted summaries
// plus any remaining raw event files.
func AllTimeSummary() (Summary, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return Summary{}, err
	}

	combined := Summary{}
	for _, ms := range summaries {
		combined = MergeSummary(combined, SummaryFromMonthly(ms))
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return Summary{}, err
	}
	for _, sessions := range rawSessions {
		combined = MergeSummary(combined, Summarize(sessions))
	}

	return combined, nil
}

// AllTimeSummaryForProject returns a combined Summary filtered for a
// specific project, across compacted summaries and raw events.
func AllTimeSummaryForProject(project string) (Summary, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return Summary{}, err
	}

	combined := Summary{}
	for _, ms := range summaries {
		combined = MergeSummary(combined, SummaryFromMonthlyForProject(ms, project))
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return Summary{}, err
	}
	for _, sessions := range rawSessions {
		filtered := FilterByProject(sessions, project)
		combined = MergeSummary(combined, Summarize(filtered))
	}

	return combined, nil
}

// AllTimeSummaryForLanguage returns a combined Summary filtered for a
// specific language, across compacted summaries and raw events.
func AllTimeSummaryForLanguage(lang string) (Summary, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return Summary{}, err
	}

	combined := Summary{}
	for _, ms := range summaries {
		combined = MergeSummary(combined, SummaryFromMonthlyForLanguage(ms, lang))
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return Summary{}, err
	}
	for _, sessions := range rawSessions {
		filtered := FilterByLanguage(sessions, lang)
		combined = MergeSummary(combined, Summarize(filtered))
	}

	return combined, nil
}

// AllTimeProjectNames returns all unique project names across compacted
// summaries and raw event files.
func AllTimeProjectNames() ([]string, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	for _, ms := range summaries {
		for name := range ms.Projects {
			seen[name] = true
		}
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return nil, err
	}
	for _, sessions := range rawSessions {
		for _, s := range sessions {
			seen[s.Project] = true
		}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// AllTimeLanguageNames returns all unique language names across compacted
// summaries and raw event files.
func AllTimeLanguageNames() ([]string, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	for _, ms := range summaries {
		for name := range ms.Languages {
			seen[name] = true
		}
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return nil, err
	}
	for _, sessions := range rawSessions {
		for _, s := range sessions {
			seen[s.Language] = true
		}
	}

	names := make([]string, 0, len(seen))
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names, nil
}

// ProfileData holds the data needed by the profile command.
type ProfileData struct {
	Summary     Summary
	DaysTracked int
	FirstDay    time.Time
}

// LoadProfileData returns combined profile data from compacted summaries
// and raw event files.
func LoadProfileData() (ProfileData, error) {
	summaries, err := ReadAllMonthlySummaries()
	if err != nil {
		return ProfileData{}, err
	}

	combined := Summary{}
	totalDays := 0
	var firstDay time.Time

	for _, ms := range summaries {
		combined = MergeSummary(combined, SummaryFromMonthly(ms))
		totalDays += ms.DaysTracked
		if ms.FirstDay != "" {
			d, err := time.Parse(time.DateOnly, ms.FirstDay)
			if err == nil && (firstDay.IsZero() || d.Before(firstDay)) {
				firstDay = d
			}
		}
	}

	rawSessions, err := readRawEventSessions()
	if err != nil {
		return ProfileData{}, err
	}

	days := make(map[string]bool)
	for _, sessions := range rawSessions {
		combined = MergeSummary(combined, Summarize(sessions))
		for _, s := range sessions {
			days[s.Start.Format(time.DateOnly)] = true
			if firstDay.IsZero() || s.Start.Before(firstDay) {
				firstDay = s.Start
			}
		}
	}
	totalDays += len(days)

	return ProfileData{
		Summary:     combined,
		DaysTracked: totalDays,
		FirstDay:    firstDay,
	}, nil
}
