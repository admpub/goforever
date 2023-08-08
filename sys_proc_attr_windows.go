//go:build windows

package goforever

import (
	"fmt"
	"os"
	"syscall"

	"github.com/webx-top/com"
)

func buildOption(options map[string]interface{}) map[string]interface{} {
	if options == nil {
		options = map[string]interface{}{}
	}
	options[`HideWindow`] = false
	return options
}

func SetSysProcAttr(attr *syscall.SysProcAttr, userName string, options map[string]interface{}) (func(), error) {
	token, err := getToken(os.Getpid())
	if err != nil {
		return nil, err
	}
	if v, y := options[`HideWindow`]; y {
		attr.HideWindow = com.Bool(v)
	}
	attr.Token = token
	return func() {
		token.Close()
	}, nil
}

func getToken(pid int) (syscall.Token, error) {
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return token, fmt.Errorf("Token Process Error: %w", err)
	}
	defer syscall.CloseHandle(handle)

	// Find process token via win32
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)
	if err != nil {
		return token, fmt.Errorf("Open Token Process Error: %w", err)
	}
	return token, err
}
