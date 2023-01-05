package core

type Config struct {
	BackupDir  string `mapstructure:"backup_dir"`
	MapFile    string `mapstructure:"map_file"`
	IgnoreFile string `mapstructure:"ignore_file"`
}
