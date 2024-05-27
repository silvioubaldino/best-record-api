package domain

import "time"

type Recording struct {
	ID        string
	StartTime time.Time
	EndTime   time.Time
	Status    string
	FilePath  string
}
