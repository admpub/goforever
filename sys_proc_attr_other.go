//go:build !windows

package goforever

import (
	"errors"
	"fmt"
	"os/user"
	"strconv"
	"syscall"
)

func SetSysProcAttr(attr *syscall.SysProcAttr, userName string, hideWindow bool) error {
	userInfo, err := user.Lookup(userName)
	if err != nil {
		return errors.New("failed to get user: " + err.Error())
	}
	uid, err := strconv.ParseUint(userInfo.Uid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to ParseUint(userInfo.Uid=%q): %w", userInfo.Uid, err)
	}
	gid, err := strconv.ParseUint(userInfo.Gid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to ParseUint(userInfo.Gid=%q): %w", userInfo.Gid, err)
	}
	attr.Credential = &syscall.Credential{
		Uid:         uint32(uid),
		Gid:         uint32(gid),
		NoSetGroups: true,
	}
	return err
}
