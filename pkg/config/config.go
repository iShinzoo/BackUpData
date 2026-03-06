package config

import "github.com/spf13/viper"

type Config struct {
	PostgresURL string
	BackupDir   string
}

func LoadConfig() (*Config, error) {

	viper.SetDefault("backup_dir", "./backups")

	viper.AutomaticEnv()

	cfg := &Config{
		PostgresURL: viper.GetString("POSTGRES_URL"),
		BackupDir:   viper.GetString("backup_dir"),
	}

	return cfg, nil
}
