// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"
"path/filepath"
	"github.com/gooops/go-ps"
)

var ping = "1m"

//Run the process
func RunProcess(name string, p *Process) chan *Process {
	ch := make(chan *Process)
	go func() {
		proc, msg, err := p.find()
		_, _ = msg, err
		// proc, err := ps.FindProcess(p.Pid)
		if proc == nil {
			p.start(name)
		}
		p.ping(ping, func(time time.Duration, p *Process) {
			if p.Pid > 0 {
				p.respawns = 0
				fmt.Printf("%s refreshed after %s.\n", p.Name, time)
				p.Status = "running"
			}
		})
		go p.watch()
		ch <- p
	}()
	return ch
}

type Process struct {
	Name     string
	Command  string
	Env      []string
	Dir      string
	Args     []string
	Pidfile  Pidfile
	Logfile  string
	Errfile  string
	Path     string
	Respawn  int
	Delay    string
	Ping     string
	Pid      int
	Status   string
	x        *os.Process
	respawns int
	children children
}

func (p *Process) String() string {
	js, err := json.Marshal(p)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(js)
}

//Find a process by name
func (p *Process) find() (*os.Process, string, error) {
	if p.Pidfile == "" {
		return nil, "", errors.New("Pidfile is empty.")
	}
	if pid := p.Pidfile.read(); pid > 0 {
		proc, err := ps.FindProcess(pid)
		if err != nil || proc == nil {
			return nil, "", err
		}
		process, err := os.FindProcess(pid)
		if err != nil {
			return nil, "", err
		}
		p.x = process
		p.Pid = process.Pid
		p.Status = "running"
		message := fmt.Sprintf("%s is %#v\n", p.Name, process.Pid)
		return process, message, nil
	}
	message := fmt.Sprintf("%s not running.\n", p.Name)
	return nil, message, errors.New(fmt.Sprintf("Could not find process %s.", p.Name))
}

//Start the process
func (p *Process) start(name string) string {
	p.Name = name
	// wd, _ := os.Getwd()
	wd := p.Dir
	if p.Dir == "" {
		wd, _ = os.Getwd()
	}
	abspath := filepath.Join(wd, p.Command)
	dirpath := filepath.Dir(abspath)
	basepath := filepath.Base(abspath)
	fmt.Println(dirpath)
	proc := &os.ProcAttr{
		Dir: dirpath,
		Env: append(os.Environ()[:], p.Env...),
		Files: []*os.File{
			os.Stdin,
			NewLog(p.Logfile),
			NewLog(p.Errfile),
		},
	}
	args := append([]string{basepath}, p.Args...)
	basepath = "./" + basepath
	fmt.Printf("Args: %v %v %v", basepath, args, proc)
	process, err := os.StartProcess(basepath, args, proc)
	if err != nil {
		log.Fatalf("%s failed. %s\n", p.Name, err)
		return ""
	}
	err = p.Pidfile.write(process.Pid)
	if err != nil {
		log.Printf("%s pidfile error: %s\n", p.Name, err)
		return ""
	}
	p.x = process
	p.Pid = process.Pid
	p.Status = "started"
	return fmt.Sprintf("%s is %#v\n", p.Name, process.Pid)
}

//Stop the process
func (p *Process) stop() string {
	if p.x != nil {
		// Initial code has the following comment: "p.x.Kill() this seems to cause trouble"
		// I want this to work on windows where AFAIK the existing code was not portable
		if err := p.x.Kill(); err != nil {
			log.Println(err)
		} else {
			fmt.Println("Stop command seemed to work")
		}
		// p.children.stop("all")
	}
	p.release("stopped")
	message := fmt.Sprintf("%s stopped.\n", p.Name)
	return message
}

//Release process and remove pidfile
func (p *Process) release(status string) {
	// debug.PrintStack()
	if p.x != nil {
		p.x.Release()
	}
	p.Pid = 0
	// 去掉删除pid文件的动作，用于goforever进程重启后继续监控，防止启动重复进程
	// p.Pidfile.delete()
	p.Status = status
}

