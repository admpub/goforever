package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting")
	for {
		fmt.Println("Sleeping for 5s...")
		time.Sleep(5 * time.Second)
	}
}
