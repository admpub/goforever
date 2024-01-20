// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/admpub/goforever"
	cfg "github.com/admpub/goforever/config"
	httpF "github.com/admpub/goforever/http"
	"github.com/admpub/greq"
)

var conf = "goforever.toml"
var config *cfg.Config
var daemon *goforever.Process
var version = `v0.0.1`
var enableHTTP bool

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
  version           Show version information.
  example           Display configuration file example.
  generate          Generate sample configuration file.
`
	fmt.Fprintln(os.Stderr, usage)
}

func main() {
	flag.StringVar(&conf, "conf", conf, "Path to config file.")
	flag.BoolVar(&enableHTTP, "http", enableHTTP, "Enable HTTP server")
	flag.Usage = Usage
	flag.Parse()

	if len(flag.Args()) > 0 {
		sub := flag.Arg(0)
		switch sub {
		case "version":
			fmt.Println(version)
			return
		case "example":
			fmt.Println(exampleConfig)
			return
		case "generate":
			err := os.WriteFile(conf, []byte(exampleConfig), os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(`The sample configuration file is generated successfully`)
			}
			return
		}
	}

	setConfig()
	initDaemon()
	// if err := config.Export(conf + ".test"); err != nil {
	// 	log.Fatalln(err.Error())
	// }
	// return
	if len(flag.Args()) > 0 {
		fmt.Printf("%s", Cli())
		return
	}
	if len(flag.Args()) == 0 {
		RunDaemon()
		if enableHTTP {
			httpF.New(config, daemon).HttpServer()
		} else {
			<-make(chan struct{})
		}
		return
	}
}

func initDaemon() {
	daemon = goforever.NewProcess("goforever", filepath.Base(os.Args[0]))
	daemon.Pidfile = config.Pidfile
	daemon.Logfile = config.Logfile
	daemon.Errfile = config.Errfile
	daemon.Respawn = 1
	daemon.Debug = config.Debug
}

func Cli() string {
	var o []byte
	var err error
	sub := flag.Arg(0)
	name := flag.Arg(1)
	req := greq.New(host(), true)
	if sub == "list" {
		o, _, err = req.Get("/")
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
			proc, message := daemon.Restart()
			fmt.Print(message)
			return fmt.Sprintf("%s\n", proc)
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
	children := map[string]*goforever.Process{}
	for _, name := range config.Keys() {
		children[name] = config.Get(name).Init()
		children[name].Debug = config.Debug
	}
	daemon.SetChildren(children)
	daemon.Run()
}

func setConfig() {
	var err error
	config, err = cfg.Load(conf)
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
