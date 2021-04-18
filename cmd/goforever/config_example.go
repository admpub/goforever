package main

var exampleConfig = `
ip = "127.0.0.1"
port = "2224"
username = "go"
password = "forever"
pidfile = "goforever.pid"
logfile = "goforever.log"
errfile = "goforever.log"

[[process]]
name = "example-panic"
command = "./example/example-panic"
pidfile = "example/example-panic.pid"
logfile = "example/logs/example-panic.debug.log"
errfile = "example/logs/example-panic.errors.log"
respawn = 1
delay = "1m"
ping = "30s"

[[process]]
name = "example"
dir = "../../"
env = ["aaa=aaa","bbb=bbb"]
command = "./example/example"
args = ["-name=foo"]
pidfile = "example/example.pid"
logfile = "example/logs/example.debug.log"
errfile = "example/logs/example.errors.log"
respawn = 1

[[process]]
name = "ll"
dir = "../../"
env = []
command = "ls -alh ."
args = []
pidfile = "example/example.pid"
logfile = "example/logs/example.debug.log"
errfile = "example/logs/example.errors.log"
respawn = 1
`
