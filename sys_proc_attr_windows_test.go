//go:build windows

package goforever

import (
	"os"
	"testing"

	"golang.org/x/sys/windows"
)

func TestWindowsSID(t *testing.T) {
	sid, a, b, err := windows.LookupSID(`Hank-MiniPC`, `PC`)
	if err != nil {
		t.Error(err)
	}
	t.Logf(`sid: %v, %v, %v`, sid, a, b)
	t.Logf(`uid: %v`, sid.String())
}

func TestWindowsToken(t *testing.T) {
	token, err := getToken(os.Getpid())
	if err != nil {
		t.Error(err)
	}
	t.Logf(`token: %v`, token)
}
