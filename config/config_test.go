package config

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	r, err := Load("../cmd/goforever/goforever.toml")

	if err != nil {
		t.Errorf("Error creating config %s.", err)
		return
	}
	if r == nil {
		t.Errorf("Expected %#v. Result %#v\n", r, nil)
	}
}

func TestConfigGet(t *testing.T) {
	c, _ := Load("../cmd/goforever/goforever.toml")
	ex := "example/example.pid"
	r := string(c.Get("example").Pidfile)
	if ex != r {
		t.Errorf("Expected %#v. Result %#v\n", ex, r)
	}
}

func TestConfigKeys(t *testing.T) {
	c, _ := Load("../cmd/goforever/goforever.toml")
	ex := []string{"example", "example-panic", "ll"}
	r := c.Keys()
	if len(ex) != len(r) {
		t.Errorf("Expected %#v. Result %#v\n", ex, r)
	}
}
