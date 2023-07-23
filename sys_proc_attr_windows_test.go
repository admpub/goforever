//go:build windows

package goforever

import (
	"testing"
	
	"golang.org/x/sys/windows"
)

func TestWindowsSID(t *testing.T) {
	sid,a,b,err:=windows.LookupSID(`Hank-MiniPC`,`PC`)
	if err != nil {
		t.Error(err)
	}
	t.Logf(`sid: %v, %v, %v`,sid, a,b)
	t.Logf(`token: %v`,sid.Token)
	
}