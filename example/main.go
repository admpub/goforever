package main

import (
	"encoding/json"
	"fmt"
	"os/user"
	"time"
)

func main() {
	var userName string
	u, err := user.Current()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		b, _ := json.MarshalIndent(u, ``, `  `)
		fmt.Println("user:", string(b))
		userName = u.Username
	}
	fmt.Println("Starting for user:", userName)
	for {
		msg := time.Now().Format("[2006-01-02 15:04:05]") + "Sleeping for 5s..."
		fmt.Println(msg)
		time.Sleep(5 * time.Second)
	}
}
