package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting")
	for {
		fmt.Println(time.Now().Format("[2006-01-02 15:04:05]") + "Sleeping for 5s...")
		time.Sleep(5 * time.Second)
	}
}
