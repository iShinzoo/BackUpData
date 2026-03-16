package metrics

import "github.com/prometheus/client_golang/prometheus"

var BackupSuccess = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "backup_success_total",
		Help: "Total successful backups",
	},
)

var BackupFailure = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "backup_failure_total",
		Help: "Total failed backups",
	},
)

func Register() {
	prometheus.MustRegister(BackupSuccess)
	prometheus.MustRegister(BackupFailure)
}
