package core

import "time"

type BackupJob struct {
	Name string
	URL  string
}

type BackupResult struct {
	Name     string
	Status   string
	Error    error
	Duration time.Duration
	Size     int64
}
