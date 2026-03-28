package entities

import "time"

type (
	Entry struct {
		ID          int64
		Text        string
		Tag         string
		Repo        string
		Branch      string
		CommitHash  string
		DurationSec int
		CreatedAt   time.Time
	}
)
