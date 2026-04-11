package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"
	"github.com/adrg/xdg"
	"github.com/jiang-zhexin/animedb/internal/app"
	"github.com/jiang-zhexin/animedb/internal/config"
	"github.com/jiang-zhexin/animedb/internal/model"
	"github.com/jiang-zhexin/animedb/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:  app.Appname,
	RunE: runRoot,
}

func init() {
	configPath, _ := xdg.ConfigFile(filepath.Join(app.Appname, "config.toml"))
	cachePath, _ := xdg.CacheFile(filepath.Join(app.Appname, "cache.json"))
	logPath, _ := xdg.StateFile(filepath.Join(app.Appname, "log.jsonl"))

	rootCmd.PersistentFlags().StringP("config", "c", configPath, "config file path")
	rootCmd.PersistentFlags().StringP("db", "d", cachePath, "db file path")
	rootCmd.PersistentFlags().StringP("log", "l", logPath, "log file path")
}

func runRoot(cmd *cobra.Command, args []string) error {
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

	p := tea.NewProgram(tui.NewRootModel(&tui.Ctx{Model: m, WindowSizeMsg: tea.WindowSizeMsg{}}, tui.NewMainMenuModel))
	_, err = p.Run()
	return err
}

func Execute() error {
	return rootCmd.Execute()
}
