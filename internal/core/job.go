package core

type BackupJob struct {
	Name string
	URL  string
}

type BackupResult struct {
	Name   string
	Status string
	Error  error
}
