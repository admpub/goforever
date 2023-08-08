package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/admpub/goforever"
	cfg "github.com/admpub/goforever/config"
)

func New(config *cfg.Config, daemon *goforever.Process) *HTTP {
	return &HTTP{
		config: config,
		daemon: daemon,
	}
}

type HTTP struct {
	config *cfg.Config
	daemon *goforever.Process
}

func (h *HTTP) HttpServer() {
	http.HandleFunc("/favicon.ico", http.NotFound)
	http.HandleFunc("/", h.AuthHandler(h.Handler))
	log.Printf("goforever serving port %s\n", h.config.Port)
	log.Printf("goforever serving IP %s\n", h.config.IP)
	bindAddress := fmt.Sprintf("%s:%s", h.config.IP, h.config.Port)
	if !h.isHttps() {
		if err := http.ListenAndServe(bindAddress, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		return
	}
	log.Println("SSL enabled.")
	if err := http.ListenAndServeTLS(bindAddress, "cert.pem", "key.pem", nil); err != nil {
		log.Fatal("ListenAndServeTLS: ", err)
	}
}

func (h *HTTP) isHttps() bool {
	return len(h.config.TLSCertfile) > 0 && len(h.config.TLSKeyfile) > 0
}

func (h *HTTP) Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "DELETE":
		h.DeleteHandler(w, r)
		return
	case "POST":
		h.PostHandler(w, r)
		return
	case "PUT":
		h.PutHandler(w, r)
		return
	case "GET":
		h.GetHandler(w, r)
		return
	}
}

type ProcessInfo struct {
	*goforever.Process
	Pid    int
	Status string
}

func (p ProcessInfo) String() string {
	js, err := json.Marshal(p)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(js)
}

func NewProcessInfo(p *goforever.Process) ProcessInfo {
	return ProcessInfo{Process: p, Pid: p.Pid(), Status: p.Status()}
}

func (h *HTTP) GetHandler(w http.ResponseWriter, r *http.Request) {
	var output []byte
	var err error
	switch r.URL.Path[1:] {
	case "":
		output, err = json.Marshal(h.daemon.ChildKeys())
	default:
		p := h.daemon.Child(r.URL.Path[1:])
		output, err = json.Marshal(NewProcessInfo(p))
	}
	if err != nil {
		log.Printf("Get Error: %#v\n", err)
		return
	}
	fmt.Fprintf(w, "%s", output)
}

func (h *HTTP) PostHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	ch, err := h.daemon.StartChild(name)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, "%s", NewProcessInfo(ch))
}

func (h *HTTP) PutHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	ch, err := h.daemon.RestartChild(name)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, "%s", NewProcessInfo(ch))
}

func (h *HTTP) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[1:]
	err := h.daemon.StopChild(name)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	fmt.Fprintf(w, "%s stopped.", name)
}

func (h *HTTP) AuthHandler(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		url := r.URL
		for k, v := range r.Header {
			fmt.Printf("  %s = %s\n", k, v[0])
		}
		auth, ok := r.Header["Authorization"]
		if !ok {
			log.Printf("Unauthorized access to %s\n", url)
			w.Header().Add("WWW-Authenticate", "basic realm=\"host\"")
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, "Not Authorized.")
			return
		}
		encoded := strings.SplitN(auth[0], " ", 2)
		if len(encoded) != 2 || encoded[0] != "Basic" {
			log.Printf("Strange Authorization %q\n", auth)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(encoded[1])
		if err != nil {
			log.Printf("Cannot decode %q: %s\n", auth, err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		parts := strings.SplitN(string(decoded), ":", 2)
		if len(parts) != 2 {
			log.Printf("Unknown format for credentials %q\n", decoded)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if parts[0] == h.config.Username && parts[1] == h.config.Password {
			fn(w, r)
			return
		}
		log.Printf("Unauthorized access to %s\n", url)
		w.Header().Add("WWW-Authenticate", "basic realm=\"host\"")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Not Authorized.")
	}
}
