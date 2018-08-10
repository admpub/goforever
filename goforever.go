package goforever

func New(config *Config) *Process {
	return &Process{
		Name:    "goforever",
		Args:    []string{},
		Command: "goforever",
		Pidfile: config.Pidfile,
		Logfile: config.Logfile,
		Errfile: config.Errfile,
		Respawn: 1,
	}
}

func NewConfg() *Config {
	return &Config{
		Processes: []*Process{},
	}
}
