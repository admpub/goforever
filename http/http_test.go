package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/admpub/goforever"
	cfg "github.com/admpub/goforever/config"
	"github.com/gwoo/greq"
)

var config = &cfg.Config{
	IP:       `0.0.0.0`,
	Port:     `2224`,
	Username: `admin`,
	Password: `admin`,
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
	daemon.Children = goforever.Children{
		"test": &goforever.Process{Name: "test"},
	}
	body, _ := newTestResponse("GET", "/", nil)
	ex := fmt.Sprintf("%s", string([]byte(`["test"]`)))
	r := fmt.Sprintf("%s", string(body))
	if ex != r {
		t.Errorf("\nExpected = %v\nResult = %v\n", ex, r)
	}
}

func TestShowHandler(t *testing.T) {
	daemon.Children = goforever.Children{
		"test": &goforever.Process{Name: "test"},
	}
	body, _ := newTestResponse("GET", "/test", nil)
	e := []byte(`{"Name":"test","Command":"","Args":null,"Pidfile":"","Logfile":"","Errfile":"","Path":"","Respawn":0,"Delay":"","Ping":"","Pid":0,"Status":""}`)
	ex := fmt.Sprintf("%s", e)
	r := fmt.Sprintf("%s", body)
	if ex != r {
		t.Errorf("\nExpected = %v\nResult = %v\n", ex, r)
	}
}

func TestPostHandler(t *testing.T) {
	daemon.Children = goforever.Children{
		"test": &goforever.Process{Name: "test", Command: "/bin/echo", Args: []string{"woohoo"}},
	}
	body, _ := newTestResponse("POST", "/test", nil)
	e := []byte(`{"Name":"test","Command":"/bin/echo","Args":["woohoo"],"Pidfile":"","Logfile":"","Errfile":"","Path":"","Respawn":0,"Delay":"","Ping":"","Pid":0,"Status":"stopped"}`)
	ex := fmt.Sprintf("%s", e)
	r := fmt.Sprintf("%s", body)
	if ex != r {
		t.Errorf("\nExpected = %v\nResult = %v\n", ex, r)
	}
}

func newTestResponse(method string, path string, body io.Reader) ([]byte, *http.Response) {
	ts := httptest.NewServer(http.HandlerFunc(New(config, daemon).Handler))
	defer ts.Close()
	url := ts.URL + path
	b, r, _ := greq.Do(method, url, nil, body)
	return b, r
}
