package goforever

import (
	"reflect"
	"syscall"
)

func NewSysProcAttr() *syscall.SysProcAttr {
	procAttr := &syscall.SysProcAttr{}
	r := reflect.ValueOf(procAttr)
	v := reflect.Indirect(r)
	if f := v.FieldByName(`Pdeathsig`); f.IsValid() {
		f.Set(reflect.ValueOf(syscall.Signal(0)))
	}
	if f := v.FieldByName(`Setpgid`); f.IsValid() {
		f.SetBool(true)
	}
	return procAttr
}
