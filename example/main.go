package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"time"
)

func main() {
	f, err := os.OpenFile(`test.log`, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	var userName string
	u, err := user.Current()
	if err != nil {
		f.WriteString(err.Error() + "\n")
		fmt.Println(err.Error())
	} else {
		b, _ := json.MarshalIndent(u, ``, `  `)
		f.WriteString("user:\n")
		f.Write(b)
		f.WriteString("\n")
		userName = u.Username
	}
	f.WriteString("Starting for user:" + userName + "\n")
	fmt.Println("Starting for user:", userName)
	for {
		msg := time.Now().Format("[2006-01-02 15:04:05]") + "Sleeping for 5s..."
		f.WriteString(msg + "\n")
		fmt.Println(msg)
		time.Sleep(5 * time.Second)
	}
}
