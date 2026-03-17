package internal

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// WakaTime JSON export structs — only the fields we need.

type wakaTimeExport struct {
	Days []wakaTimeDay `json:"days"`
}

type wakaTimeDay struct {
	Date       string              `json:"date"` // "2025-03-27"
	GrandTotal wakaTimeGrandTotal  `json:"grand_total"`
	Projects   []wakaTimeProject   `json:"projects"`
	Languages  []wakaTimeTimeEntry `json:"languages"`
}

type wakaTimeGrandTotal struct {
	TotalSeconds float64 `json:"total_seconds"`
}

type wakaTimeProject struct {
	Name       string              `json:"name"`
	GrandTotal wakaTimeGrandTotal  `json:"grand_total"`
	Languages  []wakaTimeTimeEntry `json:"languages"`
}

type wakaTimeTimeEntry struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
}

// ParseWakaTimeExport reads a WakaTime JSON export file and returns
// one MonthlySummary per month found in the data.
func ParseWakaTimeExport(filePath string) ([]MonthlySummary, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading wakatime export: %w", err)
	}

	var export wakaTimeExport
	if err := json.Unmarshal(data, &export); err != nil {
		return nil, fmt.Errorf("parsing wakatime JSON: %w", err)
	}

	if len(export.Days) == 0 {
		return nil, fmt.Errorf("no days found in wakatime export")
	}

	// Group days by month (YYYY-MM).
	monthDays := make(map[string][]wakaTimeDay)
	for _, day := range export.Days {
		if len(day.Date) < 7 {
			continue
		}
		month := day.Date[:7] // "2025-03"
		monthDays[month] = append(monthDays[month], day)
	}

	var summaries []MonthlySummary
	for month, days := range monthDays {
		ms := buildSummaryFromWakaTime(month, days)
		summaries = append(summaries, ms)
	}

	// Sort by month ascending.
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Month < summaries[j].Month
	})

	return summaries, nil
}

func buildSummaryFromWakaTime(month string, days []wakaTimeDay) MonthlySummary {
	ms := MonthlySummary{
		Month:            month,
		Projects:         make(map[string]int64),
		Languages:        make(map[string]int64),
		ProjectLanguages: make(map[string]map[string]int64),
	}

	uniqueDays := make(map[string]bool)
	var firstDay, lastDay string

	for _, day := range days {
		if day.GrandTotal.TotalSeconds <= 0 {
			continue
		}

		uniqueDays[day.Date] = true

		if firstDay == "" || day.Date < firstDay {
			firstDay = day.Date
		}
		if lastDay == "" || day.Date > lastDay {
			lastDay = day.Date
		}

		ms.TotalSeconds += int64(math.Round(day.GrandTotal.TotalSeconds))

		// Top-level languages.
		for _, lang := range day.Languages {
			name := strings.ToLower(lang.Name)
			ms.Languages[name] += int64(math.Round(lang.TotalSeconds))
		}

		// Projects and per-project languages.
		for _, proj := range day.Projects {
			ms.Projects[proj.Name] += int64(math.Round(proj.GrandTotal.TotalSeconds))

			if ms.ProjectLanguages[proj.Name] == nil {
				ms.ProjectLanguages[proj.Name] = make(map[string]int64)
			}
			for _, lang := range proj.Languages {
				name := strings.ToLower(lang.Name)
				ms.ProjectLanguages[proj.Name][name] += int64(math.Round(lang.TotalSeconds))
			}
		}
	}

	ms.DaysTracked = len(uniqueDays)
	ms.FirstDay = firstDay
	ms.LastDay = lastDay

	return ms
}
