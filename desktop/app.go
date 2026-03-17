package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/arnaudhrt/devtime/internal"
)

var Version = "dev"

type App struct {
	ctx     context.Context
	Version string
}

func NewApp(version string) *App {
	return &App{Version: version}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// --- Data types returned to frontend ---

type ItemData struct {
	Name     string  `json:"name"`
	Duration string  `json:"duration"`
	Percent  float64 `json:"percent"`
}

type SummaryData struct {
	Total     string     `json:"total"`
	Projects  []ItemData `json:"projects"`
	Languages []ItemData `json:"languages"`
}

type StatusData struct {
	Active   bool   `json:"active"`
	Project  string `json:"project"`
	Language string `json:"language"`
	Editor   string `json:"editor"`
	Session  string `json:"session"`
	LastEnd  string `json:"lastEnd"`
}

type ProfileResponse struct {
	TrackingSince string     `json:"trackingSince"`
	TotalTime     string     `json:"totalTime"`
	DailyAverage  string     `json:"dailyAverage"`
	DaysTracked   int        `json:"daysTracked"`
	TopProjects   []ItemData `json:"topProjects"`
	TopLanguages  []ItemData `json:"topLanguages"`
}

type DetailResponse struct {
	Name      string     `json:"name"`
	AllTime   string     `json:"allTime"`
	ThisMonth string     `json:"thisMonth"`
	ThisWeek  string     `json:"thisWeek"`
	Items     []ItemData `json:"items"`
}

// --- Helper ---

func summaryToData(s internal.Summary) SummaryData {
	s = internal.FilterShort(s)
	data := SummaryData{
		Total:     internal.FormatDuration(s.Total),
		Projects:  make([]ItemData, 0, len(s.Projects)),
		Languages: make([]ItemData, 0, len(s.Languages)),
	}
	for _, p := range s.Projects {
		pct := float64(0)
		if s.Total > 0 {
			pct = float64(p.Duration) / float64(s.Total) * 100
		}
		data.Projects = append(data.Projects, ItemData{
			Name:     p.Name,
			Duration: internal.FormatDuration(p.Duration),
			Percent:  pct,
		})
	}
	for _, l := range s.Languages {
		pct := float64(0)
		if s.Total > 0 {
			pct = float64(l.Duration) / float64(s.Total) * 100
		}
		data.Languages = append(data.Languages, ItemData{
			Name:     l.Name,
			Duration: internal.FormatDuration(l.Duration),
			Percent:  pct,
		})
	}
	return data
}

func rangeToSummary(start, end time.Time) (SummaryData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return SummaryData{Total: "0h 00m"}, nil
	}
	events, err := internal.ReadEventsForRange(start, end)
	if err != nil {
		return SummaryData{}, err
	}
	sessions := internal.ComputeSessions(events)
	summary := internal.Summarize(sessions)
	return summaryToData(summary), nil
}

// --- Bound methods ---

func (a *App) GetToday() (SummaryData, error) {
	start, end := internal.TodayRange()
	return rangeToSummary(start, end)
}

func (a *App) GetWeek() (SummaryData, error) {
	start, end := internal.WeekRange()
	return rangeToSummary(start, end)
}

func (a *App) GetMonth() (SummaryData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return SummaryData{Total: "0h 00m"}, nil
	}

	now := time.Now()
	year := now.Year()
	month := now.Month()

	dir, err := internal.EventsDir()
	if err != nil {
		return SummaryData{}, err
	}

	combined := internal.Summary{}

	// Read compacted summary if it exists
	monthStr := fmt.Sprintf("%04d-%02d", year, int(month))
	summaryPath := filepath.Join(dir, fmt.Sprintf("summary-%s.json", monthStr))
	if _, err := os.Stat(summaryPath); err == nil {
		ms, err := internal.ReadMonthlySummary(summaryPath)
		if err != nil {
			return SummaryData{}, err
		}
		combined = internal.MergeSummary(combined, internal.SummaryFromMonthly(ms))
	}

	// Read raw events
	start := time.Date(year, month, 1, 0, 0, 0, 0, now.Location())
	events, err := internal.ReadEventsForRange(start, now)
	if err != nil {
		return SummaryData{}, err
	}
	if len(events) > 0 {
		sessions := internal.ComputeSessions(events)
		combined = internal.MergeSummary(combined, internal.Summarize(sessions))
	}

	return summaryToData(combined), nil
}

