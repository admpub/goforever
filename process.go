// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package goforever

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	ps "github.com/admpub/go-ps"
)

var ping = "1m"

//RunProcess Run the process
func RunProcess(name string, p *Process) chan *Process {
	ch := make(chan *Process)
	go func() {
		proc, msg, err := p.Find()
		_, _ = msg, err
		// proc, err := ps.FindProcess(p.Pid)
		if proc == nil {
			p.Start(name)
		}
		p.ping(ping, func(time time.Duration, p *Process) {
			if p.Pid > 0 {
				p.respawns = 0
				fmt.Printf("%s refreshed after %s.\n", p.Name, time)
				p.Status = StatusRunning
			}
		})
		go p.watch()
		ch <- p
	}()
	return ch
}

const (
	StatusStarted   = `started`
	StatusRunning   = `running`
	StatusStopped   = `stopped`
	StatusRestarted = `restarted`
)

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
	Debug    bool
	x        *os.Process
	respawns int
	Children Children
}

func (p *Process) String() string {
	js, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(js)
}

//Find a process by name
func (p *Process) Find() (*os.Process, string, error) {
	if len(p.Pidfile) == 0 {
		return nil, "", errors.New("Pidfile is empty")
	}
	if pid := p.Pidfile.Read(); pid > 0 {
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
		p.Status = StatusRunning
		message := fmt.Sprintf("%s is %#v\n", p.Name, process.Pid)
		return process, message, nil
	}
	message := fmt.Sprintf("%s not running.\n", p.Name)
	return nil, message, fmt.Errorf("Could not find process %s", p.Name)
}

//Start the process
func (p *Process) Start(name string) string {
	p.Name = name
	// wd, _ := os.Getwd()
	wd := p.Dir
	if len(p.Dir) == 0 {
		wd, _ = os.Getwd()
	}
	abspath := filepath.Join(wd, p.Command)
	dirpath := filepath.Dir(abspath)
	basepath := filepath.Base(abspath)
	logPrefix := `[Process:` + name + `]`
	if p.Debug {
		log.Println(logPrefix+`Dir:`, dirpath)
	}
	files := []*os.File{
		os.Stdin,
		os.Stdout,
		os.Stderr,
	}
	if len(p.Logfile) > 0 {
		logDir := filepath.Dir(p.Logfile)
		os.MkdirAll(logDir, os.ModePerm)
		files[1] = NewLog(p.Logfile)
	}
	if len(p.Errfile) > 0 {
		logDir := filepath.Dir(p.Errfile)
		os.MkdirAll(logDir, os.ModePerm)
		files[2] = NewLog(p.Errfile)
	}
	proc := &os.ProcAttr{
		Dir:   dirpath,
		Env:   append(os.Environ()[:], p.Env...),
		Files: files,
	}
	args := append([]string{basepath}, p.Args...)
	basepath = "./" + basepath
	if p.Debug {
		log.Printf(logPrefix+"Args: %v\n", args)
	}
	process, err := os.StartProcess(basepath, args, proc)
	if err != nil {
		log.Fatalf("%s failed. %s\n", p.Name, err)
		return ""
	}
	err = p.Pidfile.Write(process.Pid)
	if err != nil {
		log.Printf("%s pidfile error: %s\n", p.Name, err)
		return ""
	}
	p.x = process
	p.Pid = process.Pid
	p.Status = StatusStarted
	return fmt.Sprintf("%s is %#v\n", p.Name, process.Pid)
}

//Stop the process
func (p *Process) Stop() string {
	if p.x != nil {
		// Initial code has the following comment: "p.x.Kill() this seems to cause trouble"
		// I want this to work on windows where AFAIK the existing code was not portable
		if err := p.x.Kill(); err != nil { //err := syscall.Kill(p.x.Pid, syscall.SIGTERM)
			log.Println(err)
		} else {
			fmt.Println("Stop command seemed to work")
		}
		p.Children.Stop()
	}
	p.release(StatusStopped)
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
	//p.Pidfile.Delete()
	p.Status = status
}

//Restart the process
func (p *Process) Restart() (chan *Process, string) {
	p.Stop()
	message := fmt.Sprintf("%s restarted.\n", p.Name)
	ch := RunProcess(p.Name, p)
	return ch, message
}

//Run callback on the process after given duration.
func (p *Process) ping(duration string, f func(t time.Duration, p *Process)) {
	if len(p.Ping) > 0 {
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
		p.release(StatusStopped)
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
		if p.Status == StatusStopped {
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
		if len(p.Delay) > 0 {
			t, _ := time.ParseDuration(p.Delay)
			time.Sleep(t)
		}
		p.Restart()
		p.Status = StatusRestarted
	case err := <-died:
		p.release("killed")
		log.Printf("%d %s killed = %#v\n", p.x.Pid, p.Name, err)
	}
}

//Run child processes
func (p *Process) Run() {
	for name, p := range p.Children {
		RunProcess(name, p)
	}
}

func (p *Process) StartChild(name string) (*Process, error) {
	cp := Child(name)
	if cp == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	cpp, _, _ := cp.Find()
	if cpp != nil {
		return nil, fmt.Errorf("%s already running", name)
	}
	ch := RunProcess(name, cp)
	procs := <-ch
	return procs, nil
}

func (p *Process) RestartChild(name string) (*Process, error) {
	cp := p.Child(name)
	if p == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	cp.Find()
	ch, _ := cp.Restart()
	procs := <-ch
	return procs, nil
}

func (p *Process) StopChild(name string) error {
	cp := p.Child(name)
	if cp == nil {
		return fmt.Errorf("%s does not exist", name)
	}
	cp.Find()
	cp.Stop()
	return nil
}

func (p *Process) Child(name string) *Process {
	return p.Children.Get(name)
}

func (p *Process) Add(name string, procs *Process, run ...bool) *Process {
	p.Children[name] = procs
	if len(run) > 0 && run[0] {
		RunProcess(name, procs)
	}
	return p
}

func (p *Process) ChildKeys() []string {
	return p.Children.Keys()
}
