package internal

import (
	"fmt"
	"strings"
	"time"
)

const BarWidth = 20
const MinDisplayDuration = 3 * time.Minute

// FilterShort removes projects and languages with less than MinDisplayDuration
// from a Summary. The Total is not affected.
func FilterShort(s Summary) Summary {
	var projects []ProjectSummary
	for _, p := range s.Projects {
		if p.Duration >= MinDisplayDuration {
			projects = append(projects, p)
		}
	}
	var languages []LanguageSummary
	for _, l := range s.Languages {
		if l.Duration >= MinDisplayDuration {
			languages = append(languages, l)
		}
	}
	return Summary{
		Total:     s.Total,
		Projects:  projects,
		Languages: languages,
	}
}

// FormatDuration formats a duration as "Xh Ym" (e.g. "2h 45m").
func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %02dm", h, m)
}

// PrintSummary prints a full summary with bar chart.
// Header example: "Today: 4h 23m"
func PrintSummary(header string, summary Summary) {
	if summary.Total == 0 {
		fmt.Printf("%s: no data\n", header)
		return
	}

	summary = FilterShort(summary)

	fmt.Printf("\n  %s: %s\n\n", header, FormatDuration(summary.Total))

	// Find the longest project name for alignment
	maxNameLen := 0
	for _, p := range summary.Projects {
		if len(p.Name) > maxNameLen {
			maxNameLen = len(p.Name)
		}
	}

	// Print project breakdown with bar chart
	fmt.Println("  Projects:")
	for _, p := range summary.Projects {
		pct := float64(p.Duration) / float64(summary.Total) * 100
		filled := int(float64(BarWidth) * float64(p.Duration) / float64(summary.Total))
		if filled < 0 {
			filled = 0
		}
		if filled > BarWidth {
			filled = BarWidth
		}
		empty := BarWidth - filled
		bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
		fmt.Printf("    %-*s  %s  %s  %3.0f%%\n", maxNameLen, p.Name, FormatDuration(p.Duration), bar, pct)
	}

	// Print language breakdown
	if len(summary.Languages) > 0 {
		maxLangLen := 0
		for _, l := range summary.Languages {
			if len(l.Name) > maxLangLen {
				maxLangLen = len(l.Name)
			}
		}

		fmt.Println()
		fmt.Println("  Languages:")
		for _, l := range summary.Languages {
			pct := float64(l.Duration) / float64(summary.Total) * 100
			filled := int(float64(BarWidth) * float64(l.Duration) / float64(summary.Total))
			if filled < 0 {
				filled = 0
			}
			if filled > BarWidth {
				filled = BarWidth
			}
			empty := BarWidth - filled
			bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
			fmt.Printf("    %-*s  %s  %s  %3.0f%%\n", maxLangLen, l.Name, FormatDuration(l.Duration), bar, pct)
		}
	}

	fmt.Println()
}

// PrintStatus prints the active/inactive status.
func PrintStatus(active bool, session *Session) {
	if active && session != nil {
		dur := time.Since(session.Start)
		fmt.Printf("\n")
		fmt.Printf("  Status: active\n")
		fmt.Printf("  Project:  %s\n", session.Project)
		fmt.Printf("  Language: %s\n", session.Language)
		fmt.Printf("  Editor:   %s\n", session.Editor)
		fmt.Printf("  Session:  %s\n\n", FormatDuration(dur))
	} else if session != nil {
		fmt.Printf("\n")
		fmt.Printf("  Status: not active\n")
		fmt.Printf("  Last session: %s on %s (%s)\n\n",
			FormatDuration(session.Duration),
			session.Project,
			session.End.Format("15:04"))
	} else {
		fmt.Printf("\n")
		fmt.Printf(" Status: not active\n")
		fmt.Printf(" No recent sessions\n")
	}
}
