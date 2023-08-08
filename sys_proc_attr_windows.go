//go:build windows

package goforever

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/shirou/gopsutil/process"
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
	parts := strings.SplitN(userName, `\`, 2)
	var system string
	var err error
	if len(parts) != 2 {
		userName = parts[0]
	} else {
		system = parts[0]
		userName = parts[1]
	}
	token, err := getToken(system, userName)
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

func getToken(system, user string) (token syscall.Token, err error) {
	if len(system) == 0 {
		system, err = os.Hostname()
		if err != nil {
			err = fmt.Errorf(`failed to query os.Hostname(): %w`, err)
			return
		}
	}
	pid, err := getPidByUsername(system + `\` + user)
	if err != nil {
		return 0, err
	}
	return getTokenByPid(uint32(pid))
}

func getPidByUsername(username string, exename ...string) (int32, error) {
	var name string
	if len(exename) > 0 {
		name = exename[0]
	}
	pids, err := process.Pids()
	if err != nil {
		return 0, err
	}
	var pname, pusername string
	for _, pid := range pids {
		var proc *process.Process
		proc, err = process.NewProcess(pid)
		if err != nil {
			return 0, err
		}
		if len(name) > 0 {
			pname, err = proc.Name()
			if err != nil {
				return 0, fmt.Errorf(`failed to query proc.Name(): %w`, err)
			}
			if !strings.EqualFold(pname, name) {
				continue
			}
		}
		pusername, err = proc.Username()
		if err != nil {
			err = fmt.Errorf(`failed to query proc.Username(): %w`, err)
			continue
		}
		fmt.Println(`pname:`, pname, `username:`, username, `pusername:`, pusername)
		if strings.EqualFold(pusername, username) {
			return pid, nil
		}
	}
	if err != nil {
		return 0, err
	}
	err = fmt.Errorf(`the process(username: %v, name: %v) not found`, username, name)
	return 0, err
}

func getTokenByPid(pid uint32) (syscall.Token, error) {
	var err error
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, pid)
	if err != nil {
		return 0, fmt.Errorf("failed to OpenProcess(%d): %w", pid, err)
	}
	defer syscall.CloseHandle(handle)
	// Find process token via win32
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)
	if err != nil {
		return 0, fmt.Errorf("failed to OpenProcessToken(%d): %w", handle, err)
	}
	return token, err
}
