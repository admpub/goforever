// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/admpub/goforever"
	cfg "github.com/admpub/goforever/config"
	httpF "github.com/admpub/goforever/http"
	"github.com/admpub/greq"
)

//go:embed goforever.toml
var exampleConfig []byte

var conf = flag.String("conf", "goforever.toml", "Path to config file.")
var config *cfg.Config
var daemon *goforever.Process
var version = `v0.0.1`

var Usage = func() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	usage := `
Commands
  list              List processes.
  show [name]       Show main proccess or named process.
  start [name]      Start main proccess or named process.
  stop [name]       Stop main proccess or named process.
  restart [name]    Restart main proccess or named process.
  version 			Show version information.
  example 			Display configuration file example.
  generate 			Generate sample configuration file.
`
	fmt.Fprintln(os.Stderr, usage)
}

func init() {
	flag.Usage = Usage
	flag.Parse()
	setConfig()
	daemon = &goforever.Process{
		Name:    "goforever",
		Args:    []string{},
		Command: filepath.Base(os.Args[0]),
		Pidfile: config.Pidfile,
		Logfile: config.Logfile,
		Errfile: config.Errfile,
		Respawn: 1,
		Debug:   true,
	}
}

func main() {
	if len(flag.Args()) > 0 {
		fmt.Printf("%s", Cli())
		return
	}
	httpS := httpF.New(config, daemon)
	if len(flag.Args()) == 0 {
		RunDaemon()
		httpS.HttpServer()
		return
	}
}

func Cli() string {
	var o []byte
	var err error
	sub := flag.Arg(0)
	name := flag.Arg(1)
	req := greq.New(host(), true)
	if sub == "list" {
		o, _, err = req.Get("/")
	} else if sub == "version" {
		o = []byte(version)
	} else if sub == "example" {
		o = exampleConfig
	} else if sub == "generate" {
		err = ioutil.WriteFile(*conf, exampleConfig, os.ModePerm)
		if err != nil {
			o = []byte(err.Error())
		} else {
			o = []byte(`The sample configuration file is generated successfully`)
		}
	} else if name == "" {
		if sub == "start" {
			daemon.Args = append(daemon.Args, os.Args[2:]...)
			return daemon.Start(daemon.Name)
		}
		_, _, err = daemon.Find()
		if err != nil {
			return fmt.Sprintf("Error: %s.\n", err)
		}
		if sub == "show" {
			return fmt.Sprintf("%s.\n", daemon.String())
		}
		if sub == "stop" {
			message := daemon.Stop()
			return message
		}
		if sub == "restart" {
			ch, message := daemon.Restart()
			fmt.Print(message)
			return fmt.Sprintf("%s\n", <-ch)
		}
	} else {
		path := fmt.Sprintf("/%s", name)
		switch sub {
		case "show":
			o, _, err = req.Get(path)
		case "start":
			o, _, err = req.Post(path, nil)
		case "stop":
			o, _, err = req.Delete(path)
		case "restart":
			o, _, err = req.Put(path, nil)
		}
	}
	if err != nil {
		fmt.Printf("Process error: %s", err)
	}
	return fmt.Sprintf("%s\n", o)
}

func RunDaemon() {
	daemon.Children = make(map[string]*goforever.Process, 0)
	for _, name := range config.Keys() {
		daemon.Children[name] = config.Get(name)
		daemon.Children[name].Debug = true
	}
	daemon.Run()
}

func setConfig() {
	var err error
	config, err = cfg.Load(*conf)
	if err != nil {
		log.Fatalf("%s", err)
		return
	}
	if config.Username == "" {
		log.Fatalf("Config error: %s", "Please provide a username.")
		return
	}
	if config.Password == "" {
		log.Fatalf("Config error: %s", "Please provide a password.")
		return
	}
	if config.Port == "" {
		config.Port = "2224"
	}
	if config.IP == "" {
		config.IP = "0.0.0.0"
	}
}

func host() string {
	scheme := "https"
	if len(config.TLSCertfile) == 0 || len(config.TLSKeyfile) == 0 {
		scheme = "http"
	}
	return fmt.Sprintf("%s://%s:%s@0.0.0.0:%s",
		scheme, config.Username, config.Password, config.Port,
	)
}
