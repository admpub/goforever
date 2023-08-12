// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package goforever

import (
	"flag"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
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
		Name:    "example",
		Command: "./example",
		Dir:     `./example`,
		Args:    []string{"foo", "bar"},
		Pidfile: "echo.pid",
		Logfile: "debug.log",
		Errfile: "error.log",
		Respawn: 3,
		Debug:   true,
		Ping:    "1s",
	}
	bin := `./example/example`
	cmd := exec.Command(`go`, `build`, `-o`, bin, `./example`)
	err := cmd.Run()
	if err != nil {
		t.Error(cmd.String() + `: ` + err.Error())
	}
	p.Start(p.Name)
	//<-RunProcess(p.Name, p)
	//time.Sleep(30 * time.Second)
	ex := 0
	r := p.Pid()
	if ex >= r {
		t.Errorf("Expected %#v < %#v\n", ex, r)
	}
	p.Stop()
}

var testuser string = `hank-minipc\test`
var testpass string

func TestMain(t *testing.M) {
	if len(testuser) == 0 {
		u, err := user.Current()
		if err == nil {
			testuser = u.Username
		}
	}
	flag.StringVar(&testuser, `user`, testuser, `--user `+testuser)
	flag.StringVar(&testpass, `pass`, testpass, `--pass `+testpass)
	t.Run()
}

// sudo go test -v -count=1 -run "TestProcessStartByUser" --user=aaa --pass=yourWindowsPassword
func TestProcessStartByUser(t *testing.T) {
	os.Remove("debug.log")
	p := &Process{
		Name:    "exampleByUser",
		Command: "./example",
		Dir:     `./example`,
		Args:    []string{},
		Pidfile: "echo.pid",
		Logfile: "debug.log",
		Errfile: "error.log",
		Respawn: 3,
		Debug:   true,
		User:    testuser,
		Options: map[string]interface{}{
			`HideWindow`: true,
			`Password`:   testpass,
		},
	}
	//com.Dump(p)
	bin := `./example/example`
	if com.IsWindows {
		bin = `C:\Users\test\example.exe`
		p.Dir = `C:\Users\test`
		p.Command = bin
		p.Pidfile = Pidfile(filepath.Join(p.Dir, "echo.pid"))
		p.Logfile = filepath.Join(p.Dir, "debug.log")
		p.Errfile = filepath.Join(p.Dir, "error.log")
	}
	cmd := exec.Command(`go`, `build`, `-o`, bin, `./example`)
	err := cmd.Run()
	if err != nil {
		t.Error(cmd.String() + `: ` + err.Error())
	}
	p.Start("exampleByUser") // 此测试用例必须用root身份执行，否则报错：fork/exec ./example: operation not permitted
	ex := 0
	r := p.Pid()
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
