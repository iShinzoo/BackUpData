package cmd

import (
	"os"

	"github.com/iShinzoo/BackUpData/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "backup-cli",
	Short: "Databse backup CLI utility",
}

func Execute() {

	logg, err := logger.New()
	if err != nil {
		panic(err)
	}

	defer logg.Sync()

	if err := rootCmd.Execute(); err != nil {
		logg.Error("operation failed", zap.Error(err))
		os.Exit(1)
	}
}
