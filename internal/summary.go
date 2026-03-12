package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// MonthlySummary holds pre-aggregated data for a compacted month.
type MonthlySummary struct {
	Month            string                      `json:"month"`
	TotalSeconds     int64                       `json:"total_seconds"`
	DaysTracked      int                         `json:"days_tracked"`
	FirstDay         string                      `json:"first_day"`
	LastDay          string                      `json:"last_day"`
	Projects         map[string]int64            `json:"projects"`
	Languages        map[string]int64            `json:"languages"`
	ProjectLanguages map[string]map[string]int64 `json:"project_languages"`
}

// BuildMonthlySummary creates a MonthlySummary from sessions.
func BuildMonthlySummary(month string, sessions []Session) MonthlySummary {
	ms := MonthlySummary{
		Month:            month,
		Projects:         make(map[string]int64),
		Languages:        make(map[string]int64),
		ProjectLanguages: make(map[string]map[string]int64),
	}

	if len(sessions) == 0 {
		return ms
	}

	days := make(map[string]bool)
	var earliest, latest time.Time

	for _, s := range sessions {
		secs := int64(math.Round(s.Duration.Seconds()))
		ms.TotalSeconds += secs
		ms.Projects[s.Project] += secs
		ms.Languages[s.Language] += secs

		if ms.ProjectLanguages[s.Project] == nil {
			ms.ProjectLanguages[s.Project] = make(map[string]int64)
		}
		ms.ProjectLanguages[s.Project][s.Language] += secs

		day := s.Start.Format(time.DateOnly)
		days[day] = true

		if earliest.IsZero() || s.Start.Before(earliest) {
			earliest = s.Start
		}
		if latest.IsZero() || s.Start.After(latest) {
			latest = s.Start
		}
	}

	ms.DaysTracked = len(days)
	ms.FirstDay = earliest.Format(time.DateOnly)
	ms.LastDay = latest.Format(time.DateOnly)

	return ms
}

// SummaryFilePath returns the path for a monthly summary file.
func SummaryFilePath(dir, month string) string {
	return filepath.Join(dir, fmt.Sprintf("summary-%s.json", month))
}

// WriteSummary writes a MonthlySummary to a JSON file.
func WriteSummary(dir string, ms MonthlySummary) error {
	data, err := json.MarshalIndent(ms, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SummaryFilePath(dir, ms.Month), data, 0644)
}

// ReadMonthlySummary reads a single summary JSON file.
func ReadMonthlySummary(path string) (MonthlySummary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return MonthlySummary{}, err
	}
	var ms MonthlySummary
	if err := json.Unmarshal(data, &ms); err != nil {
		return MonthlySummary{}, err
	}
	return ms, nil
}

// ReadAllMonthlySummaries reads all summary-*.json files from ~/.devtime/.
func ReadAllMonthlySummaries() ([]MonthlySummary, error) {
	dir, err := EventsDir()
	if err != nil {
		return nil, err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "summary-*.json"))
	if err != nil {
		return nil, err
	}

	var summaries []MonthlySummary
	for _, path := range matches {
		ms, err := ReadMonthlySummary(path)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, ms)
	}
	return summaries, nil
}

// CompactMonth reads a month's raw events, computes sessions, writes a
// summary file, and deletes the raw events file.
func CompactMonth(dir string, year int, month time.Month) error {
	eventsPath := EventFilePath(dir, year, month)
	events, err := ReadEvents(eventsPath)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		// No events — just remove the empty file
		return os.Remove(eventsPath)
	}

	sessions := ComputeSessions(events)
	monthStr := fmt.Sprintf("%04d-%02d", year, int(month))
	ms := BuildMonthlySummary(monthStr, sessions)

	if err := WriteSummary(dir, ms); err != nil {
		return err
	}

	return os.Remove(eventsPath)
}

