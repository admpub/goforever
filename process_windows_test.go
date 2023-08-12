//go:build windows

package goforever

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
)

func init() {
	//debug = true
}

func TestWindowsSID(t *testing.T) {
	sid, domain, accType, err := windows.LookupSID(`Hank-MiniPC`, `PC`)
	if err != nil {
		t.Error(err)
	}
	t.Logf(`sid: %v, domain: %v, accType: %v`, sid, domain, accType)
	t.Logf(`uid: %v`, sid.String())
}

func TestWindowsToken(t *testing.T) {
	//token, err := getTokenByPid(2476)
	token, err := getToken(``, `PC`)
	if err != nil {
		t.Error(err)
	} else {
		defer token.Close()
	}
	t.Logf(`token: %v`, token)
}

func TestGetPidByUsername(t *testing.T) {
	pid, err := getPidByUsername(`Hank-MiniPC\test`)
	if err != nil {
		t.Error(err)
	}
	assert.Greater(t, pid, int32(0))
	t.Logf(`pid: %v`, pid)
}

func TestGetTokenByPid(t *testing.T) {
	pid, err := getPidByUsername(`Hank-MiniPC\test`)
	if err != nil {
		t.Error(err)
	}
	token, err := getTokenByPid(uint32(pid))
	if err != nil {
		t.Error(err)
	} else {
		defer token.Close()
	}
	t.Logf(`token: %v`, token)

}
