package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jiang-zhexin/animedb/internal/config"
	"github.com/jiang-zhexin/animedb/internal/model"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var syncCmd = &cobra.Command{
	Use:  "sync",
	RunE: runSync,
}

func init() {
	rootCmd.AddCommand(syncCmd)

	syncCmd.Flags().BoolP("dry-run", "n", false, "preview changes without making them")
}

func runSync(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	cachePath, err := cmd.Flags().GetString("db")
	logPath, err := cmd.Flags().GetString("log")
	if err != nil {
		return err
	}

	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer logFile.Close()
	slog.SetDefault(slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug})))

	viper.SetConfigType("toml")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		return err
	}

	m, err := model.LoadFromFile(cachePath)
	if err != nil {
		return fmt.Errorf("failed to load cache: %w", err)
	}
	defer m.SaveToFile(cachePath)

	m.SourceDirPath = conf.SourcePath
	m.TargetDirPath = conf.TargetPath

	dryRun, err := cmd.Flags().GetBool("dry-run")
	if err != nil {
		return err
	}
	m.DryRun = dryRun

	fmt.Println("Animedb - Sync Mode")
	fmt.Println("source path:", m.SourceDirPath)
	fmt.Println("target path:", m.TargetDirPath)
	fmt.Println("dry-run:", dryRun)
	fmt.Println()

	sourcePath, err := m.WalkFile(true)
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	for _, sp := range sourcePath {
		if m.IsProcessed(sp) {
			continue
		}

		pf, err := model.NewProcessedFiles(sp)
		if err != nil {
			return fmt.Errorf("fail to parser %s: %w", sp, err)
		}

		s, err := m.GetSubjectBySeriesName(pf.SeriesName)
		if err != nil {
			return fmt.Errorf("fail to sreach subject for %s: %w", pf.SeriesName, err)
		}

		pf.UpdateTargetPath(s)
		if err = m.UpdateProcessedFiles(pf); err != nil {
			return fmt.Errorf("fail to update processed files for %s: %w", sp, err)
		}
	}
	return nil
}
