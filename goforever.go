package goforever

import (
	"os"
	"path/filepath"
)

var Default = New()

func New() *Process {
	return &Process{
		Name:     "goforever",
		Args:     []string{},
		Command:  filepath.Base(os.Args[0]),
		Respawn:  1,
		Children: make(map[string]*Process, 0),
	}
}

func NewWithConfig(config *Config) *Process {
	p := New()
	p.Pidfile = config.Pidfile
	p.Logfile = config.Logfile
	p.Errfile = config.Errfile
	return p
}

func NewConfg() *Config {
	return &Config{
		Processes: []*Process{},
	}
}
