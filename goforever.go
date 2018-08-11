package goforever

import (
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

func StartChild(name string) (*Process, error) {
	return Default.StartChild(name)
}

func RestartChild(name string) (*Process, error) {
	return Default.RestartChild(name)
}

func StopChild(name string) error {
	return Default.StopChild(name)
}

func Child(name string) *Process {
	return Default.Children.Get(name)
}

func ChildKeys() []string {
	return Default.Children.Keys()
}

func Add(name string, procs *Process) *Process {
	return Default.Add(name, procs)
}

func Run() {
	Default.Run()
}
