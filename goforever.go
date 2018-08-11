package goforever

import (
	"fmt"
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

func Start(name string) (*Process, error) {
	p := Get(name)
	if p == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	cp, _, _ := p.Find()
	if cp != nil {
		return nil, fmt.Errorf("%s already running", name)
	}
	ch := RunProcess(name, p)
	procs := <-ch
	return procs, nil
}

func Restart(name string) (*Process, error) {
	p := Get(name)
	if p == nil {
		return nil, fmt.Errorf("%s does not exist", name)
	}
	p.Find()
	ch, _ := p.Restart()
	procs := <-ch
	return procs, nil
}

func Stop(name string) error {
	p := Get(name)
	if p == nil {
		return fmt.Errorf("%s does not exist", name)
	}
	p.Find()
	p.Stop()
	return nil
}

func Get(name string) *Process {
	return Default.Children.Get(name)
}

func Keys() []string {
	return Default.Children.Keys()
}
