// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package goforever

import (
	"os/user"
	"testing"
	"time"

	"github.com/webx-top/com"
)

func TestPidfile(t *testing.T) {
	p := &Process{
		Name:    "test",
		Pidfile: "test.pid",
		Debug:   true,
	}
	err := p.Pidfile.Write(100)
	if err != nil {
		t.Errorf("Error: %s.", err)
		return
	}
	ex := 100
	r := p.Pidfile.Read()
	if ex != r {
		t.Errorf("Expected %#v. Result %#v\n", ex, r)
	}

	s := p.Pidfile.Delete()
	if s != true {
		t.Error("Failed to remove pidfile.")
		return
	}
}

func TestProcessStart(t *testing.T) {
	p := &Process{
		Name:    "bash",
		Command: "/bin/bash",
		Args:    []string{"foo", "bar"},
		Pidfile: "echo.pid",
		Logfile: "debug.log",
		Errfile: "error.log",
		Respawn: 3,
		Debug:   true,
	}
	p.Start("bash")
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	p.Stop()
}

func TestProcessStartByUser(t *testing.T) {
	p := &Process{
		Name:    "bash",
		Command: "./example",
		Dir:     `./example`,
		Args:    []string{"foo", "bar"},
		Pidfile: "echo.pid",
		Logfile: "debug.log",
		Errfile: "error.log",
		Respawn: 3,
		Debug:   true,
		User:    `admpub`,
	}
	p.Start("bash")
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	time.Sleep(10 * time.Second)
	p.Stop()
}

func TestUser(t *testing.T) {
	u, err := user.Lookup(`nobody`)
	if err != nil {
		t.Error(err)
	} else {
		com.Dump(u)
	}
}
