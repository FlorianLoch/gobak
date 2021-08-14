package config

import "time"

type Config struct {
	Source              string        `yaml:"source" required:"" help:"Source directory to be monitored for changes."`
	Target              string        `yaml:"target" required:"" help:"Target directory to backup files to."`
	Interval            time.Duration `yaml:"interval" default:"5m" help:"Interval to check source dir for changes."`
	ExcludePattern      string        `yaml:"exclude_pattern" help:"Files with matching filename (including path) will not be considered"`
	// TODO: Actually this does not need to be configurable...
	MonitoredOperations []string      `yaml:"monitored_operations" default:"CREATE,WRITE,MOVE,REMOVE" help:"Operations to be monitored. Possibles values are 'CREATE', 'WRITE', 'REMOVE', 'RENAME', 'CHMOD' (not on Windows) and 'MOVE'."`
}
