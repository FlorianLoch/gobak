package config

import "time"

type Config struct {
	Source string `yaml:"source" required:"" help:"Source directory to be monitored for changes."`
	Target string `yaml:"target" required:"" help:"Target directory to backup files to."`
	Interval time.Duration `yaml:"interval" default:"5m" help:"Interval to check source dir for changes."`
}
