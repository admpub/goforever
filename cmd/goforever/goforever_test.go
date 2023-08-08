// goforever - processes management
// Copyright (c) 2013 Garrett Woodworth (https://github.com/gwoo).

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func Test_main(t *testing.T) {
	setConfig()
	initDaemon()
	if daemon.Name != "goforever" {
		t.Error("Daemon name is not goforever")
	}
	daemon.Args = []string{"foo"}
	daemon.Start(daemon.Name)
	if daemon.Args[0] != "foo" {
		t.Error("First arg not foo")
	}
	daemon.Find()
	daemon.Stop()

	if com.IsWindows {
		assert.Equal(t, map[string]interface{}{
			`HideWindow`: false,
		}, config.Processes[len(config.Processes)-1].Options)
	}
}
