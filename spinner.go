package main

import (
	"fmt"
	"os"
	"time"
)

var stopSpinner bool
var wait chan struct{} = make(chan struct{})

func spinner(delay time.Duration) {
	for !stopSpinner {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
	fmt.Fprint(os.Stdout, "\r \r")
	wait <- struct{}{}
}

func StartSpinner() {
	stopSpinner = false
	go spinner(100 * time.Millisecond)
}

func StopSpinner() {
	stopSpinner = true
	<-wait
}
