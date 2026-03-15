package config

import (
	"errors"

	"github.com/spf13/viper"
)

type Config struct {
	PostgresURL    string
	BackupDir      string
	SlackWebhook   string
	BackupSchedule string
}

func LoadConfig() (*Config, error) {

	viper.SetDefault("backup_dir", "./backups")

	viper.AutomaticEnv()

	cfg := &Config{
		PostgresURL:    viper.GetString("POSTGRES_URL"),
		BackupDir:      viper.GetString("backup_dir"),
		SlackWebhook:   viper.GetString("SLACK_WEBHOOK_URL"),
		BackupSchedule: viper.GetString("BACKUP_SCHEDULE"),
	}

	return cfg, nil
}

func (c *Config) Validate() error {

	if c.PostgresURL == "" {
		return errors.New("POSTGRES_URL is required")
	}

	if c.BackupDir == "" {
		return errors.New("BACKUP_DIR is required")
	}

	if c.SlackWebhook == "" {
		return errors.New("SLACK_WEBHOOK_URL is required")
	}

	return nil
}
