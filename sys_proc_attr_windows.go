//go:build windows

package goforever

import (
	"golang.org/x/sys/windows"
	"syscall"
)

func (p *Process) setSysProcAttr(attr *syscall.SysProcAttr) error {
	attr.Token = windows.Token(0)
}
