package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/admpub/goforever"
)

func NewProcess(config *Config) *goforever.Process {
	p := goforever.New()
	p.Pidfile = config.Pidfile
	p.Logfile = config.Logfile
	p.Errfile = config.Errfile
	return p
}

func New() *Config {
	return &Config{
		Processes: []*goforever.Process{},
	}
}

type Config struct {
	IP          string
	Port        string
	Username    string
	Password    string
	Daemonize   bool
	Pidfile     goforever.Pidfile
	Logfile     string
	Errfile     string
	TLSCertfile string
	TLSKeyfile  string
	Processes   []*goforever.Process `toml:"process"`
}

func (c *Config) Keys() []string {
	keys := []string{}
	for _, p := range c.Processes {
		keys = append(keys, p.Name)
	}
	return keys
}

func (c *Config) Get(key string) *goforever.Process {
	for _, p := range c.Processes {
		if p.Name == key {
			return p
		}
	}
	return nil
}

func (c *Config) Add(processes ...*goforever.Process) *Config {
	c.Processes = append(c.Processes, processes...)
	return c
}

func Load(file string) (*Config, error) {
	if !strings.HasPrefix(file, "/") {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		file = filepath.Join(wd, file)
	}
	c := &Config{}
	if _, err := toml.DecodeFile(file, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (cfg *Config) Export(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := toml.NewEncoder(f)
	for _, proc := range cfg.Processes {
		proc.Init()
	}
	return enc.Encode(cfg)
}
