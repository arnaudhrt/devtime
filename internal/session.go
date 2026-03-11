package internal

import (
	"sort"
	"time"
)

const (
	sessionGap         = 5 * time.Minute
	singleEventDuration = 30 * time.Second
)

type Session struct {
	Project  string
	Language string
	Editor   string
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// ComputeSessions takes a slice of events and returns sessions.
// Events are sorted by timestamp ascending before processing.
func ComputeSessions(events []Event) []Session {
	if len(events) == 0 {
		return nil
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Timestamp.Before(events[j].Timestamp)
	})

	var sessions []Session

	sessionStart := events[0].Timestamp
	sessionEnd := events[0].Timestamp
	project := events[0].Project
	lang := events[0].Language
	editor := events[0].Editor
	needNewSession := false

	// If the first event is a blur, record a single-event session and move on
	if events[0].Type == "blur" {
		sessions = append(sessions, Session{
			Project:  project,
			Language: lang,
			Editor:   editor,
			Start:    sessionStart,
			End:      sessionEnd,
			Duration: singleEventDuration,
		})
		needNewSession = true
	}

	for i := 1; i < len(events); i++ {
		ev := events[i]

		if needNewSession {
			// Start a fresh session from this event
			sessionStart = ev.Timestamp
			sessionEnd = ev.Timestamp
			project = ev.Project
			lang = ev.Language
			editor = ev.Editor
			needNewSession = false

			if ev.Type == "blur" {
				sessions = append(sessions, Session{
					Project:  project,
					Language: lang,
					Editor:   editor,
					Start:    sessionStart,
					End:      sessionEnd,
					Duration: singleEventDuration,
				})
				needNewSession = true
			}
			continue
		}

		gap := ev.Timestamp.Sub(sessionEnd)
		projectChanged := ev.Project != project
		langChanged := ev.Language != lang

		// Close current session and start a new one if:
		// - gap > 5 minutes
		// - project or language changed
		if gap > sessionGap || projectChanged || langChanged {
			dur := sessionEnd.Sub(sessionStart)
			if dur == 0 {
				dur = singleEventDuration
			}
			sessions = append(sessions, Session{
				Project:  project,
				Language: lang,
				Editor:   editor,
				Start:    sessionStart,
				End:      sessionEnd,
				Duration: dur,
			})
			sessionStart = ev.Timestamp
			sessionEnd = ev.Timestamp
			project = ev.Project
			lang = ev.Language
			editor = ev.Editor
		}

		if ev.Type == "blur" {
			// Blur ends the current session
			sessionEnd = ev.Timestamp
			dur := sessionEnd.Sub(sessionStart)
			if dur == 0 {
				dur = singleEventDuration
			}
			sessions = append(sessions, Session{
				Project:  project,
				Language: lang,
				Editor:   editor,
				Start:    sessionStart,
				End:      sessionEnd,
				Duration: dur,
			})
			needNewSession = true
		} else {
			// heartbeat or focus: continue the session
			sessionEnd = ev.Timestamp
		}
	}

	// Close any remaining open session
	if !needNewSession {
		dur := sessionEnd.Sub(sessionStart)
		if dur == 0 {
			dur = singleEventDuration
		}
		sessions = append(sessions, Session{
			Project:  project,
			Language: lang,
			Editor:   editor,
			Start:    sessionStart,
			End:      sessionEnd,
			Duration: dur,
		})
	}

	return sessions
}