// AutoCompact compacts event files for months that ended more than 7 days ago.
func AutoCompact() error {
	dir, err := EventsDir()
	if err != nil {
		return err
	}

	matches, err := filepath.Glob(filepath.Join(dir, "events-*.jsonl"))
	if err != nil {
		return err
	}

	now := time.Now()

	for _, path := range matches {
		base := filepath.Base(path)
		var year, month int
		if _, err := fmt.Sscanf(base, "events-%d-%d.jsonl", &year, &month); err != nil {
			continue // skip files that don't match the naming pattern
		}

		// First day of the month after this one
		nextMonth := time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, now.Location())

		// Compact if the month ended more than 7 days ago
		if now.Sub(nextMonth) >= 7*24*time.Hour {
			if err := CompactMonth(dir, year, time.Month(month)); err != nil {
				return fmt.Errorf("compacting %s: %w", base, err)
			}
		}
	}

	return nil
}

// SummaryFromMonthly converts a MonthlySummary into a Summary.
func SummaryFromMonthly(ms MonthlySummary) Summary {
	total := time.Duration(ms.TotalSeconds) * time.Second

	ps := make([]ProjectSummary, 0, len(ms.Projects))
	for name, secs := range ms.Projects {
		ps = append(ps, ProjectSummary{Name: name, Duration: time.Duration(secs) * time.Second})
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Duration > ps[j].Duration })

	ls := make([]LanguageSummary, 0, len(ms.Languages))
	for name, secs := range ms.Languages {
		ls = append(ls, LanguageSummary{Name: name, Duration: time.Duration(secs) * time.Second})
	}
	sort.Slice(ls, func(i, j int) bool { return ls[i].Duration > ls[j].Duration })

	return Summary{Total: total, Projects: ps, Languages: ls}
}

// SummaryFromMonthlyForProject extracts a project-specific Summary
// from a MonthlySummary using the project_languages cross-product.
func SummaryFromMonthlyForProject(ms MonthlySummary, project string) Summary {
	projectSecs, ok := ms.Projects[project]
	if !ok {
		return Summary{}
	}

	total := time.Duration(projectSecs) * time.Second

	var ls []LanguageSummary
	if pl, ok := ms.ProjectLanguages[project]; ok {
		for lang, secs := range pl {
			ls = append(ls, LanguageSummary{Name: lang, Duration: time.Duration(secs) * time.Second})
		}
		sort.Slice(ls, func(i, j int) bool { return ls[i].Duration > ls[j].Duration })
	}

	return Summary{
		Total:     total,
		Projects:  []ProjectSummary{{Name: project, Duration: total}},
		Languages: ls,
	}
}

// SummaryFromMonthlyForLanguage extracts a language-specific Summary
// from a MonthlySummary using the project_languages cross-product.
func SummaryFromMonthlyForLanguage(ms MonthlySummary, lang string) Summary {
	langSecs, ok := ms.Languages[lang]
	if !ok {
		return Summary{}
	}

	total := time.Duration(langSecs) * time.Second

	var ps []ProjectSummary
	for project, langs := range ms.ProjectLanguages {
		if secs, ok := langs[lang]; ok {
			ps = append(ps, ProjectSummary{Name: project, Duration: time.Duration(secs) * time.Second})
		}
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Duration > ps[j].Duration })

	return Summary{
		Total:     total,
		Projects:  ps,
		Languages: []LanguageSummary{{Name: lang, Duration: total}},
	}
}

// MergeSummary combines two Summary objects by summing durations.
func MergeSummary(a, b Summary) Summary {
	total := a.Total + b.Total

	pm := make(map[string]time.Duration)
	for _, p := range a.Projects {
		pm[p.Name] += p.Duration
	}
	for _, p := range b.Projects {
		pm[p.Name] += p.Duration
	}
	ps := make([]ProjectSummary, 0, len(pm))
	for name, dur := range pm {
		ps = append(ps, ProjectSummary{Name: name, Duration: dur})
	}
	sort.Slice(ps, func(i, j int) bool { return ps[i].Duration > ps[j].Duration })

	lm := make(map[string]time.Duration)
	for _, l := range a.Languages {
		lm[l.Name] += l.Duration
	}
	for _, l := range b.Languages {
		lm[l.Name] += l.Duration
	}
	ls := make([]LanguageSummary, 0, len(lm))
	for name, dur := range lm {
		ls = append(ls, LanguageSummary{Name: name, Duration: dur})
	}
	sort.Slice(ls, func(i, j int) bool { return ls[i].Duration > ls[j].Duration })

	return Summary{Total: total, Projects: ps, Languages: ls}
}
