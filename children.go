package goforever

import (
	"encoding/json"
	"log"
	"sync"
)

func NewChildren(m map[string]*Process) *Children {
	if m == nil {
		m = map[string]*Process{}
	}
	return &Children{m: m, l: &sync.RWMutex{}}
}

// Children Child processes.
type Children struct {
	m map[string]*Process
	l *sync.RWMutex
}

// String Stringify
func (c Children) String() string {
	js, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(js)
}

// Keys Get child processes names.
func (c Children) Keys() []string {
	keys := []string{}
	c.l.RLock()
	for k := range c.m {
		keys = append(keys, k)
	}
	c.l.RUnlock()
	return keys
}

func (c Children) Run() {
	c.l.RLock()
	for name, p := range c.m {
		RunProcess(name, p)
	}
	c.l.RUnlock()
}

// Get child process.
func (c Children) Get(key string) *Process {
	c.l.RLock()
	v, ok := c.m[key]
	c.l.RUnlock()
	if ok {
		return v
	}
	return nil
}

// Get child process.
func (c Children) Set(key string, proc *Process) {
	c.l.Lock()
	c.m[key] = proc
	c.l.Unlock()
}

func (c Children) Stop(names ...string) {
	if len(names) < 1 {
		c.l.Lock()
		for name, p := range c.m {
			p.Stop()
			delete(c.m, name)
		}
		c.l.Unlock()
		return
	}
	name := names[0]
	p := c.Get(name)
	if p == nil {
		return
	}

	c.l.Lock()
	delete(c.m, name)
	c.l.Unlock()
}
