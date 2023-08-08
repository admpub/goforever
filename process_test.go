// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package goforever

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

var testuser string

func TestMain(t *testing.M) {
	u, err := user.Current()
	if err == nil {
		testuser = u.Username
	}
	flag.StringVar(&testuser, `user`, testuser, `--user `+testuser)
	t.Run()
}

// sudo go test -v -count=1 -run "TestProcessStartByUser" --user=aaa
func TestProcessStartByUser(t *testing.T) {
	os.Remove("debug.log")
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
		User:    testuser,
	}
	bin := `./example/example`
	if com.IsWindows {
		bin = `C:\Users\test\example.exe`
	}
	cmd := exec.Command(`go`, `build`, `-o`, bin, `./example`)
	err := cmd.Run()
	if err != nil {
		t.Fatal(err.Error())
	}
	p.Start("bash") // 此测试用例必须用root身份执行，否则报错：fork/exec ./example: operation not permitted
	ex := 0
	r := p.x.Pid
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	time.Sleep(10 * time.Second)
	b, err := os.ReadFile(p.Logfile)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Contains(t, string(b), `Starting for user: `+p.User)
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
