package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		fmt.Println("Hi from worker")
		time.Sleep(5 * time.Second)
	}
}