//Restart the process
func (p *Process) restart() (chan *Process, string) {
	p.stop()
	message := fmt.Sprintf("%s restarted.\n", p.Name)
	ch := RunProcess(p.Name, p)
	return ch, message
}

//Run callback on the process after given duration.
func (p *Process) ping(duration string, f func(t time.Duration, p *Process)) {
	if p.Ping != "" {
		duration = p.Ping
	}
	t, err := time.ParseDuration(duration)
	if err != nil {
		t, _ = time.ParseDuration(ping)
	}
	go func() {
		select {
		case <-time.After(t):
			f(t, p)
		}
	}()
}

//Watch the process
func (p *Process) watch() {
	if p.x == nil {
		p.release("stopped")
		return
	}
	status := make(chan *os.ProcessState)
	died := make(chan error)
	go func() {
		// state, err := p.x.Wait()
		proc, err := ps.FindProcess(p.Pid)
		var ppid int
		var state = &os.ProcessState{}
		if proc != nil {
			ppid = proc.PPid()
		}
		// 如果是当前进程fork的子进程，则阻塞等待获取子进程状态，否则循环检测进程状态（1s一次，直到状态变更）
		if ppid == os.Getpid() {
			state, err = p.x.Wait()
		} else {
			for {
				time.Sleep(1 * time.Second)
				proc, err = ps.FindProcess(p.Pid)
				if err != nil || proc == nil {
					break
				}
			}
		}
		if err != nil {
			died <- err
			return
		}
		status <- state
	}()
	select {
	case s := <-status:
		if p.Status == "stopped" {
			return
		}

		fmt.Fprintf(os.Stderr, "%s %s\n", p.Name, s)
		fmt.Fprintf(os.Stderr, "%s success = %#v\n", p.Name, s.Success())
		fmt.Fprintf(os.Stderr, "%s exited = %#v\n", p.Name, s.Exited())
		p.respawns++
		if p.respawns > p.Respawn {
			p.release("exited")
			log.Printf("%s respawn limit reached.\n", p.Name)
			return
		}
		fmt.Fprintf(os.Stderr, "%s respawns = %#v\n", p.Name, p.respawns)
		if p.Delay != "" {
			t, _ := time.ParseDuration(p.Delay)
			time.Sleep(t)
		}
		p.restart()
		p.Status = "restarted"
	case err := <-died:
		p.release("killed")
		log.Printf("%d %s killed = %#v", p.x.Pid, p.Name, err)
	}
}

//Run child processes
func (p *Process) run() {
	for name, p := range p.children {
		RunProcess(name, p)
	}
}

//Child processes.
type children map[string]*Process

//Stringify
func (c children) String() string {
	js, err := json.Marshal(c)
	if err != nil {
		log.Print(err)
		return ""
	}
	return string(js)
}

//Get child processes names.
func (c children) keys() []string {
	keys := []string{}
	for k, _ := range c {
		keys = append(keys, k)
	}
	return keys
}

//Get child process.
func (c children) get(key string) *Process {
	if v, ok := c[key]; ok {
		return v
	}
	return nil
}

func (c children) stop(name string) {
	if name == "all" {
		for name, p := range c {
			p.stop()
			delete(c, name)
		}
		return
	}
	p := c.get(name)
	p.stop()
	delete(c, name)
}

type Pidfile string

//Read the pidfile.
func (f *Pidfile) read() int {
	data, err := ioutil.ReadFile(string(*f))
	if err != nil {
		return 0
	}
	pid, err := strconv.ParseInt(string(data), 0, 32)
	if err != nil {
		return 0
	}
	return int(pid)
}

//Write the pidfile.
func (f *Pidfile) write(data int) error {
	err := ioutil.WriteFile(string(*f), []byte(strconv.Itoa(data)), 0660)
	if err != nil {
		return err
	}
	return nil
}

//Delete the pidfile
func (f *Pidfile) delete() bool {
	_, err := os.Stat(string(*f))
	if err != nil {
		return true
	}
	err = os.Remove(string(*f))
	if err == nil {
		return true
	}
	return false
}

//Create a new file for logging
func NewLog(path string) *os.File {
	if path == "" {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("%s", err)
		return nil
	}
	return file
}
