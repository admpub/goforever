package goforever

import (
	"os"
	"path/filepath"
)

func New(config *Config) *Process {
	return &Process{
		Name:     "goforever",
		Args:     []string{},
		Command:  filepath.Base(os.Args[0]),
		Pidfile:  config.Pidfile,
		Logfile:  config.Logfile,
		Errfile:  config.Errfile,
		Respawn:  1,
		Children: make(map[string]*Process, 0),
	}
}

func NewConfg() *Config {
	return &Config{
		Processes: []*Process{},
	}
}
