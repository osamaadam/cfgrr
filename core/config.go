package core

type GitAuth struct {
	Type           string `mapstructure:"type"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
	Username       string `mapstructure:"username"`
}

type GitConfig struct {
	Remote string  `mapstructure:"remote"`
	Branch string  `mapstructure:"branch"`
	Auth   GitAuth `mapstructure:"auth"`
}

type Config struct {
	BackupDir  string    `mapstructure:"backup_dir"`
	MapFile    string    `mapstructure:"map_file"`
	IgnoreFile string    `mapstructure:"ignore_file"`
	Git        GitConfig `mapstructure:"git"`
}