func (a *App) GetProfile() (ProfileResponse, error) {
	if err := internal.CheckDataExists(); err != nil {
		return ProfileResponse{}, nil
	}

	data, err := internal.LoadProfileData()
	if err != nil {
		return ProfileResponse{}, err
	}

	if data.Summary.Total == 0 {
		return ProfileResponse{}, nil
	}

	dailyAvg := data.Summary.Total / time.Duration(data.DaysTracked)

	resp := ProfileResponse{
		TrackingSince: data.FirstDay.Format("Jan 2, 2006"),
		TotalTime:     internal.FormatDuration(data.Summary.Total),
		DailyAverage:  internal.FormatDuration(dailyAvg),
		DaysTracked:   data.DaysTracked,
		TopProjects:   make([]ItemData, 0),
		TopLanguages:  make([]ItemData, 0),
	}

	filtered := internal.FilterShort(data.Summary)

	projects := filtered.Projects
	if len(projects) > 5 {
		projects = projects[:5]
	}
	for _, p := range projects {
		resp.TopProjects = append(resp.TopProjects, ItemData{
			Name:     p.Name,
			Duration: internal.FormatDuration(p.Duration),
		})
	}

	languages := filtered.Languages
	if len(languages) > 5 {
		languages = languages[:5]
	}
	for _, l := range languages {
		resp.TopLanguages = append(resp.TopLanguages, ItemData{
			Name:     l.Name,
			Duration: internal.FormatDuration(l.Duration),
		})
	}

	return resp, nil
}

func (a *App) GetStatus() (StatusData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return StatusData{}, nil
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	events, err := internal.ReadEventsForRange(start, now)
	if err != nil {
		return StatusData{}, err
	}

	sessions := internal.ComputeSessions(events)
	if len(sessions) == 0 {
		return StatusData{Active: false}, nil
	}

	last := sessions[len(sessions)-1]
	active := time.Since(last.End) < 5*time.Minute

	data := StatusData{
		Active:   active,
		Project:  last.Project,
		Language: last.Language,
		Editor:   last.Editor,
		LastEnd:  last.End.Format("15:04"),
	}

	if active {
		data.Session = internal.FormatDuration(time.Since(last.Start))
	} else {
		data.Session = internal.FormatDuration(last.Duration)
	}

	return data, nil
}

func (a *App) GetVersion() string {
	return a.Version
}

func (a *App) GetYear(year int) (SummaryData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return SummaryData{Total: "0h 00m"}, nil
	}

	now := time.Now()
	var lastMonth time.Month
	var end time.Time

	if year == now.Year() {
		lastMonth = now.Month()
		end = now
	} else {
		lastMonth = time.December
		end = time.Date(year, time.December, 31, 23, 59, 59, 999999999, now.Location())
	}

	dir, err := internal.EventsDir()
	if err != nil {
		return SummaryData{}, err
	}

	combined := internal.Summary{}

	for m := time.January; m <= lastMonth; m++ {
		monthStr := fmt.Sprintf("%04d-%02d", year, int(m))
		summaryPath := filepath.Join(dir, fmt.Sprintf("summary-%s.json", monthStr))

		if _, err := os.Stat(summaryPath); err == nil {
			ms, err := internal.ReadMonthlySummary(summaryPath)
			if err != nil {
				return SummaryData{}, err
			}
			combined = internal.MergeSummary(combined, internal.SummaryFromMonthly(ms))
		}

		start := time.Date(year, m, 1, 0, 0, 0, 0, now.Location())
		monthEnd := start.AddDate(0, 1, 0).Add(-time.Nanosecond)
		if m == lastMonth && year == now.Year() {
			monthEnd = end
		}
		events, err := internal.ReadEventsForRange(start, monthEnd)
		if err != nil {
			return SummaryData{}, err
		}
		if len(events) > 0 {
			sessions := internal.ComputeSessions(events)
			combined = internal.MergeSummary(combined, internal.Summarize(sessions))
		}
	}

	return summaryToData(combined), nil
}

func (a *App) GetProjectNames() ([]ItemData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return nil, nil
	}

	summary, err := internal.AllTimeSummary()
	if err != nil {
		return nil, err
	}

	summary = internal.FilterShort(summary)
	items := make([]ItemData, 0, len(summary.Projects))
	for _, p := range summary.Projects {
		items = append(items, ItemData{
			Name:     p.Name,
			Duration: internal.FormatDuration(p.Duration),
		})
	}
	return items, nil
}

