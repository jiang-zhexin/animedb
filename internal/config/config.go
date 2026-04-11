package config

type Config struct {
	SourcePath string `mapstructure:"source_path"`
	TargetPath string `mapstructure:"target_path"`
}
