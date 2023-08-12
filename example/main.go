package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/user"
	"time"
)

func main() {
	panicDelay := flag.Duration(`panic.delay`, 0, `--panic.delay`)
	flag.Parse()

	var userName string
	u, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.MarshalIndent(u, ``, `  `)
		fmt.Println("user:", string(b))
		userName = u.Username
	}
	start := time.Now()
	fmt.Println("Starting for user:", userName)
	for {
		if panicDelay != nil && *panicDelay > 0 && start.Before(time.Now().Add(-*panicDelay)) {
			panic(`Trigger panic: ` + time.Now().Format(time.DateTime))
		}
		msg := time.Now().Format("[2006-01-02 15:04:05]") + "Sleeping for 5s..."
		fmt.Println(msg)
		time.Sleep(5 * time.Second)
	}
}