func (a *App) GetLanguageNames() ([]ItemData, error) {
	if err := internal.CheckDataExists(); err != nil {
		return nil, nil
	}

	summary, err := internal.AllTimeSummary()
	if err != nil {
		return nil, err
	}

	summary = internal.FilterShort(summary)
	items := make([]ItemData, 0, len(summary.Languages))
	for _, l := range summary.Languages {
		items = append(items, ItemData{
			Name:     l.Name,
			Duration: internal.FormatDuration(l.Duration),
		})
	}
	return items, nil
}

func (a *App) GetProjectDetail(name string) (DetailResponse, error) {
	if err := internal.CheckDataExists(); err != nil {
		return DetailResponse{Name: name}, nil
	}

	allSummary, err := internal.AllTimeSummaryForProject(name)
	if err != nil {
		return DetailResponse{}, err
	}

	if allSummary.Total == 0 {
		return DetailResponse{Name: name, AllTime: "0h 00m", ThisMonth: "0h 00m", ThisWeek: "0h 00m"}, nil
	}

	mStart, mEnd := internal.MonthRange()
	monthEvents, err := internal.ReadEventsForRange(mStart, mEnd)
	if err != nil {
		return DetailResponse{}, err
	}
	monthSummary := internal.Summarize(internal.FilterByProject(internal.ComputeSessions(monthEvents), name))

	wStart, wEnd := internal.WeekRange()
	weekEvents, err := internal.ReadEventsForRange(wStart, wEnd)
	if err != nil {
		return DetailResponse{}, err
	}
	weekSummary := internal.Summarize(internal.FilterByProject(internal.ComputeSessions(weekEvents), name))

	allSummary = internal.FilterShort(allSummary)
	items := make([]ItemData, 0, len(allSummary.Languages))
	for _, l := range allSummary.Languages {
		pct := float64(0)
		if allSummary.Total > 0 {
			pct = float64(l.Duration) / float64(allSummary.Total) * 100
		}
		items = append(items, ItemData{
			Name:     l.Name,
			Duration: internal.FormatDuration(l.Duration),
			Percent:  pct,
		})
	}

	return DetailResponse{
		Name:      name,
		AllTime:   internal.FormatDuration(allSummary.Total),
		ThisMonth: internal.FormatDuration(monthSummary.Total),
		ThisWeek:  internal.FormatDuration(weekSummary.Total),
		Items:     items,
	}, nil
}

func (a *App) GetLanguageDetail(name string) (DetailResponse, error) {
	if err := internal.CheckDataExists(); err != nil {
		return DetailResponse{Name: name}, nil
	}

	allSummary, err := internal.AllTimeSummaryForLanguage(name)
	if err != nil {
		return DetailResponse{}, err
	}

	if allSummary.Total == 0 {
		return DetailResponse{Name: name, AllTime: "0h 00m", ThisMonth: "0h 00m", ThisWeek: "0h 00m"}, nil
	}

	mStart, mEnd := internal.MonthRange()
	monthEvents, err := internal.ReadEventsForRange(mStart, mEnd)
	if err != nil {
		return DetailResponse{}, err
	}
	monthSummary := internal.Summarize(internal.FilterByLanguage(internal.ComputeSessions(monthEvents), name))

	wStart, wEnd := internal.WeekRange()
	weekEvents, err := internal.ReadEventsForRange(wStart, wEnd)
	if err != nil {
		return DetailResponse{}, err
	}
	weekSummary := internal.Summarize(internal.FilterByLanguage(internal.ComputeSessions(weekEvents), name))

	allSummary = internal.FilterShort(allSummary)
	items := make([]ItemData, 0, len(allSummary.Projects))
	for _, p := range allSummary.Projects {
		pct := float64(0)
		if allSummary.Total > 0 {
			pct = float64(p.Duration) / float64(allSummary.Total) * 100
		}
		items = append(items, ItemData{
			Name:     p.Name,
			Duration: internal.FormatDuration(p.Duration),
			Percent:  pct,
		})
	}

	return DetailResponse{
		Name:      name,
		AllTime:   internal.FormatDuration(allSummary.Total),
		ThisMonth: internal.FormatDuration(monthSummary.Total),
		ThisWeek:  internal.FormatDuration(weekSummary.Total),
		Items:     items,
	}, nil
}
