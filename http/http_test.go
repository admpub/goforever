package http

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/admpub/goforever"
	cfg "github.com/admpub/goforever/config"
	"github.com/admpub/greq"
)

var config = &cfg.Config{
	IP:       `0.0.0.0`,
	Port:     `2224`,
	Username: `admin`,
	Password: `admin`,
	Pidfile:  goforever.Pidfile(filepath.Join(os.TempDir(), `goforeverTest.pid`)),
}
var daemon = &goforever.Process{
	Name:    "goforever",
	Args:    []string{},
	Command: "goforever",
	Pidfile: config.Pidfile,
	Logfile: config.Logfile,
	Errfile: config.Errfile,
	Respawn: 1,
}

func TestListHandler(t *testing.T) {
	daemon.SetChildren(goforever.Children{
		"test": &goforever.Process{Name: "test"},
	})
	body, _ := newTestResponse("GET", "/", nil)
	ex := string([]byte(`["test"]`))
	r := string(body)
	if ex != r {
		t.Errorf("\nExpected = %v\nResult = %v\n", ex, r)
	}
}

func TestShowHandler(t *testing.T) {
	daemon.SetChildren(goforever.Children{
		"test": &goforever.Process{Name: "test"},
	})
	body, _ := newTestResponse("GET", "/test", nil)
	e := []byte(`{"Name":"test","Command":"","Env":null,"Dir":"","Args":null,"User":"","HideWindow":false,"Pidfile":"","Logfile":"","Errfile":"","Respawn":0,"Delay":"","Ping":"","Debug":false,"Pid":0,"Status":""}`)
	if !bytes.Equal(e, body) {
		t.Errorf("\nExpected = %s\nResult = %s\n", e, body)
	}
}

func TestPostHandler(t *testing.T) {
	pidfile := filepath.Join(os.TempDir(), `goforeverTestEcho.pid`)
	daemon.SetChildren(goforever.Children{
		"test": &goforever.Process{Name: "test", Command: "/bin/echo", Args: []string{"woohoo"}, Pidfile: goforever.Pidfile(pidfile)},
	})
	body, _ := newTestResponse("POST", "/test", nil)
	b, _ := os.ReadFile(pidfile)
	pid := string(b)
	e := []byte(`{"Name":"test","Command":"/bin/echo","Env":null,"Dir":"","Args":["woohoo"],"User":"","HideWindow":false,"Pidfile":"` + pidfile + `","Logfile":"","Errfile":"","Respawn":0,"Delay":"","Ping":"","Debug":false,"Pid":` + pid + `,"Status":"started"}`)
	if !bytes.Equal(e, body) {
		t.Errorf("\nExpected = %s\nResult = %s\n", e, body)
	}
}

func newTestResponse(method string, path string, body io.Reader) ([]byte, *http.Response) {
	ts := httptest.NewServer(http.HandlerFunc(New(config, daemon).Handler))
	defer ts.Close()
	url := ts.URL + path
	b, r, _ := greq.Do(method, url, nil, body)
	return b, r
}
